package comp

import "testing"

func TestDefine(t *testing.T) {
	isekai := "異世界"
	expected := map[string]Symbol{
		"a":    {Name: "a", Scope: GlobalScope, Index: 0},
		isekai: {Name: isekai, Scope: GlobalScope, Index: 1},
	}
	global := NewSymbolTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	異世界 := global.Define(isekai)
	if 異世界 != expected[isekai] {
		t.Errorf("expected %s=%+v, got=%+v", isekai, expected[isekai], 異世界)
	}
}

func testResolve(t *testing.T, expected []Symbol, table *SymbolTable) {
	for _, sym := range expected {
		res, ok := table.Resolve(sym.Name)
		if !ok {
			t.Errorf("sym.Name %q not resolvable", sym.Name)
			continue
		}
		if res != sym {
			t.Errorf("expected %s to resolve to %+v, got=%v",
				sym.Name, sym, res)
		}
	}
}

func TestResolveGlobal(t *testing.T) {
	isekai := "異世界"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)
	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
	}
	testResolve(t, expected, glob)
}

func TestResolveLocal(t *testing.T) {
	isekai := "異世界"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)

	local := NewFrameSymbolTable(glob)
	local.Define("c")
	local.Define("isekai")
	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: isekai, Scope: LocalScope, Index: 1},
	}
	testResolve(t, expected, local)
}

func TestResolveLocal_nested(t *testing.T) {
	isekai := "異世界"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)

	local := NewFrameSymbolTable(glob)
	local.Define("c")
	local.Define("isekai")

	local2 := NewFrameSymbolTable(glob)
	local2.Define("d")
	local2.Define("e")
	expectedL1 := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: isekai, Scope: LocalScope, Index: 1},
	}
	expectedL2 := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
		{Name: "d", Scope: LocalScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 1},
	}
	testResolve(t, expectedL1, local)
	testResolve(t, expectedL2, local2)
}
