package main

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"golang.org/x/exp/slices"
)

// An Index has maps for navigating "against the grain" of the pointers in the
// intermediate representation: from values to the instructions that use them,
// and from instructions to the blocks that contain them.
type Index struct {
	users  map[value.Value][]value.User
	blocks map[value.User]*ir.Block
}

// Add adds the instructions in f to the index.
func (i *Index) Add(f *ir.Func) {
	if i.users == nil {
		i.users = make(map[value.Value][]value.User)
	}
	if i.blocks == nil {
		i.blocks = make(map[value.User]*ir.Block)
	}

	for _, b := range f.Blocks {
		for _, inst := range b.Insts {
			for _, v := range inst.Operands() {
				i.users[*v] = append(i.users[*v], inst)
			}
			i.blocks[inst] = b
		}
		for _, v := range b.Term.Operands() {
			i.users[*v] = append(i.users[*v], b.Term)
		}
		i.blocks[b.Term] = b
	}
}

// Users returns a slice of the instructions and terminators that use v.
func (i *Index) Users(v value.Value) []value.User {
	return i.users[v]
}

// ReplaceValue replaces oldVal with newVal wherever it is used.
func (i *Index) ReplaceValue(oldVal, newVal value.Value) {
	for _, user := range i.users[oldVal] {
		for _, slot := range user.Operands() {
			if *slot == oldVal {
				*slot = newVal
				i.users[newVal] = append(i.users[newVal], user)
				break
			}
		}
	}
	delete(i.users, oldVal)
}

// deleteElement deletes the first element in s that is equal to v, and returns
// the modified slice.
func deleteElement[E comparable](s []E, v E) []E {
	i := slices.Index(s, v)
	if i == -1 {
		return s
	}
	return slices.Delete(s, i, i+1)
}

// DeleteInstruction deletes inst.
func (i *Index) DeleteInstruction(inst ir.Instruction) {
	b := i.blocks[inst]
	if b == nil {
		return
	}
	b.Insts = deleteElement(b.Insts, inst)
	for _, op := range inst.Operands() {
		v := *op
		i.users[v] = deleteElement[value.User](i.users[v], inst)
	}
}
