// compiler/symbol_table_test.go

package compiler

import "testing"

// TestDefine is a test case for Define
func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
		"c": Symbol{Name: "c", Scope: LocalScope, Index: 0},
		"d": Symbol{Name: "d", Scope: LocalScope, Index: 1},
		"e": Symbol{Name: "e", Scope: LocalScope, Index: 0},
		"f": Symbol{Name: "f", Scope: LocalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	firstLocal := NewEnclosedSymbolTable(global)

	c := firstLocal.Define("c")
	if c != expected["c"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	d := firstLocal.Define("d")
	if d != expected["d"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	secondLocal := NewEnclosedSymbolTable(global)

	e := secondLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	f := secondLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
}

// TestResolveNestedLocal is a test case for nested locals
func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}

}

// TestResolveLocal is a test case for locals
func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		Symbol{Name: "b", Scope: GlobalScope, Index: 1},
		Symbol{Name: "c", Scope: LocalScope, Index: 0},
		Symbol{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}
		if result != sym {
			t.Errorf("expected %s to resolve %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveGlobal(t *testing.T) {
	s := NewSymbolTable()

	s.Define("a")
	s.Define("b")

	expected := []Symbol{
		{Name: "a", Index: 0, Scope: GlobalScope},
		{Name: "b", Index: 1, Scope: GlobalScope},
	}

	for i, sym := range expected {
		sym, ok := s.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %q not resolvable", sym.Name)
		}
		if sym != expected[i] {
			t.Errorf("expected=%+v, got=%+v", expected[i], sym)
		}
	}
}
