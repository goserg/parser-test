package parser

import (
	"errors"
	"testing"

	"github.com/Masterminds/squirrel"
)

type addTest struct {
	arg      string
	expected interface{}
	err      error
}

func TestParserEmpty(t *testing.T) {
	subQ := squirrel.Select("aa").From("bb")
	got, err := Parse("", subQ)
	if err != nil {
		t.Errorf(`Parse("", db) error = %v; want error = nil`, err)
	}
	gotSQL, gotArgs, _ := got.ToSql()
	wantSQL, wantArgs, _ := subQ.ToSql()
	if wantSQL != gotSQL {
		t.Errorf(`Parse("", db): got sql = "%v", want "%v"; got arg = %v, want %v`, gotSQL, wantSQL, gotArgs, wantArgs)
	}
}

func TestParserError(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X 'hello'", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AN Field2 != 7", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AND", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AND Field2 7", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' OR Field2 7", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' OR Field2 != 7 AND", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AN Field2 != 7 OR Field3 > 11.7", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AN Field2 != 7 AN Field3 > 11.7", nil, errors.New("Parsing error")},
		{"Field1 = 'foo' AN Field2 != 7 O Field3", nil, errors.New("Parsing error")},
		{"0Foo.Bar.X = 'hello'", nil, errors.New("Invalid column name")},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserEqual(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X = 'hello'", squirrel.Eq{"Foo.Bar.X": "hello"}, nil},
		{"Bar.Alpha = 7", squirrel.Eq{"Bar.Alpha": "7"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func runTest(t *testing.T, test addTest) {
	subQ := squirrel.Select("aa").From("bb")
	got, err := Parse(test.arg, subQ)
	if err != nil && err.Error() != test.err.Error() {
		t.Errorf(`Parse(%s, db) error = %v; want error = %v`, test.arg, err, test.err.Error())
	}
	if test.err == nil {
		gotSQL, gotArgs, _ := got.ToSql()
		want := subQ.Where(test.expected)
		wantSQL, wantArgs, _ := want.ToSql()
		if wantSQL != gotSQL || wantArgs[0] != gotArgs[0] {
			t.Errorf(`Parse(%s, db): got sql = "%v", want "%v"; got arg = %v, want %v`, test.arg, gotSQL, wantSQL, gotArgs, wantArgs)
		}
	}
}

func TestParserNotEqual(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X != 'hello'", squirrel.NotEq{"Foo.Bar.X": "hello"}, nil},
		{"Foo.Bar.X <> 'hello '", squirrel.NotEq{"Foo.Bar.X": "hello "}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserLike(t *testing.T) {
	runTest(t, addTest{"Alice.Name ~ 'A.'", squirrel.Like{"Alice.Name": "A."}, nil})
}

func TestParserNotLike(t *testing.T) {
	runTest(t, addTest{"Bob.LastName !~ 'Bill.'", squirrel.NotLike{"Bob.LastName": "Bill."}, nil})
}

func TestParserILike(t *testing.T) {
	runTest(t, addTest{"Alice.Name ~* 'A.'", squirrel.ILike{"Alice.Name": "A."}, nil})
}

func TestParserNotILike(t *testing.T) {
	runTest(t, addTest{"Alice.Name !~* 'A.'", squirrel.NotILike{"Alice.Name": "A."}, nil})
}

func TestParserGt(t *testing.T) {
	runTest(t, addTest{"Price > 1000", squirrel.Gt{"Price": "1000"}, nil})
}

func TestParserLt(t *testing.T) {
	runTest(t, addTest{"Price < 1000", squirrel.Lt{"Price": "1000"}, nil})
}

func TestParserGtOrEq(t *testing.T) {
	runTest(t, addTest{"Price >= 1000", squirrel.GtOrEq{"Price": "1000"}, nil})
}

func TestParserLtOrEq(t *testing.T) {
	runTest(t, addTest{"Price <= 1000", squirrel.LtOrEq{"Price": "1000"}, nil})
}

func TestParserAND(t *testing.T) {
	runTest(t, addTest{
		"Foo.Bar.Beta > 21 AND Alpha.Bar != 'hello'",
		squirrel.And{
			squirrel.Gt{"Foo.Bar.Beta": "21"},
			squirrel.NotEq{"Alpha.Bar": "hello"},
		},
		nil,
	})
}

func TestParserOR(t *testing.T) {
	runTest(t, addTest{
		"Foo.Bar.Beta > 21 OR Alpha.Bar != 'hello'",
		squirrel.Or{
			squirrel.Gt{"Foo.Bar.Beta": "21"},
			squirrel.NotEq{"Alpha.Bar": "hello"},
		},
		nil,
	})
}

func TestParserANDOR(t *testing.T) {
	runTest(t, addTest{
		"Field1 = 'foo' AND Field2 != 7 OR Field3 > 11.7",
		squirrel.Or{
			squirrel.And{
				squirrel.Eq{"Field1": "foo"},
				squirrel.NotEq{"Field2": "7"},
			},
			squirrel.Gt{"Field3": "11.7"},
		},
		nil,
	})
}

func TestParserORAND(t *testing.T) {
	runTest(t, addTest{
		"Field1 = 'foo' OR Field2 != 7 AND Field3 > 11.7",
		squirrel.Or{
			squirrel.Eq{"Field1": "foo"},
			squirrel.And{
				squirrel.NotEq{"Field2": "7"},
				squirrel.Gt{"Field3": "11.7"},
			},
		},
		nil,
	})
}

func TestParserANDAND(t *testing.T) {
	runTest(t, addTest{
		"Field1 = 'foo' AND Field2 != 7 AND Field3 > 11.7",
		squirrel.And{
			squirrel.And{
				squirrel.Eq{"Field1": "foo"},
				squirrel.NotEq{"Field2": "7"},
			},
			squirrel.Gt{"Field3": "11.7"},
		},
		nil,
	})
}

func TestParserOROR(t *testing.T) {
	runTest(t, addTest{
		"Field1 = 'foo' OR Field2 != 7 OR Field3 > 11.7",
		squirrel.Or{
			squirrel.Or{
				squirrel.Eq{"Field1": "foo"},
				squirrel.NotEq{"Field2": "7"},
			},
			squirrel.Gt{"Field3": "11.7"},
		},
		nil,
	})
}

func TestParseBool(t *testing.T) {
	runTest(t, addTest{
		"Alice.IsActive AND Bob.LastHash = 'ab5534b'",
		squirrel.And{
			squirrel.Eq{"Alice.IsActive": true},
			squirrel.Eq{"Bob.LastHash": "ab5534b"},
		},
		nil,
	})
}

func TestParseLongQuery(t *testing.T) {
	runTest(t, addTest{
		"a = 1 AND b = 2 AND c != 1 OR a = 2 AND b = 1 OR c = 3",
		squirrel.Or{
			squirrel.Or{
				squirrel.And{
					squirrel.And{
						squirrel.Eq{"a": "1"},
						squirrel.Eq{"b": "2"},
					},
					squirrel.NotEq{"c": "1"},
				},
				squirrel.And{
					squirrel.Eq{"a": "2"},
					squirrel.Eq{"b": "1"},
				},
			},
			squirrel.Eq{"c": "3"},
		},
		nil,
	})
}

func TestParseWithQuotesAndSpaceInColumnName(t *testing.T) {
	runTest(t, addTest{
		`"Hello world" = 'world'`,
		squirrel.Eq{`"Hello world"`: "world"},
		nil,
	})
}
