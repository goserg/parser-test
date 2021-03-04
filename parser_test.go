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

func TestParserError(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X 'hello'", nil, errors.New("Parsing error")},
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
			t.Errorf(`Parse(%s, db): got sql = "%v", want "%v"; got arg = %v, want %v`, test.arg, gotSQL, wantSQL, gotArgs[0], wantArgs[0])
		}
	}
}

func TestParserNotEqual(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X != 'hello'", squirrel.NotEq{"Foo.Bar.X": "hello"}, nil},
		{"Foo.Bar.X <> 'hello'", squirrel.NotEq{"Foo.Bar.X": "hello"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserLike(t *testing.T) {
	var addTests = []addTest{
		{"Alice.Name ~ 'A.'", squirrel.Like{"Alice.Name": "A."}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserNotLike(t *testing.T) {
	var addTests = []addTest{
		{"Bob.LastName !~ 'Bill.'", squirrel.NotLike{"Bob.LastName": "Bill."}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserILike(t *testing.T) {
	var addTests = []addTest{
		{"Alice.Name ~* 'A.'", squirrel.ILike{"Alice.Name": "A."}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserNotILike(t *testing.T) {
	var addTests = []addTest{
		{"Alice.Name !~* 'A.'", squirrel.NotILike{"Alice.Name": "A."}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserGt(t *testing.T) {
	var addTests = []addTest{
		{"Price > 1000", squirrel.Gt{"Price": "1000"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserLt(t *testing.T) {
	var addTests = []addTest{
		{"Price < 1000", squirrel.Lt{"Price": "1000"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserGtOrEq(t *testing.T) {
	var addTests = []addTest{
		{"Price >= 1000", squirrel.GtOrEq{"Price": "1000"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserLtOrEq(t *testing.T) {
	var addTests = []addTest{
		{"Price <= 1000", squirrel.LtOrEq{"Price": "1000"}, nil},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}
