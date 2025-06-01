package model

type Usuario struct {
	Id    int    `json:"id_usuario"`
	Nome  string `json:"nome_usuario"`
	Login string `json:"login_usuario"`
	Senha string `json:"senha_usuario"`
}
