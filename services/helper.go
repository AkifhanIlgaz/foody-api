package services

import (
	"fmt"
)

func columnWithDot(table, field string) string {
	return fmt.Sprintf("%v.%v", table, field)
}
