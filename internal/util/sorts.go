package util

import "fmt"

func ToSorts(sortCol string, sortOrd string) []string {
	if len(sortCol) == 0 {
		sortCol = "id"
	}
	if len(sortOrd) == 0 {
		sortOrd = "asc"
	}
	return []string{fmt.Sprintf("%s %s", sortCol, sortOrd)}
}
