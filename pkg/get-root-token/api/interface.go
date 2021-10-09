package api

type TokenInterface interface {
	Token() (string, error)
	TokenName() string
}
