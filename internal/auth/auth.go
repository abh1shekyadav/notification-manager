package auth

type AuthValidator interface {
	ValidateToken(token string) (string, error)
}
