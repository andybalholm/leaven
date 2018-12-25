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

	fmt.Fprint(out, `package main

import "unsafe"
var _ unsafe.Pointer

`)

	for i, t := range m.TypeDefs {
		name := t.Name()
		if name == "" {
			name = fmt.Sprintf("type%d", i)
		}
		name = strings.TrimPrefix(name, "struct.")
		t.SetName(name)

		def, err := TypeDefinition(t)
		if err != nil {
			log.Fatalf("Error generating type definition for %v: %v", t, err)
		}

		fmt.Fprintf(out, "type %s %s\n\n", name, def)
	}

	// TODO: Globals

	for _, f := range m.Funcs {
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

		// Declare variables.
		vars := make(map[string][]string)
		for _, b := range f.Blocks {
			for _, inst := range b.Insts {
				if inst, ok := inst.(value.Named); ok {
					t, err := TypeSpec(inst.Type())
					if err != nil {
						log.Fatalf("Error translating type of %s in %s: %v", inst.Ident(), f.Name(), err)
					}
					vars[t] = append(vars[t], VariableName(inst))
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
				fmt.Fprintf(out, "\t%s\n", translated)
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
				// TODO: Assign values to phi nodes
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
				fmt.Fprintf(out, "\treturn %s\n", retVal)

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
