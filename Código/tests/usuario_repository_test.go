package main

import (
	"regexp"
	"testing"

	"go-api/model"
	"go-api/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsuarioById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)
	expected := model.Usuario{
		Id:    1,
		Nome:  "João",
		Login: "joao123",
		Senha: "senha",
	}

	rows := sqlmock.NewRows([]string{"id", "nome", "login", "senha"}).
		AddRow(expected.Id, expected.Nome, expected.Login, expected.Senha)

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, nome, login, senha FROM usuario WHERE id = ?")).
		ExpectQuery().
		WithArgs(expected.Id).
		WillReturnRows(rows)

	result, err := repo.GetUsuarioById(expected.Id)
	assert.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestGetUsuarios(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nome", "login", "senha"}).
		AddRow(1, "João", "joao123", "senha").
		AddRow(2, "Maria", "maria123", "senha123")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, nome, login, senha FROM usuario")).
		WillReturnRows(rows)

	result, err := repo.GetUsuarios()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Maria", result[1].Nome)
}

func TestCreateUsuario(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)

	input := model.Usuario{
		Nome:  "Maria",
		Login: "maria123",
		Senha: "senha123",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO usuario (nome, login, senha) VALUES (?, ?, ?)")).
		WithArgs(input.Nome, input.Login, input.Senha).
		WillReturnResult(sqlmock.NewResult(10, 1))

	id, err := repo.CreateUsuario(input)
	assert.NoError(t, err)
	assert.Equal(t, 10, id)
}

func TestUpdateUsuarioById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)

	user := &model.Usuario{
		Nome:  "Atualizado",
		Login: "atualizado123",
		Senha: "novaSenha",
	}

	mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE usuario 
		SET nome = ?, login = ?, senha = ?
		WHERE id = ?`)).
		ExpectExec().
		WithArgs(user.Nome, user.Login, user.Senha, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUsuarioById(1, user)
	assert.NoError(t, err)
}

func TestSoftDeleteUsuarioById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)

	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE usuario SET ativo = 'I' WHERE id = ? AND ativo = 'A'")).
		ExpectExec().
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SoftDeleteUsuarioById(5)
	assert.NoError(t, err)
}

func TestGetUsuarioByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUsuarioRepository(db)

	expected := model.Usuario{
		Id:    1,
		Nome:  "João",
		Login: "joao123",
		Senha: "senha",
	}

	rows := sqlmock.NewRows([]string{"id", "nome", "login", "senha"}).
		AddRow(expected.Id, expected.Nome, expected.Login, expected.Senha)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, nome, login, senha FROM usuario WHERE login = ?")).
		WithArgs("joao123").
		WillReturnRows(rows)

	result, err := repo.GetUsuarioByLogin("joao123")
	assert.NoError(t, err)
	assert.Equal(t, &expected, result)
}
