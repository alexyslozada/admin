package login

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	EDtimer "gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
	portUser "gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
)

type UseCase struct {
	userUseCase portUser.UserUseCase
	timer       EDtimer.Timer
}

func NewUseCase(userUseCase portUser.UserUseCase, timer EDtimer.Timer) UseCase {
	return UseCase{
		userUseCase: userUseCase,
		timer:       timer,
	}
}

func (uc UseCase) Login(email, password string) (string, error) {
	u, err := uc.userUseCase.Login(email, password)
	if err != nil {
		return "", err
	}

	return uc.CreateToken(u)
}

func (uc UseCase) CreateToken(u domain.User) (string, error) {
	claims := domain.JWTCustomClaims{
		ID:    u.ID,
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(uc.timer.Now().Add(24 * time.Hour)),
			Issuer:    "EDteam",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (uc UseCase) ValidateToken(tokenString string) (bool, error) {
	claims := domain.JWTCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, uc.getJWTSecret)
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, err
	}

	return true, nil
}

func (uc UseCase) getJWTSecret(_ *jwt.Token) (any, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}
