package ports

type LoginUseCase interface {
	Login(email, password string) (string, error)
	ValidateToken(tokenString string) (bool, error)
}
