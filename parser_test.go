package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
)

type addTest struct {
	arg      string
	expected interface{}
}

var addTests = []addTest{
	{"Foo.Bar.X = 'hello'", squirrel.Eq{"Foo.Bar.X": "hello"}},
	{"Bar.Alpha = 7", squirrel.Eq{"Bar.Alpha": "7"}},
}

func TestParsers(t *testing.T) {
	for _, test := range addTests {
		subQ := squirrel.Select("aa").From("bb")
		got, err := Parse(test.arg, subQ)
		gotSQL, gotArgs, _ := got.ToSql()
		if err != nil {
			t.Errorf(`Parse(%s, db) error = %d; want error = nil`, test.arg, err)
		}

		want := subQ.Where(test.expected)
		wantSQL, wantArgs, _ := want.ToSql()
		if wantSQL != gotSQL || wantArgs[0] != gotArgs[0] {
			t.Errorf(`Parse(%s, db): got sql = "%s", want "%s"; got arg = %d, want %d`, test.arg, gotSQL, wantSQL, gotArgs[0], wantArgs[0])
		}
	}
}
