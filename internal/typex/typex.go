package typex

import (
	"fmt"
	"math"
	"time"

	"github.com/G0tem/go-servise-entity/internal/utilx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserId      uuid.UUID
	Username    string
	Email       string
	Roles       []string
	Permissions []string
	Exp         time.Time
}

const nsInSecond int64 = 1e9

func ParseJwtClaims(claims jwt.MapClaims) (*JwtClaims, error) {
	if claims == nil {
		return nil, fmt.Errorf("no JWT")
	}

	_, ok := claims["user_id"]
	if !ok {
		return nil, fmt.Errorf("incorrect JWT (missing user_id field)")
	}

	userId, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	_, ok = claims["exp"]
	if !ok {
		return nil, fmt.Errorf("incorrect JWT (missing exp field)")
	}

	sec, secFrac := math.Modf(claims["exp"].(float64))
	exp := time.Unix(int64(math.Floor(sec)), int64(secFrac*float64(nsInSecond)))

	roles := make([]string, 0)
	if claims["roles"] != nil {
		roles = utilx.Mapping(claims["roles"].([]interface{}), func(a interface{}) string { return a.(string) })
	}

	permissions := make([]string, 0)
	if claims["permissions"] != nil {
		permissions = utilx.Mapping(claims["permissions"].([]interface{}), func(a interface{}) string { return a.(string) })
	}

	return &JwtClaims{
		UserId:      userId,
		Username:    claims["username"].(string),
		Email:       claims["email"].(string),
		Roles:       roles,
		Permissions: permissions,
		Exp:         exp,
	}, err
}
