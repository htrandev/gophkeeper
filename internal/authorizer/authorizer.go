package authorizer

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// claims описывает набор JWT claims.
type claims struct {
	jwt.RegisteredClaims
	UserID string `json:"Id"`
}

// Authorizer — структура, отвечающая за создание,
// валидацию и расшифровку JWT-токенов, а также за хэширование паролей.
type Authorizer struct {
	secretKey []byte
	ttl       time.Duration
}

// New возвращает новый экземпляр Authorizer.
func New(key string, ttl time.Duration) *Authorizer {
	return &Authorizer{
		secretKey: []byte(key),
		ttl:       ttl,
	}
}

// HashPassword хэширует пароль с помощью bcrypt.
func (a *Authorizer) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("generate hash from password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword проверяет соответствие введённого пароля ранее сохранённому хэшу.
func (a *Authorizer) ValidatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

// Token создаёт новый JWT-токен для пользователя с заданным id.
func (a *Authorizer) Token(id string) (string, error) {
	ttl := time.Now().Add(a.ttl)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(ttl),
		},
		UserID: id,
	})

	token, err := t.SignedString(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}
	return token, nil
}

// GetIDFromToken парсит JWT-токен и возвращает ID пользователя,
// если токен валиден и не истёк.
func (a *Authorizer) GetIDFromToken(token string) (string, error) {
	parseen, err := jwt.ParseWithClaims(
		token,
		&claims{},
		func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedMethod
			}

			return a.secretKey, nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, ok := parseen.Claims.(*claims)
	if !ok {
		return "", ErrInvalidClaims
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return "", ErrTokenExpired
	}

	return claims.UserID, nil
}
