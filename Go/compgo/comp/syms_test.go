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
	lsekai := "isekai"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)

	local := NewFrameSymbolTable(glob)
	local.Define("c")
	local.Define(lsekai)
	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: lsekai, Scope: LocalScope, Index: 1},
	}
	testResolve(t, expected, local)
}

func TestResolveLocal_nested(t *testing.T) {
	isekai := "異世界"
	lsekai := "isekai"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)

	local := NewFrameSymbolTable(glob)
	local.Define("c")
	local.Define(lsekai)

	local2 := NewFrameSymbolTable(glob)
	local2.Define("d")
	local2.Define("e")
	expectedL1 := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: lsekai, Scope: LocalScope, Index: 1},
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

func TestResolveFree(t *testing.T) {
	isekai := "異世界"
	lsekai := "isekai"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)

	local := NewFrameSymbolTable(glob)
	local.Define("c")
	local.Define(lsekai)

	local2 := NewFrameSymbolTable(local)
	local2.Define("d")
	local2.Define("e")

	tests := []struct {
		table                         *SymbolTable
		expectedSymbols, expectedFree []Symbol
	}{
		{
			local,
			[]Symbol{
				{"a", GlobalScope, 0},
				{isekai, GlobalScope, 1},
				{"c", LocalScope, 0},
				{lsekai, LocalScope, 1},
			},
			[]Symbol{},
		},
		{
			local2,
			[]Symbol{
				{"a", GlobalScope, 0},
				{isekai, GlobalScope, 1},
				{"c", FreeScope, 0},
				{lsekai, FreeScope, 1},
				{"d", LocalScope, 0},
				{"e", LocalScope, 1},
			},
			[]Symbol{
				{"c", LocalScope, 0},
				{lsekai, LocalScope, 1},
			},
		},
	}
	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			r, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s is not resolvable", sym.Name)
				continue
			}
			if r != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, r)
			}
		}
		if len(tt.table.FreeSymbols) != len(tt.expectedFree) {
			t.Errorf("wrong number of free symbols. got=%d want=%d",
				len(tt.table.FreeSymbols), len(tt.expectedFree))
			continue
		}

		for i, sym := range tt.expectedFree {
			r := tt.table.FreeSymbols[i]
			if r != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v", r, sym)
			}
		}
	}
}
