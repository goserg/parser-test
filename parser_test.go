package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
)

type addTest struct {
	arg      string
	expected interface{}
}

func TestParsers(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X = 'hello'", squirrel.Eq{"Foo.Bar.X": "hello"}},
		{"Bar.Alpha = 7", squirrel.Eq{"Bar.Alpha": "7"}},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func runTest(t *testing.T, test addTest) {
	subQ := squirrel.Select("aa").From("bb")
	got, err := Parse(test.arg, subQ)
	gotSQL, gotArgs, _ := got.ToSql()
	if err != nil {
		t.Errorf(`Parse(%s, db) error = %v; want error = nil`, test.arg, err)
	}

	want := subQ.Where(test.expected)
	wantSQL, wantArgs, _ := want.ToSql()
	if wantSQL != gotSQL || wantArgs[0] != gotArgs[0] {
		t.Errorf(`Parse(%s, db): got sql = "%v", want "%v"; got arg = %v, want %v`, test.arg, gotSQL, wantSQL, gotArgs[0], wantArgs[0])
	}
}

func TestParserNotEqual(t *testing.T) {
	var addTests = []addTest{
		{"Foo.Bar.X != 'hello'", squirrel.NotEq{"Foo.Bar.X": "hello"}},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserLike(t *testing.T) {
	var addTests = []addTest{
		{"Alice.Name ~ 'A.'", squirrel.Like{"Alice.Name": "A."}},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}

func TestParserNotLike(t *testing.T) {
	var addTests = []addTest{
		{"Bob.LastName !~ 'Bill.'", squirrel.NotLike{"Bob.LastName": "Bill."}},
	}
	for _, test := range addTests {
		runTest(t, test)
	}
}
