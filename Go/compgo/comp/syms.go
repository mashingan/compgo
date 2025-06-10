package comp

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store       map[string]Symbol
	numdef      int
	scoped      *SymbolTable
	FreeSymbols []Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: map[string]Symbol{}}
}

func (s *SymbolTable) Define(sym string) Symbol {
	syms := Symbol{sym, GlobalScope, s.numdef}
	if s.scoped != nil {
		syms.Scope = LocalScope
	}
	s.numdef++
	s.store[sym] = syms
	return syms
}

func (s *SymbolTable) DefineBuiltin(index int, sym string) Symbol {
	ss := Symbol{sym, BuiltinScope, index}
	s.store[sym] = ss
	return ss
}

func (s *SymbolTable) ResolveBuiltin(sym string) (Symbol, bool) {
	ss, ok := s.store[sym]
	if !ok {
		return ss, ok
	}
	if ss.Scope != BuiltinScope {
		return ss, false
	}
	return ss, true
}

func (s *SymbolTable) Resolve(sym string) (Symbol, bool) {
	var (
		ss           Symbol
		ok, captured bool
	)
	st := s
	for !ok && st != nil {
		ss, ok = st.store[sym]
		if !ok {
			st = st.scoped
		} else if st != s {
			captured = true
		}
	}
	if captured && ss.Name != "" && !(ss.Scope == GlobalScope || ss.Scope == BuiltinScope) {
		idx := len(s.FreeSymbols)
		s.FreeSymbols = append(s.FreeSymbols, ss)
		ss = Symbol{ss.Name, FreeScope, idx}
	}
	return ss, ok
}

func NewFrameSymbolTable(outer *SymbolTable) *SymbolTable {
	st := NewSymbolTable()
	st.scoped = outer
	return st
}
