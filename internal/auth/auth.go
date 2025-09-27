package auth

type AuthValidator interface {
	ValidateToken(token string) (map[string]interface{}, error)
}
