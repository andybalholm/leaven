package main

import (
	"fmt"
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

	fmt.Fprint(out, "package main\n\n")

	for i, t := range m.TypeDefs {
		name := t.Name()
		if name == "union.anon" {
			continue
		}
		if name == "" {
			name = fmt.Sprintf("type%d", i)
		}
		name = strings.TrimPrefix(name, "struct.")
		name = strings.TrimPrefix(name, "union.")
		t.SetName(name)

		def, err := TypeDefinition(t)
		if err != nil {
			log.Fatalf("Error generating type definition for %v: %v", t, err)
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
			log.Fatalf("Error translating type (%v): %v", g.ContentType, err)
		}
		val, err := FormatValue(g.Init)
		if err != nil {
			log.Fatalf("Error translating initializer (%v): %v", g.Init, err)
		}
		fmt.Fprintf(out, "var %s %s = %s\n\n", VariableName(g), t, val)
	}

	for _, f := range m.Funcs {
		if f.Blocks == nil {
			// Just a declaration, not a definition; skip it.
			continue
		}
		if f.Name() == "main" {
			fmt.Fprintln(out, "func main() {")
		} else {
			fmt.Fprintf(out, "func %s(", f.Name())
			for i, p := range f.Params {
				if i > 0 {
					fmt.Fprint(out, ", ")
				}
				pt, err := TypeSpec(p.Typ)
				if err != nil {
					log.Fatalf("Error translating type for parameter %d of %s: %v", i, f.Name(), err)
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
					log.Fatalf("Error translating return type for %s: %v", f.Name(), err)
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
						log.Fatalf("Error translating type of %s in %s: %v", inst.Ident(), f.Name(), err)
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
				fmt.Fprintf(out, "\nblock%d:\n", b.LocalID)
			}
			for _, inst := range b.Insts {
				if _, ok := inst.(*ir.InstPhi); ok {
					continue
				}
				translated, err := TranslateInstruction(inst)
				if err != nil {
					log.Fatalf("Error translating %q: %v", inst.Def(), err)
				}
				if translated != "" {
					fmt.Fprintf(out, "\t%s\n", translated)
				}
			}
			switch term := b.Term.(type) {
			case *ir.TermBr:
				phis, err := PhiAssignments(b, term.Target)
				if err != nil {
					log.Fatalf("Error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t%s\n", phis)
				}
				fmt.Fprintf(out, "\tgoto block%d\n", term.Target.LocalID)

			case *ir.TermCondBr:
				cond, err := FormatValue(term.Cond)
				if err != nil {
					log.Fatalf("Error translating condition (%v): %v", term.Cond, err)
				}
				fmt.Fprintf(out, "\tif %s {\n", cond)
				phis, err := PhiAssignments(b, term.TargetTrue)
				if err != nil {
					log.Fatalf("Error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto block%d\n", term.TargetTrue.LocalID)
				fmt.Fprintln(out, "\t} else {")
				phis, err = PhiAssignments(b, term.TargetFalse)
				if err != nil {
					log.Fatalf("Error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto block%d\n", term.TargetFalse.LocalID)
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
					log.Fatalf("Error translating return value (%v): %v", term.X, err)
				}
				if f.Name() == "main" {
					fmt.Fprintf(out, "\tos.Exit(int(%s))\n", retVal)
				} else {
					fmt.Fprintf(out, "\treturn %s\n", retVal)
				}

			case *ir.TermSwitch:
				x, err := FormatValue(term.X)
				if err != nil {
					log.Fatalf("Error translating control value (%v): %v", term.X, err)
				}
				fmt.Fprintf(out, "\tswitch %s {\n", x)
				for _, c := range term.Cases {
					x, err := FormatValue(c.X)
					if err != nil {
						log.Fatalf("Error translating case value (%v): %v", c.X, err)
					}
					fmt.Fprintf(out, "\tcase %s:\n", x)
					phis, err := PhiAssignments(b, c.Target)
					if err != nil {
						log.Fatalf("Error translating phi nodes: %v", err)
					}
					if phis != "" {
						fmt.Fprintf(out, "\t\t%s\n", phis)
					}
					fmt.Fprintf(out, "\t\tgoto block%d\n", c.Target.LocalID)
				}
				fmt.Fprint(out, "\tdefault:\n")
				phis, err := PhiAssignments(b, term.TargetDefault)
				if err != nil {
					log.Fatalf("Error translating phi nodes: %v", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto block%d\n", term.TargetDefault.LocalID)
				fmt.Fprint(out, "\t}\n")

			default:
				log.Fatalf("Unsupported block terminator type: %T", term)
			}
		}

		fmt.Fprint(out, "}\n\n")
	}
}

// PhiAssignments returns an assignment statement expressing the effects of Phi
// nodes on the branch from block a to block b. If block b has no phi nodes,
// it returns the empty string.
func PhiAssignments(a, b *ir.BasicBlock) (string, error) {
	var dest, src []string
	for _, inst := range b.Insts {
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
