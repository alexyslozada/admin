package login

type PortUseCase interface {
	Login(email, password string) (string, error)
}
