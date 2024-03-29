package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: leaven input-file.ll")
		os.Exit(1)
	}

	inFile := os.Args[1]
	m, err := asm.ParseFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	outFile := strings.TrimSuffix(inFile, ".ll") + ".go"
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = compile(out, m)
	if err != nil {
		log.Fatal(err)
	}
}

func compile(out io.Writer, m *ir.Module) error {
	fmt.Fprint(out, "package main\n\n")

	for _, t := range m.TypeDefs {
		name := TypeName(t)
		if name == "" {
			continue
		}
		if strings.Contains(name, ".") {
			// It's a definition that's beeen replaced by a reference to a standard-library type.
			continue
		}

		def, err := TypeDefinition(t)
		if err != nil {
			return fmt.Errorf("error generating type definition for %v: %v", t, err)
		}

		fmt.Fprintf(out, "type %s %s\n\n", name, def)
	}

	for _, g := range m.Globals {
		if g.Init == nil {
			// Just a declaration; skip it.
			continue
		}
		t, err := TypeSpec(g.ContentType)
		if err != nil {
			return fmt.Errorf("error translating type (%v): %v", g.ContentType, err)
		}
		val, err := FormatValue(g.Init)
		if err != nil {
			return fmt.Errorf("error translating initializer (%v): %v", g.Init, err)
		}
		fmt.Fprintf(out, "var %s %s = %s\n\n", VariableName(g), t, val)
	}

	for _, f := range m.Funcs {
		if f.Blocks == nil {
			// Just a declaration, not a definition; skip it.
			continue
		}

		fixMalloc(f)

		name := VariableName(f)

		if name == "main" {
			fmt.Fprintln(out, "func main() {")
		} else {
			fmt.Fprintf(out, "func %s(", name)
			for i, p := range f.Params {
				if i > 0 {
					fmt.Fprint(out, ", ")
				}
				pt, err := TypeSpec(p.Typ)
				if err != nil {
					return fmt.Errorf("error translating type for parameter %d of %s: %v", i, f.Name(), err)
				}
				fmt.Fprintf(out, "%s %s", VariableName(p), pt)
			}
			if f.Sig.Variadic {
				if len(f.Params) > 0 {
					fmt.Fprint(out, ", ")
				}
				fmt.Fprint(out, "varargs ...interface{}")
			}
			fmt.Fprint(out, ") ")
			rt := f.Sig.RetType
			if !types.Equal(rt, types.Void) {
				retType, err := TypeSpec(rt)
				if err != nil {
					return fmt.Errorf("error translating return type for %s: %v", f.Name(), err)
				}
				fmt.Fprintf(out, "%s ", retType)
			}
			fmt.Fprint(out, "{\n")
		}

		// Declare variables.
		vars := make(map[string][]string)
		var allVars []string
		for _, b := range f.Blocks {
			for _, inst := range b.Insts {
				if inst, ok := inst.(value.Named); ok {
					if types.Equal(inst.Type(), types.Void) {
						continue
					}
					t, err := TypeSpec(inst.Type())
					if err != nil {
						return fmt.Errorf("error translating type of %s in %s: %v", inst.Ident(), f.Name(), err)
					}
					vars[t] = append(vars[t], VariableName(inst))
					allVars = append(allVars, VariableName(inst))
				}
			}
		}
		varTypes := make([]string, 0, len(vars))
		for t := range vars {
			varTypes = append(varTypes, t)
		}
		sort.Strings(varTypes)
		for _, t := range varTypes {
			fmt.Fprintf(out, "\tvar %s %s\n", strings.Join(vars[t], ", "), t)
		}
		if len(vars) > 0 {
			fmt.Fprintln(out)
			// Get rid of unused-variable errors.
			for i := range allVars {
				if i == 0 {
					fmt.Fprint(out, "\t_")
				} else {
					fmt.Fprint(out, ", _")
				}
			}
			fmt.Fprintf(out, " = %s\n\n", strings.Join(allVars, ", "))
		}

		// Translate instructions.
		for i, b := range f.Blocks {
			if i != 0 {
				fmt.Fprintf(out, "\n%s:\n", BlockName(b))
			}
			for _, inst := range b.Insts {
				if _, ok := inst.(*ir.InstPhi); ok {
					continue
				}
				translated, err := TranslateInstruction(inst)
				if err != nil {
					return fmt.Errorf("error translating %q: %v", inst.LLString(), err)
				}
				if translated != "" {
					fmt.Fprintf(out, "\t%s\n", translated)
				}
			}
			switch term := b.Term.(type) {
			case *ir.TermBr:
				phis, err := PhiAssignments(b, term.Target)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t%s\n", phis)
				}
				fmt.Fprintf(out, "\tgoto %s\n", BlockName(term.Target))

			case *ir.TermCondBr:
				cond, err := FormatValue(term.Cond)
				if err != nil {
					return fmt.Errorf("error translating condition (%v): %v", term.Cond, err)
				}
				fmt.Fprintf(out, "\tif %s {\n", cond)
				phis, err := PhiAssignments(b, term.TargetTrue)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetTrue))
				fmt.Fprintln(out, "\t} else {")
				phis, err = PhiAssignments(b, term.TargetFalse)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetFalse))
				fmt.Fprintln(out, "\t}")

			case *ir.TermRet:
				if term.X == nil {
					// void return
					if i == len(f.Blocks)-1 {
						// Just skip the return statement, since it's the end of the function anyway.
						continue
					}
					fmt.Fprintln(out, "\treturn")
				}
				retVal, err := FormatValue(term.X)
				if err != nil {
					return fmt.Errorf("error translating return value (%v): %v", term.X, err)
				}
				if f.Name() == "main" {
					fmt.Fprintf(out, "\tos.Exit(int(%s))\n", retVal)
				} else {
					fmt.Fprintf(out, "\treturn %s\n", retVal)
				}

			case *ir.TermSwitch:
				x, err := FormatValue(term.X)
				if err != nil {
					return fmt.Errorf("error translating control value (%v): %v", term.X, err)
				}
				fmt.Fprintf(out, "\tswitch %s {\n", x)
				for _, c := range term.Cases {
					x, err := FormatValue(c.X)
					if err != nil {
						return fmt.Errorf("error translating case value (%v): %v", c.X, err)
					}
					fmt.Fprintf(out, "\tcase %s:\n", x)
					phis, err := PhiAssignments(b, c.Target)
					if err != nil {
						return fmt.Errorf("error translating phi nodes: %v", err)
					}
					if phis != "" {
						fmt.Fprintf(out, "\t\t%s\n", phis)
					}
					fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(c.Target))
				}
				fmt.Fprint(out, "\tdefault:\n")
				phis, err := PhiAssignments(b, term.TargetDefault)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetDefault))
				fmt.Fprint(out, "\t}\n")

			default:
				return fmt.Errorf("unsupported block terminator type: %T", term)
			}
		}

		fmt.Fprint(out, "}\n\n")
	}
	return nil
}

// PhiAssignments returns an assignment statement expressing the effects of Phi
// nodes on the branch from block a to block b. If block b has no phi nodes,
// it returns the empty string.
func PhiAssignments(a, b value.Value) (string, error) {
	var dest, src []string
	for _, inst := range b.(*ir.Block).Insts {
		phi, ok := inst.(*ir.InstPhi)
		if !ok {
			break
		}
		for _, inc := range phi.Incs {
			if inc.Pred == a {
				source, err := FormatValue(inc.X)
				if err != nil {
					return "", fmt.Errorf("error translating value (%v): %v", inc.X, err)
				}
				src = append(src, source)
				dest = append(dest, VariableName(phi))
				break
			}
		}
	}
	if len(src) == 0 {
		return "", nil
	}
	return strings.Join(dest, ", ") + " = " + strings.Join(src, ", "), nil
}
