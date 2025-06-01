package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type UsuarioRepository struct {
	connection *sql.DB
}

func NewUsuarioRepository(connection *sql.DB) UsuarioRepository {
	return UsuarioRepository{
		connection: connection,
	}
}

func (ur *UsuarioRepository) GetUsuarios() ([]model.Usuario, error) {
	query := "SELECT id, nome, login, senha FROM usuario"
	rows, err := ur.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Usuario{}, err
	}

	var usuarioList []model.Usuario
	var usuarioObj model.Usuario

	for rows.Next() {
		err = rows.Scan(
			&usuarioObj.Id,
			&usuarioObj.Nome,
			&usuarioObj.Login,
			&usuarioObj.Senha)

		if err != nil {
			fmt.Println(err)
			return []model.Usuario{}, err
		}

		usuarioList = append(usuarioList, usuarioObj)
	}

	rows.Close()

	return usuarioList, nil
}

func (ur *UsuarioRepository) CreateUsuario(usuario model.Usuario) (int, error) {
	result, err := ur.connection.Exec(
		"INSERT INTO usuario (nome, login, senha) VALUES (?, ?, ?)",
		usuario.Nome,
		usuario.Login,
		usuario.Senha,
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

func (ur *UsuarioRepository) GetUsuarioById(id_usuario int) (*model.Usuario, error) {
	query, err := ur.connection.Prepare("SELECT id, nome, login, senha FROM usuario WHERE id = ?")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer query.Close()

	var usuario model.Usuario

	err = query.QueryRow(id_usuario).Scan(
		&usuario.Id,
		&usuario.Nome,
		&usuario.Login,
		&usuario.Senha,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &usuario, nil
}

func (ur *UsuarioRepository) UpdateUsuarioById(id_usuario int, usuario *model.Usuario) error {
	query, err := ur.connection.Prepare(`
		UPDATE usuario 
		SET nome = ?, login = ?, senha = ?
		WHERE id = ?
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	result, err := query.Exec(
		usuario.Nome,
		usuario.Login,
		usuario.Senha,
		id_usuario,
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
		// Nenhum registro foi atualizado, talvez o ID não exista
		return sql.ErrNoRows
	}

	return nil
}

func (ur *UsuarioRepository) SoftDeleteUsuarioById(id_usuario int) error {
	query, err := ur.connection.Prepare("UPDATE usuario SET ativo = 'I' WHERE id = ? AND ativo = 'A'")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	result, err := query.Exec(id_usuario)
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

func (ur *UsuarioRepository) GetUsuarioByLogin(login string) (*model.Usuario, error) {
	query := "SELECT id, nome, login, senha FROM usuario WHERE login = ?"
	row := ur.connection.QueryRow(query, login)

	var usuario model.Usuario
	err := row.Scan(&usuario.Id, &usuario.Nome, &usuario.Login, &usuario.Senha)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &usuario, nil
}
