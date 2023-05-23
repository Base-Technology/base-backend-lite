package token

import (
	"fmt"
	"time"

	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type UserClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {
	user := &UserClaims{ID: id}
	user.IssuedAt = time.Now().Unix()
	user.ExpiresAt = time.Now().Add(time.Second * time.Duration(conf.Conf.ServerConf.TokenExpireTime)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	signedToken, err := token.SignedString([]byte(conf.Conf.ServerConf.TokenSecret))
	if err != nil {
		return "", fmt.Errorf("sign token error, %v", err)
	}
	return signedToken, nil
}

func VerifyToken(token string) (uint, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.ServerConf.TokenSecret), nil
	})
	if err != nil {
		return 0, errors.Errorf("parse token error, %v", err)
	}
	if err := t.Claims.Valid(); err != nil {
		return 0, errors.Errorf("invalid token, error: %v", err)
	}
	user, ok := t.Claims.(*UserClaims)
	if !ok {
		return 0, errors.Errorf("invalid token type: %T", t.Claims)
	}
	return user.ID, nil
}
