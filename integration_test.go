package main

import (
	"bytes"
	"os/exec"
	"testing"
)

// doTestCase runs the specified program in the testdata directory twice,
// once compiled directly with clang, and the other time compiled to Go with
// leaven. It compares the output of the two programs.
func doTestCase(t *testing.T, progName string) {
	p := "testdata/" + progName
	clang := exec.Command("clang", "-o", p+"_c", p+".c", "testdata/main.c")
	if err := clang.Run(); err != nil {
		t.Fatalf("Error in native compilation: %v", err)
	}

	prog := exec.Command(p + "_c")
	nativeOut, err := prog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running natively-compiled program: %v", err)
	}

	clang2 := exec.Command("clang", "-S", "-emit-llvm", "-o", p+".ll", p+".c")
	if err := clang2.Run(); err != nil {
		t.Fatalf("Error compiling to LLVM: %v", err)
	}

	leaven := exec.Command("leaven", p+".ll")
	if err := leaven.Run(); err != nil {
		t.Fatalf("Error running leaven: %v", err)
	}

	goimports := exec.Command("goimports", "-w", p+".go")
	if err := goimports.Run(); err != nil {
		t.Fatalf("Error running goimports: %v", err)
	}

	goRun := exec.Command("go", "run", p+".go", "testdata/main.go")
	goOut, err := goRun.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running Go program: %v", err)
	}

	if !bytes.Equal(goOut, nativeOut) {
		t.Fatalf("Output does not match. C = %q, Go = %q", nativeOut, goOut)
	}
}

func TestHelloWorld(t *testing.T) {
	doTestCase(t, "hello")
}

func TestHelloWorldPuts(t *testing.T) {
	doTestCase(t, "hello-puts")
}

func TestBinaryTrees(t *testing.T) {
	doTestCase(t, "binarytrees")
}

func TestFannkuch(t *testing.T) {
	doTestCase(t, "fannkuch-redux")
}
