package main

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-api/model"
	"go-api/repository"
)

func TestGetTarefaById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)
	tarefaId := 1
	expected := model.Tarefa{
		Id: tarefaId, Nome: "Teste", Conteudo: "Conteudo", UsuarioResp: "user1", Finalizado: "N",
	}

	rows := sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
		AddRow(expected.Id, expected.Nome, expected.Conteudo, expected.UsuarioResp, expected.Finalizado)

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa WHERE id = ?")).
		ExpectQuery().WithArgs(tarefaId).WillReturnRows(rows)

	tarefa, err := repo.GetTarefaById(tarefaId)
	assert.NoError(t, err)
	assert.Equal(t, &expected, tarefa)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTarefas(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
		AddRow(1, "Tarefa1", "Conteudo1", "user1", false).
		AddRow(2, "Tarefa2", "Conteudo2", "user2", true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa")).
		WillReturnRows(rows)

	tarefas, err := repo.GetTarefas()
	assert.NoError(t, err)
	assert.Len(t, tarefas, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTarefa(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)
	tarefa := model.Tarefa{Nome: "Nova", Conteudo: "Teste", UsuarioResp: "user1", Finalizado: "N"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO tarefa (nome, conteudo, usuario_responsavel, finalizado) VALUES (?, ?, ?, ?)")).
		WithArgs(tarefa.Nome, tarefa.Conteudo, tarefa.UsuarioResp, tarefa.Finalizado).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.CreateTarefa(tarefa)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTarefaById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)
	tarefa := &model.Tarefa{Nome: "Atualizada", Conteudo: "Atualizado", UsuarioResp: "user1", Finalizado: "S"}

	mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE tarefa 
		SET nome = ?, conteudo = ?, usuario_responsavel = ?, finalizado = ?
		WHERE id = ?`)).
		ExpectExec().WithArgs(tarefa.Nome, tarefa.Conteudo, tarefa.UsuarioResp, tarefa.Finalizado, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateTarefaById(1, tarefa)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDeleteTarefaById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)

	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE tarefa SET ativo = 'N' WHERE id = ? AND ativo = 'A'")).
		ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SoftDeleteTarefaById(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTarefasByUsuarioId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTarefaRepository(db)
	usuarioId := "user1"

	rows := sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
		AddRow(1, "Tarefa1", "Conteudo1", usuarioId, false)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, nome, conteudo, usuario_responsavel, finalizado 
		FROM tarefa 
		WHERE usuario_responsavel = ? AND ativo = 'A'`)).
		WithArgs(usuarioId).
		WillReturnRows(rows)

	tarefas, err := repo.GetTarefasByUsuarioId(usuarioId)
	assert.NoError(t, err)
	assert.Len(t, tarefas, 1)
	assert.Equal(t, usuarioId, tarefas[0].UsuarioResp)
	assert.NoError(t, mock.ExpectationsWereMet())
}
