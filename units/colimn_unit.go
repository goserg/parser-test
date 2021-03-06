package units

import "strings"

//ColumnUnit структура для хранения имени колонки
type ColumnUnit struct {
	Name string
}

//IsValid проверка валидности имени колонки
func (c *ColumnUnit) IsValid() bool {
	if strings.Index("0123456789", string(c.Name[0])) != -1 {
		return false
	}

	return true
}
