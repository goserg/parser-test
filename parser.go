package parser

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
)

func parseExpression(left string, middle string, right string) interface{} {
	left = strings.Trim(left, `"`)
	right = strings.Trim(right, `'`)
	var builder interface{}
	switch middle {
	case "=":
		builder = squirrel.Eq{left: right}
		break
	case "!=":
		builder = squirrel.NotEq{left: right}
		break
	case "~":
		builder = squirrel.Like{left: right}
		break
	case "!~":
		builder = squirrel.NotLike{left: right}
		break
	case "~*":
		builder = squirrel.ILike{left: right}
		break
	case "!~*":
		builder = squirrel.NotILike{left: right}
		break
	case "<>":
		builder = squirrel.NotEq{left: right}
		break
	case ">":
		builder = squirrel.Gt{left: right}
		break
	case ">=":
		builder = squirrel.GtOrEq{left: right}
		break
	case "<":
		builder = squirrel.Lt{left: right}
		break
	case "<=":
		builder = squirrel.LtOrEq{left: right}
		break
	}
	return builder
}

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	units := strings.Split(query, " ")
	if len(units) < 3 {
		return nil, errors.New("Parsing error")
	}
	if len(units) == 3 {
		expression := parseExpression(units[0], units[1], units[2])
		db = db.Where(expression)
	}
	return &db, nil
}
