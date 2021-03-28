package verify

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	tm "github.com/sta-golang/go-lib-utils/time"
	"time"
)

const (
	iUser = "sta-golang"
)

var authErr = errors.New("auth err")

type jwtService struct {
	priKey []byte
}

type CustomClaims struct {
	Username string
	jwt.StandardClaims
}

var globalJWTService = jwtService{priKey: []byte("sta-golang")}

func NewJWTService() VerifyAuthService {
	return &globalJWTService
}

func (js *jwtService) CreateToken(info string, expire time.Duration) (string ,error)  {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Username:       info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tm.GetNowTime().Add(expire).Unix(),
			Issuer:    iUser,
		},
	})

	return token.SignedString(js.priKey)
}

func (js *jwtService) VerifyAuth(tokenStr string) (string, bool, error)  {
	token, err := jwt.ParseWithClaims(tokenStr,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return js.priKey, nil
	})
	if token == nil {
		return "", false, authErr
	}
	if claims, ok := token.Claims.(*CustomClaims); ok  {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return "", false, nil
		}
		if !token.Valid {
			return "", false, err
		}
		return claims.Username, true, nil
	}
	return "", false, err
}