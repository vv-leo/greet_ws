package utils

import "strings"

type StringUtils struct {
}

func (*StringUtils) InStrArray(arr []string, str string) bool {
	str = strings.TrimSpace(str)
	for _, s := range arr {
		if strings.TrimSpace(s) == str {
			return true
		}
	}
	return false
}
