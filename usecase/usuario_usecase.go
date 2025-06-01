package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type UsuarioUsecase struct {
	repository repository.UsuarioRepository
}

func NewUsuarioUseCase(repo repository.UsuarioRepository) UsuarioUsecase {
	return UsuarioUsecase{
		repository: repo,
	}
}

func (uu *UsuarioUsecase) GetUsuarios() ([]model.Usuario, error) {
	return uu.repository.GetUsuarios()
}

func (uu *UsuarioUsecase) CreateUsuario(usuario model.Usuario) (model.Usuario, error) {
	id, err := uu.repository.CreateUsuario(usuario)
	if err != nil {
		return model.Usuario{}, err
	}

	usuario.Id = id

	return usuario, nil
}

func (uu *UsuarioUsecase) GetUsuarioById(id_usuario int) (*model.Usuario, error) {
	usuario, err := uu.repository.GetUsuarioById(id_usuario)
	if err != nil {
		return nil, err
	}

	return usuario, nil
}

func (uu *UsuarioUsecase) UpdateUsuarioById(id_usuario int, usuario *model.Usuario) error {
	err := uu.repository.UpdateUsuarioById(id_usuario, usuario)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UsuarioUsecase) SoftDeleteUsuarioById(id_usuario int) error {
	return uu.repository.SoftDeleteUsuarioById(id_usuario)
}
