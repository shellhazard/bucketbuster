package utils

import (
	"reflect"
	"strings"
)

func ReverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func CleanStringSlice(input []string) []string {
	var output []string
	for _, elem := range input {
		trimmed := strings.TrimSpace(elem)
		if len(elem) > 0 {
			output = append(output, trimmed)
		}
	}
	return output
}
