package parser

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
)

func parseExpression(left string, middle string, right string) squirrel.Sqlizer {
	right = strings.Trim(right, `'`)
	var builder squirrel.Sqlizer
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
	var sign, left, right string

	left, query = getNextUnit(query)
	if !isColumnNameValid(left) {
		return nil, "", errors.New("Invalid column name")
	}

	sign, _ = getNextUnit(query)
	if isAComparisonSign(sign) {
		_, query = getNextUnit(query)
		right, query = getNextUnit(query)
		return parseExpression(left, sign, right), query, nil
	}
	if isABooleanSign(sign) {
		return squirrel.Eq{left: "true"}, query, nil
	}
	return nil, query, errors.New("Parsing error")
}

func isColumnNameValid(name string) bool {
	if strings.Index("0123456789", string(name[0])) != -1 {
		return false
	}

	return true
}

func getNextANDBlock(query string) (squirrel.Sqlizer, string, error) {
	var expr, nextExpr squirrel.Sqlizer
	var err error
	expr, query, err = getNextExpr(query)
	if err != nil {
		return nil, "", err
	}
	sign, _ := getNextUnit(query)
	for sign == "AND" || sign == "and" {
		_, query = getNextUnit(query)
		if query == "" {
			return nil, "", err
		}
		nextExpr, query, err = getNextExpr(query)
		if err != nil {
			return nil, "", err
		}
		expr = squirrel.And{expr, nextExpr}
		sign, _ = getNextUnit(query)
	}
	return expr, query, nil
}

func isAComparisonSign(elem string) bool {
	for _, comp := range []string{"=", "!=", "<>", "~", "!~", "~*", "!~*", ">", ">=", "<", "<="} {
		if comp == elem {
			return true
		}
	}
	return false
}

func isABooleanSign(elem string) bool {
	for _, comp := range []string{"AND", "OR", "and", "or"} {
		if comp == elem {
			return true
		}
	}
	return false
}

//Parse - простой парсер WHERE части запросов SQL
func Parse(query string, db squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
	var expression, rightExpr, empty squirrel.Sqlizer
	var sign string
	var err error
	for query != "" {
		if expression == empty {
			expression, query, err = getNextANDBlock(query)
			if err != nil {
				return nil, err
			}
		}
		if query == "" {
			break
		}
		sign, query = getNextUnit(query)
		if isABooleanSign(sign) {
			if sign == "OR" || sign == "or" {
				rightExpr, query, err = getNextANDBlock(query)
				if err != nil {
					return nil, err
				}
				expression = squirrel.Or{expression, rightExpr}
			}
		}
	}
	db = db.Where(expression)
	return &db, nil
}

func isColAndValueValid(colName string, value interface{}) {

}
