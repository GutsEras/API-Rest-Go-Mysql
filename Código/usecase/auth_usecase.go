package usecase

import (
	"errors"
	"go-api/config"
	"go-api/repository"
)

type AuthUsecase struct {
	UsuarioRepo repository.UsuarioRepository
}

func NewAuthUsecase(repo repository.UsuarioRepository) *AuthUsecase {
	return &AuthUsecase{UsuarioRepo: repo}
}

func (uc *AuthUsecase) Login(login, senha string) (string, error) {
	usuario, err := uc.UsuarioRepo.GetUsuarioByLogin(login)
	if err != nil {
		return "", err
	}
	if usuario == nil || usuario.Senha != senha {
		return "", errors.New("login ou senha inválidos")
	}

	return config.GenerateToken(usuario.Id)
}
