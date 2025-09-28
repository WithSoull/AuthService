package service

type AuthService interface {
	Login(string, string) string
	GetRefreshToken(string) string
	GetAccessToken(string) string
}

type AccessService interface {
	Check(string) error
}
