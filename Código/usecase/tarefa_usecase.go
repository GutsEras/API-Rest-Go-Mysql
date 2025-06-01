package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type TarefaUsecase struct {
	repository repository.TarefaRepository
}

func NewTarefaUseCase(repo repository.TarefaRepository) TarefaUsecase {
	return TarefaUsecase{
		repository: repo,
	}
}

func (tu *TarefaUsecase) GetTarefas() ([]model.Tarefa, error) {
	return tu.repository.GetTarefas()
}

func (tu *TarefaUsecase) CreateTarefa(tarefa model.Tarefa) (model.Tarefa, error) {
	id, err := tu.repository.CreateTarefa(tarefa)
	if err != nil {
		return model.Tarefa{}, err
	}

	tarefa.Id = id
	return tarefa, nil
}

func (tu *TarefaUsecase) GetTarefaById(id_tarefa int) (*model.Tarefa, error) {
	tarefa, err := tu.repository.GetTarefaById(id_tarefa)
	if err != nil {
		return nil, err
	}
	return tarefa, nil
}

func (tu *TarefaUsecase) UpdateTarefaById(id_tarefa int, tarefa *model.Tarefa) error {
	return tu.repository.UpdateTarefaById(id_tarefa, tarefa)
}

func (tu *TarefaUsecase) SoftDeleteTarefaById(id_tarefa int) error {
	return tu.repository.SoftDeleteTarefaById(id_tarefa)
}

func (tu *TarefaUsecase) GetTarefasByUsuarioId(usuarioId string) ([]model.Tarefa, error) {
	return tu.repository.GetTarefasByUsuarioId(usuarioId)
}
