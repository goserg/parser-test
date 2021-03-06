package parser

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"

	"github.com/goserg/parser-test/units"
)

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	var expression, nextExpression squirrel.Sqlizer
	var sign units.SignUnit
	var err error

	for query != "" {
		if expression == nil {
			expression, query, err = getNextANDBlock(query)
			if err != nil {
				return nil, err
			}
		}
		if query == "" {
			break
		}
		sign.Value, query = getNextUnit(query)
		if sign.IsABoolSign() {
			if sign.IsORSign() {
				nextExpression, query, err = getNextANDBlock(query)
				if err != nil {
					return nil, err
				}
				expression = squirrel.Or{expression, nextExpression}
			}
		}
	}
	db = db.Where(expression)
	return &db, nil
}

func parseExpression(column units.ColumnUnit, sign units.SignUnit, value string) squirrel.Sqlizer {
	value = strings.Trim(value, `'`)
	switch sign.Value {
	case "=":
		return squirrel.Eq{column.Name: value}
	case "!=":
		return squirrel.NotEq{column.Name: value}
	case "~":
		return squirrel.Like{column.Name: value}
	case "!~":
		return squirrel.NotLike{column.Name: value}
	case "~*":
		return squirrel.ILike{column.Name: value}
	case "!~*":
		return squirrel.NotILike{column.Name: value}
	case "<>":
		return squirrel.NotEq{column.Name: value}
	case ">":
		return squirrel.Gt{column.Name: value}
	case ">=":
		return squirrel.GtOrEq{column.Name: value}
	case "<":
		return squirrel.Lt{column.Name: value}
	default: // case "<="
		return squirrel.LtOrEq{column.Name: value}
	}
}

func getNextUnit(query string) (string, string) {
	delimiter := " "
	if strings.HasPrefix(query, `"`) {
		delimiter = `"`
	} else if strings.HasPrefix(query, `'`) {
		delimiter = `'`
	}
	query = strings.TrimLeft(query, delimiter)
	indexOfFirstDelimiter := strings.Index(query, delimiter)
	if indexOfFirstDelimiter == -1 {
		return query, ""
	}
	left := query[:indexOfFirstDelimiter]
	rest := query[indexOfFirstDelimiter+1:]
	if delimiter == `"` {
		left = `"` + left + `"`
	}
	return left, rest
}

func getNextExpr(query string) (squirrel.Sqlizer, string, error) {
	var column units.ColumnUnit
	var sign units.SignUnit
	var value string

	column.Name, query = getNextUnit(query)
	if !column.IsValid() {
		return nil, "", errors.New("Invalid column name")
	}

	sign.Value, _ = getNextUnit(query)
	if sign.IsAComparisonSign() {
		_, query = getNextUnit(query)
		value, query = getNextUnit(query)
		return parseExpression(column, sign, value), query, nil
	}
	if sign.IsABoolSign() {
		return squirrel.Eq{column.Name: "true"}, query, nil
	}
	return nil, query, errors.New("Parsing error")
}

func getNextANDBlock(query string) (squirrel.Sqlizer, string, error) {
	var expr, nextExpr squirrel.Sqlizer
	var err error
	var sign units.SignUnit

	expr, query, err = getNextExpr(query)
	if err != nil {
		return nil, "", err
	}
	sign.Value, _ = getNextUnit(query)
	for sign.IsANDSign() {
		_, query = getNextUnit(query)
		if query == "" {
			return nil, "", err
		}
		nextExpr, query, err = getNextExpr(query)
		if err != nil {
			return nil, "", err
		}
		expr = squirrel.And{expr, nextExpr}
		sign.Value, _ = getNextUnit(query)
	}
	return expr, query, nil
}
