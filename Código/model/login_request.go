package model

type LoginRequest struct {
	Login string `json:"login" example:"usuario123"`
	Senha string `json:"senha" example:"senhaSegura"`
}
