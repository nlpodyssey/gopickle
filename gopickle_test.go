// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopickle

import (
	"github.com/nlpodyssey/gopickle/types"
	"math/big"
	"strings"
	"testing"
)

func TestNoneP1(t *testing.T) {
	// pickle.dumps(None, protocol=1)
	loadsNoErrEqual(t, "N.", nil)
}

func TestNoneP2(t *testing.T) {
	// pickle.dumps(None, protocol=2)
	loadsNoErrEqual(t, "\x80\x02N.", nil)
}

func TestTrueP1(t *testing.T) {
	// pickle.dumps(True, protocol=1)
	loadsNoErrEqual(t, "I01\n.", true)
}

func TestTrueP2(t *testing.T) {
	// pickle.dumps(True, protocol=2)
	loadsNoErrEqual(t, "\x80\x02\x88.", true)
}

func TestFalseP1(t *testing.T) {
	// pickle.dumps(False, protocol=1)
	loadsNoErrEqual(t, "I00\n.", false)
}

func TestFalseP2(t *testing.T) {
	// pickle.dumps(False, protocol=2)
	loadsNoErrEqual(t, "\x80\x02\x89.", false)
}

func TestIntP0Positive(t *testing.T) {
	// pickle.dumps(42, protocol=0)
	loadsNoErrEqual(t, "I42\n.", 42)
}

func TestIntP0Negative(t *testing.T) {
	// pickle.dumps(-42, protocol=0)
	loadsNoErrEqual(t, "I-42\n.", -42)
}

func TestFloatP0Positive(t *testing.T) {
	// pickle.dumps(4.2, protocol=0)
	loadsNoErrEqual(t, "F4.2\n.", 4.2)
}

func TestFloatP0Negative(t *testing.T) {
	// pickle.dumps(-4.2, protocol=0)
	loadsNoErrEqual(t, "F-4.2\n.", -4.2)
}

func TestBinIntP1Positive(t *testing.T) {
	// pickle.dumps(100200, protocol=1)
	loadsNoErrEqual(t, "Jh\x87\x01\x00.", 100200)
}

func TestBinIntP4Positive(t *testing.T) {
	// pickle.dumps(70100, protocol=4)
	loadsNoErrEqual(t,
		"\x80\x04\x95\x06\x00\x00\x00\x00\x00\x00\x00J\xd4\x11\x01\x00.", 70100)
}

func TestBinIntP1Negative(t *testing.T) {
	// pickle.dumps(-100200, protocol=1)
	loadsNoErrEqual(t, "J\x98x\xfe\xff.", -100200)
}

func TestBinIntP4Negative(t *testing.T) {
	// pickle.dumps(-70100, protocol=4)
	loadsNoErrEqual(t,
		"\x80\x04\x95\x06\x00\x00\x00\x00\x00\x00\x00J,\xee\xfe\xff.", -70100)
}

func TestBinInt1P2(t *testing.T) {
	// pickle.dumps(42, protocol=2)
	loadsNoErrEqual(t, "\x80\x02K*.", 42)
}

func TestBinInt2P2(t *testing.T) {
	// pickle.dumps(300, protocol=2)
	loadsNoErrEqual(t, "\x80\x02M,\x01.", 300)
}

func TestLongP1Positive(t *testing.T) {
	// pickle.dumps(100200300400, protocol=1)
	loadsNoErrEqual(t, "L100200300400L\n.", 100200300400)
}

func TestLongP1Negative(t *testing.T) {
	// pickle.dumps(-100200300400, protocol=1)
	loadsNoErrEqual(t, "L-100200300400L\n.", -100200300400)
}

func TestLongP1BigPositive(t *testing.T) {
	// pickle.dumps(100200300400500600700, protocol=1)
	actual := loadsNoErr(t, "L100200300400500600700L\n.")
	switch v := actual.(type) {
	case *big.Int:
		expected := "100200300400500600700"
		if v.String() != expected {
			t.Errorf("expected %s, actual %s", expected, v.String())
		}
	default:
		t.Error("expected big Int", actual)
	}
}

func TestLongP1BigNegative(t *testing.T) {
	// pickle.dumps(-100200300400500600700, protocol=1)
	actual := loadsNoErr(t, "L-100200300400500600700L\n.")
	switch v := actual.(type) {
	case *big.Int:
		expected := "-100200300400500600700"
		if v.String() != expected {
			t.Errorf("expected %s, actual %s", expected, v.String())
		}
	default:
		t.Error("expected big Int", actual)
	}
}

func TestStringPython27P0(t *testing.T) {
	// pickle.dumps('Café', protocol=0)  # Python 2.7
	// TODO: the string should be decoded
	loadsNoErrEqual(t, "S'Caf\\xc3\\xa9'\np0\n.", "Caf\\xc3\\xa9")
}

func TestBinStringPython27P1(t *testing.T) {
	// pickle.dumps(b'1234567890'*26, protocol=1)  # Python 2.7
	loadsNoErrEqual(t,
		"T\x04\x01\x00\x0012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890q\x00.",
		strings.Repeat("1234567890", 26))
}

func TestShortBinStringPython27P1(t *testing.T) {
	// pickle.dumps(b"Café", protocol=1)  # Python 2.7
	loadsNoErrEqual(t, "U\x05Caf\xc3\xa9q\x00.", "Café")
}

func TestUnicodePython27P0(t *testing.T) {
	// pickle.dumps(u"Café", protocol=0)  # Python 2.7
	loadsNoErrEqual(t, "VCaf\xe9\np0\n.", "Caf\xe9")
}

func TestBinUnicodeP1(t *testing.T) {
	// pickle.dumps('Café', protocol=1)
	loadsNoErrEqual(t, "X\x05\x00\x00\x00Caf\xc3\xa9q\x00.", "Café")
}

func TestShortBinUnicodeP4(t *testing.T) {
	// pickle.dumps('Café', protocol=4)
	loadsNoErrEqual(t,
		"\x80\x04\x95\t\x00\x00\x00\x00\x00\x00\x00\x8c\x05Caf\xc3\xa9\x94.",
		"Café")
}

func TestDictP0Empty(t *testing.T) {
	// pickle.dumps({}, protocol=0)
	actual := loadsNoErr(t, "(dp0\n.")
	switch v := actual.(type) {
	case *types.Dict:
		if v.Len() != 0 {
			t.Error("expected empty Dict, actual:", actual)
		}
	default:
		t.Error("expected Dict, actual:", actual)
	}
}

func TestDictP0OneKeyValue(t *testing.T) {
	// pickle.dumps({'a': 1}, protocol=0)
	actual := loadsNoErr(t, "(dp0\nVa\np1\nI1\ns.")
	switch v := actual.(type) {
	case *types.Dict:
		if x, ok := v.Get("a"); v.Len() != 1 || !ok || x != 1 {
			t.Error("expected {'a': 1}, actual:", actual)
		}
	default:
		t.Error("expected Dict, actual:", actual)
	}
}

func TestEmptyDictP2(t *testing.T) {
	// pickle.dumps({}, protocol=2)
	actual := loadsNoErr(t, "\x80\x02}q\x00.")
	switch v := actual.(type) {
	case *types.Dict:
		if v.Len() != 0 {
			t.Error("expected empty Dict, actual:", actual)
		}
	default:
		t.Error("expected Dict, actual:", actual)
	}
}

func TestTupleP0EmptyTuple(t *testing.T) {
	// pickle.dumps(tuple(), protocol=0)
	actual := loadsNoErr(t, "(t.")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 0 {
			t.Error("expected empty Tuple, actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestTupleP0OneItem(t *testing.T) {
	// pickle.dumps((1,), protocol=0)
	actual := loadsNoErr(t, "(I1\ntp0\n.")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 1 || v.Get(0) != 1 {
			t.Error("expected (1,), actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestEmptyTupleP2(t *testing.T) {
	// pickle.dumps(tuple(), protocol=2)
	actual := loadsNoErr(t, "\x80\x02).")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 0 {
			t.Error("expected empty Tuple, actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestTuple1P2(t *testing.T) {
	// pickle.dumps((1,), protocol=2)
	actual := loadsNoErr(t, "\x80\x02K\x01\x85q\x00.")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 1 || v.Get(0) != 1 {
			t.Error("expected (1,), actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestTuple2P2(t *testing.T) {
	// pickle.dumps((1, 2), protocol=2)
	actual := loadsNoErr(t, "\x80\x02K\x01K\x02\x86q\x00.")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 2 || v.Get(0) != 1 || v.Get(1) != 2 {
			t.Error("expected (1, 2), actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestTuple3P2(t *testing.T) {
	// pickle.dumps((1, 2, 3), protocol=2)
	actual := loadsNoErr(t, "\x80\x02K\x01K\x02K\x03\x87q\x00.")
	switch v := actual.(type) {
	case *types.Tuple:
		if v.Len() != 3 || v.Get(0) != 1 || v.Get(1) != 2 || v.Get(2) != 3 {
			t.Error("expected (1, 2, 3), actual:", actual)
		}
	default:
		t.Error("expected Tuple, actual:", actual)
	}
}

func TestListP0EmptyList(t *testing.T) {
	// pickle.dumps([], protocol=0)
	actual := loadsNoErr(t, "(lp0\n.")
	switch v := actual.(type) {
	case *types.List:
		if v.Len() != 0 {
			t.Error("expected empty List, actual:", actual)
		}
	default:
		t.Error("expected List, actual:", actual)
	}
}

func TestEmptyListP2(t *testing.T) {
	// pickle.dumps([], protocol=2)
	actual := loadsNoErr(t, "\x80\x02]q\x00.")
	switch v := actual.(type) {
	case *types.List:
		if v.Len() != 0 {
			t.Error("expected empty List, actual:", actual)
		}
	default:
		t.Error("expected List, actual:", actual)
	}
}

func TestListP2OneItem(t *testing.T) {
	// pickle.dumps([1], protocol=2)
	actual := loadsNoErr(t, "\x80\x02]q\x00K\x01a.")
	switch v := actual.(type) {
	case *types.List:
		if v.Len() != 1 || v.Get(0) != 1 {
			t.Error("expected [1], actual:", actual)
		}
	default:
		t.Error("expected List, actual:", actual)
	}
}

func TestListP2TwoItems(t *testing.T) {
	// pickle.dumps([1, 2], protocol=2)
	actual := loadsNoErr(t, "\x80\x02]q\x00(K\x01K\x02e.")
	switch v := actual.(type) {
	case *types.List:
		if v.Len() != 2 || v.Get(0) != 1 || v.Get(1) != 2 {
			t.Error("expected [1, 2], actual:", actual)
		}
	default:
		t.Error("expected List, actual:", actual)
	}
}

func TestBinFloatP2Positive(t *testing.T) {
	// pickle.dumps(1.2, protocol=2)
	loadsNoErrEqual(t, "\x80\x02G?\xf3333333.", 1.2)
}

func TestBinFloaP2tNegative(t *testing.T) {
	// pickle.dumps(-1.2, protocol=2)
	loadsNoErrEqual(t, "\x80\x02G\xbf\xf3333333.", -1.2)
}

func TestLong1P2SmallPositive(t *testing.T) {
	// pickle.dumps(100200300400, protocol=2)
	loadsNoErrEqual(t, "\x80\x02\x8a\x05p?gT\x17.", 100200300400)
}

func TestLong1P2SmallNegative(t *testing.T) {
	// pickle.dumps(-100200300400, protocol=2)
	loadsNoErrEqual(t, "\x80\x02\x8a\x05\x90\xc0\x98\xab\xe8.", -100200300400)
}

func TestLong1P2BigPositive(t *testing.T) {
	// pickle.dumps(100200300400500600700, protocol=2)
	actual := loadsNoErr(t, "\x80\x02\x8a\t|\xefD\x8fT\xfa\x8en\x05.")
	switch v := actual.(type) {
	case *big.Int:
		expected := "100200300400500600700"
		if v.String() != expected {
			t.Errorf("expected %s, actual %s", expected, v.String())
		}
	default:
		t.Error("expected big Int", actual)
	}
}

func TestLong1P2BigNegative(t *testing.T) {
	// pickle.dumps(-100200300400500600700, protocol=2)
	actual := loadsNoErr(t, "\x80\x02\x8a\t\x84\x10\xbbp\xab\x05q\x91\xfa.")
	switch v := actual.(type) {
	case *big.Int:
		expected := "-100200300400500600700"
		if v.String() != expected {
			t.Errorf("expected %s, actual %s", expected, v.String())
		}
	default:
		t.Error("expected big Int", actual)
	}
}

func TestBinBytesP3(t *testing.T) {
	// pickle.dumps(b'1234567890'*26, protocol=3)
	actual := loadsNoErr(t,
		"\x80\x03B\x04\x01\x00\x001234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"123456789012345678901234567890123456789012345678901234567890"+
			"1234567890123456789012345678901234567890q\x00.")
	switch v := actual.(type) {
	case []byte:
		expected := []byte(strings.Repeat("1234567890", 26))
		if string(v) != string(expected) {
			t.Errorf("expected %v actual: %v", expected, actual)
		}
	default:
		t.Error("expected []byte, actual:", actual)
	}
}

func TestShortBinBytesP3(t *testing.T) {
	// pickle.dumps(b'ab', protocol=3)
	actual := loadsNoErr(t, "\x80\x03C\x02abq\x00.")
	switch v := actual.(type) {
	case []byte:
		expected := []byte{'a', 'b'}
		if string(v) != string(expected) {
			t.Errorf("expected %v actual: %v", expected, actual)
		}
	default:
		t.Error("expected []byte, actual:", actual)
	}
}

func TestEmptySetP4(t *testing.T) {
	// pickle.dumps(set(), protocol=4)
	actual := loadsNoErr(t, "\x80\x04\x8f\x94.")
	switch v := actual.(type) {
	case *types.Set:
		if v.Len() != 0 {
			t.Error("expected empty Set, actual:", actual)
		}
	default:
		t.Error("expected Set, actual:", actual)
	}
}

func TestP4SetWithOneItem(t *testing.T) {
	// pickle.dumps(set([1]), protocol=4)
	actual := loadsNoErr(t,
		"\x80\x04\x95\x07\x00\x00\x00\x00\x00\x00\x00\x8f\x94(K\x01\x90.")
	switch v := actual.(type) {
	case *types.Set:
		if v.Len() != 1 || !v.Has(1) {
			t.Error("expected [1], actual:", actual)
		}
	default:
		t.Error("expected Set, actual:", actual)
	}
}

func TestFrozenSetP4EmptyFrozenSet(t *testing.T) {
	// pickle.dumps(frozenset(), protocol=4)
	actual := loadsNoErr(t,
		"\x80\x04\x95\x04\x00\x00\x00\x00\x00\x00\x00(\x91\x94.")
	switch v := actual.(type) {
	case *types.FrozenSet:
		if v.Len() != 0 {
			t.Error("expected empty FrozenSet, actual:", actual)
		}
	default:
		t.Error("expected FrozenSet, actual:", actual)
	}
}

func TestFrozenSetP4OneItem(t *testing.T) {
	// pickle.dumps(frozenset([1]), protocol=4)
	actual := loadsNoErr(t,
		"\x80\x04\x95\x06\x00\x00\x00\x00\x00\x00\x00(K\x01\x91\x94.")
	switch v := actual.(type) {
	case *types.FrozenSet:
		if v.Len() != 1 || !v.Has(1) {
			t.Error("expected [1], actual:", actual)
		}
	default:
		t.Error("expected FrozenSet, actual:", actual)
	}
}

func TestP0GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=0)
	actual := loadsNoErr(t, "ccopy_reg\n_reconstructor\np0\n(c__main__\nFoo\n"+
		"p1\nc__builtin__\nobject\np2\nNtp3\nRp4\n.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP1GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=1)
	actual := loadsNoErr(t, "ccopy_reg\n_reconstructor\nq\x00(c__main__\nFoo\n"+
		"q\x01c__builtin__\nobject\nq\x02Ntq\x03Rq\x04.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP2GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=2)
	actual := loadsNoErr(t, "\x80\x02c__main__\nFoo\nq\x00)\x81q\x01.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP3GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=3)
	actual := loadsNoErr(t, "\x80\x03c__main__\nFoo\nq\x00)\x81q\x01.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP4GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=4)
	actual := loadsNoErr(t, "\x80\x04\x95\x17\x00\x00\x00\x00\x00\x00\x00"+
		"\x8c\x08__main__\x94\x8c\x03Foo\x94\x93\x94)\x81\x94.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP5GenericObject(t *testing.T) {
	// class Foo(): pass
	// pickle.dumps(Foo(), protocol=5)
	actual := loadsNoErr(t, "\x80\x05\x95\x17\x00\x00\x00\x00\x00\x00\x00"+
		"\x8c\x08__main__\x94\x8c\x03Foo\x94\x93\x94)\x81\x94.")
	switch v := actual.(type) {
	case *types.GenericObject:
		if v.Class.Module != "__main__" || v.Class.Name != "Foo" ||
			len(v.ConstructorArgs) != 0 {
			t.Errorf("expected __main__.Foo(), actual: %#v", v)
		}
	default:
		t.Error("expected GenericObject, actual:", actual)
	}
}

func TestP4EmptyOrderedDict(t *testing.T) {
	// pickle.dumps(collections.OrderedDict(), protocol=4)
	actual := loadsNoErr(t, "\x80\x04\x95\"\x00\x00\x00\x00\x00\x00\x00"+
		"\x8c\x0bcollections\x94\x8c\x0bOrderedDict\x94\x93\x94)R\x94.")
	switch v := actual.(type) {
	case *types.OrderedDict:
		if v.Len() != 0 {
			t.Error("expected empty OrderedDict, actual:", actual)
		}
	default:
		t.Error("expected OrderedDict, actual:", actual)
	}
}

func TestP4OrderedDictWithOneKeyValue(t *testing.T) {
	// pickle.dumps(collections.OrderedDict({'a': 1}), protocol=4)
	actual := loadsNoErr(t, "\x80\x04\x95)\x00\x00\x00\x00\x00\x00\x00"+
		"\x8c\x0bcollections\x94\x8c\x0bOrderedDict\x94\x93\x94)R\x94"+
		"\x8c\x01a\x94K\x01s.")
	switch v := actual.(type) {
	case *types.OrderedDict:
		if x, ok := v.Get("a"); v.Len() != 1 || !ok || x != 1 {
			t.Error("expected {'a': 1}, actual:", actual)
		}
	default:
		t.Error("expected Dict, actual:", actual)
	}
}

func TestP4NestedDicts(t *testing.T) {
	// pickle.dumps({'a': 1, 'b': {'c': 2}}, protocol=4)
	actual := loadsNoErr(t, "\x80\x04\x95\x18\x00\x00\x00\x00\x00\x00\x00}"+
		"\x94(\x8c\x01a\x94K\x01\x8c\x01b\x94}\x94\x8c\x01c\x94K\x02su.")
	switch v := actual.(type) {
	case *types.Dict:
		if v.Len() != 2 {
			t.Error("expected two entries, actual:", actual)
		}
		if a, ok := v.Get("a"); !ok || a != 1 {
			t.Error("expected 'a' => 1, actual:", actual)
		}
		b, bOk := v.Get("b")
		bDict, bDictOk := b.(*types.Dict)
		if !bOk || !bDictOk {
			t.Error("expected 'b' => Dict, actual:", actual)
		}
		if c, ok := bDict.Get("c"); bDict.Len() != 1 || !ok || c != 2 {
			t.Error("expected 'c' => 2, actual:", actual)
		}
	default:
		t.Error("expected Dict, actual:", actual)
	}
}

// TODO: test BinPersId
// TODO: test Get
// TODO: test BinGet
// TODO: test LongBinPut
// TODO: test LongBinGet
// TODO: test Build
// TODO: test PersId
// TODO: test Pop
// TODO: test PopMark
// TODO: test Dup
// TODO: test Inst
// TODO: test Obj
// TODO: test Long4

func loadsNoErrEqual(t *testing.T, s string, expected interface{}) {
	actual := loadsNoErr(t, s)
	if actual != expected {
		t.Errorf("expected %v, actual: %v", expected, actual)
	}
}

func loadsNoErr(t *testing.T, s string) interface{} {
	result, err := Loads(s)
	if err != nil {
		t.Error(err)
	}
	return result
}
