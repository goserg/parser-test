package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
)

func TestParser(t *testing.T) {
	subQ := squirrel.Select("aa").From("bb")
	got, err := Parse("Foo.Bar.X = 'hello'", subQ)
	gotSQL, gotArgs, _ := got.ToSql()
	if err != nil {
		t.Errorf(`Parse("Foo.Bar.X = 'hello'", db) error = %d; want error = nil`, err)
	}

	want := subQ.Where(squirrel.Eq{"Foo.Bar.X": "helo"})
	wantSQL, wantArgs, _ := want.ToSql()
	if wantSQL != gotSQL || wantArgs[0] != gotArgs[0] {
		t.Errorf(`Parse("Foo.Bar.X = 'hello'", db): got sql = "%s", want "%s"; got arg = "%s", want "%s"`, gotSQL, wantSQL, gotArgs, wantArgs)
	}
}
