package model

type Tarefa struct {
	Id          int    `json:"id_tarefa"`
	Nome        string `json:"nome_tarefa"`
	Conteudo    string `json:"conteudo_tarefa"`
	UsuarioResp string `json:"usuario_responsavel_tarefa"`
	Finalizado  string `json:"finalizado"`
}
