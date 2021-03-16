package verify_auth

import (
	"fmt"
	"testing"
	"time"
)

func TestNewJWTService(t *testing.T) {
	service := NewJWTService()
	token, err := service.CreateToken("userName", time.Second)
	fmt.Println(err)
	sss := token
	bys := []byte(sss)
	bys[4] = '2'
	fmt.Println(service.VerifyAuth("123456"))
	fmt.Println(service.VerifyAuth(string(bys)))
	fmt.Println(service.VerifyAuth(token))
	time.Sleep(time.Second * 2)
	fmt.Println(service.VerifyAuth(token))
}
