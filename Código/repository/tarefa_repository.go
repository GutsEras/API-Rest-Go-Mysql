package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type TarefaRepository struct {
	connection *sql.DB
}

func NewTarefaRepository(connection *sql.DB) TarefaRepository {
	return TarefaRepository{
		connection: connection,
	}
}

func (tr *TarefaRepository) GetTarefas() ([]model.Tarefa, error) {
	query := "SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa"
	rows, err := tr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Tarefa{}, err
	}
	defer rows.Close()

	var tarefaList []model.Tarefa
	for rows.Next() {
		var tarefa model.Tarefa
		err := rows.Scan(
			&tarefa.Id,
			&tarefa.Nome,
			&tarefa.Conteudo,
			&tarefa.UsuarioResp,
			&tarefa.Finalizado,
		)
		if err != nil {
			fmt.Println(err)
			return []model.Tarefa{}, err
		}
		tarefaList = append(tarefaList, tarefa)
	}

	return tarefaList, nil
}

func (tr *TarefaRepository) CreateTarefa(tarefa model.Tarefa) (int, error) {
	result, err := tr.connection.Exec(
		"INSERT INTO tarefa (nome, conteudo, usuario_responsavel, finalizado) VALUES (?, ?, ?, ?)",
		tarefa.Nome,
		tarefa.Conteudo,
		tarefa.UsuarioResp,
		tarefa.Finalizado,
	)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(id), nil
}

func (tr *TarefaRepository) GetTarefaById(id_tarefa int) (*model.Tarefa, error) {
	query, err := tr.connection.Prepare("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa WHERE id = ?")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer query.Close()

	var tarefa model.Tarefa
	err = query.QueryRow(id_tarefa).Scan(
		&tarefa.Id,
		&tarefa.Nome,
		&tarefa.Conteudo,
		&tarefa.UsuarioResp,
		&tarefa.Finalizado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &tarefa, nil
}

func (tr *TarefaRepository) UpdateTarefaById(id_tarefa int, tarefa *model.Tarefa) error {
	query, err := tr.connection.Prepare(`
		UPDATE tarefa 
		SET nome = ?, conteudo = ?, usuario_responsavel = ?, finalizado = ?
		WHERE id = ?
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	result, err := query.Exec(
		tarefa.Nome,
		tarefa.Conteudo,
		tarefa.UsuarioResp,
		tarefa.Finalizado,
		id_tarefa,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (tr *TarefaRepository) SoftDeleteTarefaById(id_tarefa int) error {
	query, err := tr.connection.Prepare("UPDATE tarefa SET ativo = 'N' WHERE id = ? AND ativo = 'A'")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	result, err := query.Exec(id_tarefa)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (tr *TarefaRepository) GetTarefasByUsuarioId(usuarioId string) ([]model.Tarefa, error) {
	query := `
		SELECT id, nome, conteudo, usuario_responsavel, finalizado 
		FROM tarefa 
		WHERE usuario_responsavel = ? AND ativo = 'A'
	`

	rows, err := tr.connection.Query(query, usuarioId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var tarefas []model.Tarefa
	for rows.Next() {
		var tarefa model.Tarefa
		err := rows.Scan(
			&tarefa.Id,
			&tarefa.Nome,
			&tarefa.Conteudo,
			&tarefa.UsuarioResp,
			&tarefa.Finalizado,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tarefas = append(tarefas, tarefa)
	}

	return tarefas, nil
}
