package utils

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

type UUidUtils struct{}

func (UUidUtils) GetUUID() string {
	u2 := uuid.NewV4()

	str := strings.ReplaceAll(u2.String(), "-", "")
	str = strings.ToUpper(str)
	return str
}
