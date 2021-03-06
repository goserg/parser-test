package units

import (
	"strings"
	"unicode"
)

//ColumnUnit структура для хранения имени колонки
type ColumnUnit struct {
	Name string
}

//IsValid проверка валидности имени колонки
func (c *ColumnUnit) IsValid() bool {
	firstSymbol := string(c.Name[0])
	lastSymbol := string(c.Name[len(c.Name)-1])
	if firstSymbol == `"` && lastSymbol == `"` { // К идентификаторам в кавычках не применяются правила валидности
		return true
	}
	for i, symbol := range c.Name {
		if i == 0 {
			if !unicode.IsLetter(symbol) && symbol != rune("_"[0]) {
				return false
			}
		} else {
			if !unicode.IsLetter(symbol) && strings.Index("0123456789_.", string(symbol)) == -1 {
				return false
			}
		}
	}

	return true
}
