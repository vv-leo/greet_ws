package utils

import "strconv"

type RegexUtils struct {
}

// 是否数字
func (*RegexUtils) IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// 是否手机号
func (srv *RegexUtils) IsMobilePhone(s string) bool {
	if len(s) != 11 {
		return false
	}

	if s[0:1] != "1" {
		return false
	}

	return srv.IsNum(s)
}
