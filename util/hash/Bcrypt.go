package hash

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Bcrypt加密
func BcryptMake(pwd string) string {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// 加密验证
func BcryptMakeCheck(pwd, hashedPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
