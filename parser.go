package parser

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
)

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	left := strings.Split(query, " ")[0]
	right := strings.Split(query, " ")[2]
	right = strings.Trim(right, `'"`)

	middle := strings.Split(query, " ")[1]

	switch middle {
	case "=":
		db = db.Where(squirrel.Eq{left: right})
		break
	case "~":
		db = db.Where(squirrel.Like{left: right})
		break
	case "!=":
		db = db.Where(squirrel.NotEq{left: right})
		break
	case "!~":
		db = db.Where(squirrel.NotLike{left: right})
		break
	default:
		return nil, errors.New("can't parse it")
	}
	return &db, nil
}
