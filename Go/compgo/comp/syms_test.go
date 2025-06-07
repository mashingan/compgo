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

func TestResolveGlobal(t *testing.T) {
	isekai := "異世界"
	glob := NewSymbolTable()
	glob.Define("a")
	glob.Define(isekai)
	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: isekai, Scope: GlobalScope, Index: 1},
	}
	for _, sym := range expected {
		res, ok := glob.Resolve(sym.Name)
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
