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
	query = strings.Trim(query, " ")
	indexOfFirstSpace := strings.Index(query, " ")
	if indexOfFirstSpace == -1 {
		return query, ""
	}
	return query[:indexOfFirstSpace], query[indexOfFirstSpace+1:]
}

func getNextExpr(query string) (squirrel.Sqlizer, string, error) {
	var sign, left, right string
	left, query = getNextUnit(query)
	sign, _ = getNextUnit(query)
	if isAComparisonSign(sign) {
		_, query = getNextUnit(query)
		right, query = getNextUnit(query)
		return parseExpression(left, sign, right), query, nil
	}
	if isABooleanSign(sign) {
		return squirrel.Eq{left: true}, query, nil
	}
	return nil, query, errors.New("Parsing error")
}

func getNextANDBlock(query string) (squirrel.Sqlizer, string, error) {
	var expr, nextExpr squirrel.Sqlizer
	var err error
	expr, query, err = getNextExpr(query)
	if err != nil {
		return nil, "", errors.New("Parsing error")
	}
	sign, _ := getNextUnit(query)
	for sign == "AND" || sign == "and" {
		_, query = getNextUnit(query)
		nextExpr, query, err = getNextExpr(query)
		if err != nil {
			return nil, "", errors.New("Parsing error")
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
				return nil, errors.New("Parsing error")
			}
		}
		sign, query = getNextUnit(query)
		if isABooleanSign(sign) {
			if sign == "OR" || sign == "or" {
				rightExpr, query, err = getNextANDBlock(query)
				if err != nil {
					return nil, errors.New("Parsing error")
				}
				expression = squirrel.Or{expression, rightExpr}
			}
		}
	}
	db = db.Where(expression)
	return &db, nil
}
