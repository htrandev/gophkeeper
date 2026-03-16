package contracts


// Authorizer описывает интерфейс работы с авторизацией и аутентификацией.
type Authorizer interface {
	HashPassword(password string) (hahs string, err error)
	ValidatePassword(hash, password string) bool
	Token(id string) (token string, err error)
	GetIDFromToken(token string) (string, error)
}
