package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	EditarUmaAula(id int) (*models.Aula, []*models.Materia, error)
	TodaAula() ([]*models.Aula, error)
	UmaAula(id int) (*models.Aula, error)
	InserirAula(aula models.Aula) (int, error)
	TodasMaterias() ([]*models.Materia, error)
	AtualizarMateria(id int, materiasID []int) error
	AtualizarAula(aula models.Aula) error
	DeleteAula(id int) error
}
