package units

import (
	"strings"
)

//SignUnit структура для работы с операторами AND OR и операторами сравнения
type SignUnit struct {
	Value string
}

//IsANDSign возвращает true если значение знака AND
func (s *SignUnit) IsANDSign() bool {
	return strings.ToLower(s.Value) == "and"
}

//IsORSign возвращает true если значение знака OR
func (s *SignUnit) IsORSign() bool {
	return strings.ToLower(s.Value) == "or"
}

//IsABoolSign возвращает true если значение знака AND или OR
func (s *SignUnit) IsABoolSign() bool {
	return s.IsANDSign() || s.IsORSign()
}

//IsAComparisonSign возвращает true если знак оператор сравнения
func (s *SignUnit) IsAComparisonSign() bool {
	for _, comp := range []string{"=", "!=", "<>", "~", "!~", "~*", "!~*", ">", ">=", "<", "<="} {
		if comp == s.Value {
			return true
		}
	}
	return false
}
