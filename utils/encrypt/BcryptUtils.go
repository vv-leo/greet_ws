package encrypt

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

///bcrypt生成、比对
type BcryptUtils struct {
}

//生成密码
func (_ *BcryptUtils) GeneratePassword(src string) (string, error) {
	byteArr, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	return string(byteArr), err
}

//比对密码是否相等
func (_ *BcryptUtils) ValidatePassword(src string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(src)); err != nil {
		return false, errors.New("密码错误")
	}
	return true, nil
}
