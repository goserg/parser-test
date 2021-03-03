package parser

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	left := strings.Split(query, " ")[0]
	right := strings.Split(query, " ")[2]
	right = strings.Trim(right, `'"`)
	db = db.Where(squirrel.Eq{left: right})
	return &db, nil
}
