package parser

import "github.com/Masterminds/squirrel"

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	db = db.Where(squirrel.Eq{"Foo.Bar.X": "hello"})
	return &db, nil
}
