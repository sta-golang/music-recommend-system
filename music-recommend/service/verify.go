package service

import (
	"time"
)

/**

 */
type VerifyAuthService interface {
	CreateToken(info string, expire time.Duration) (string, error)
	VerifyAuth(token string) (string ,bool, error)
}