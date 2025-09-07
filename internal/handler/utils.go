package handler

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/G0tem/go-servise-entity/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *Handler) getUserByEmail(e string) (*model.User, error) {
	var user model.User
	if err := h.db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (h *Handler) getUserByUsername(u string) (*model.User, error) {
	var user model.User
	if err := h.db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Map applies a function to each element of the input slice
func Map[A any, B any](input []A, mapper func(A) B) []B {
	result := make([]B, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}

type JwtClaims struct {
	userId   uuid.UUID
	username string
	email    string
	roles    []string
	exp      time.Time
}

const nsInSecond int64 = 1e9

func ParseJwtClaims(claims jwt.MapClaims) (*JwtClaims, error) {
	if claims == nil {
		return nil, fmt.Errorf("no JWT")
	}

	_, ok := claims["user_id"]
	if !ok {
		return nil, fmt.Errorf("not AXA socialweb JWT")
	}

	userId, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	sec, secFrac := math.Modf(claims["exp"].(float64))
	exp := time.Unix(int64(math.Floor(sec)), int64(secFrac*float64(nsInSecond)))

	roles := make([]string, 0)
	if claims["roles"] != nil {
		roles = Map(claims["roles"].([]interface{}), func(a interface{}) string { return a.(string) })
	}

	return &JwtClaims{
		userId:   userId,
		username: claims["username"].(string),
		email:    claims["email"].(string),
		roles:    roles,
		exp:      exp,
	}, err
}

func getAuthenticatedUserId(request *fasthttp.Request) *uuid.UUID {
	authorizationHeader := request.Header.Peek("Authorization")
	if authorizationHeader == nil {
		return nil
	}
	token, _ := jwt.Parse(strings.Split(string(authorizationHeader), "Bearer ")[1], nil)
	if token == nil {
		return nil
	}
	claims, err := ParseJwtClaims(token.Claims.(jwt.MapClaims))
	if err != nil {
		return nil
	}

	return &claims.userId
}
