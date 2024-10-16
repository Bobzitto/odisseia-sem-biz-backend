package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"log"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *PostgresDBRepo) TodaAula() ([]*models.Aula, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			a.id, a.name, a.size, a.active, a.review, a.created_at, a.updated_at, 
			m.id as materia_id, m.materia as materia_name
		from
			aulas a
		left join aulas_materias am on am.aula_id = a.id
		left join materias m on am.materia_id = m.id
		order by
			a.name
	`
	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aulasMap = make(map[int]*models.Aula)

	for rows.Next() {
		var aula models.Aula
		var materiaID int
		var materiaName string

		err := rows.Scan(
			&aula.ID,
			&aula.Name,
			&aula.Size,
			&aula.Active,
			&aula.Review,
			&aula.CreatedAt,
			&aula.UpdatedAt,
			&materiaID,
			&materiaName,
		)
		if err != nil {
			return nil, err
		}

		if existingAula, exists := aulasMap[aula.ID]; exists {
			existingAula.Materias = append(existingAula.Materias, &models.Materia{
				ID:      materiaID,
				Materia: materiaName,
			})
		} else {
			aula.Materias = []*models.Materia{
				{
					ID:      materiaID,
					Materia: materiaName,
				},
			}
			aulasMap[aula.ID] = &aula
		}
	}

	var aulas []*models.Aula
	for _, aula := range aulasMap {
		aulas = append(aulas, aula)
	}
	return aulas, nil
}

func (a *PostgresDBRepo) UmaAula(id int) (*models.Aula, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, name, size, active, review, created_at, updated_at
		from
			aulas
		where id = $1`

	row := a.DB.QueryRowContext(ctx, query, id)

	var aula models.Aula

	err := row.Scan(
		&aula.ID,
		&aula.Name,
		&aula.Size,
		&aula.Active,
		&aula.Review,
		&aula.CreatedAt,
		&aula.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	//get materias
	query = `SELECT m.id, m.materia
			FROM aulas_materias am
			LEFT JOIN materias m ON am.materia_id = m.id
			WHERE am.aula_id = $1
			ORDER BY m.materia;`

	rows, err := a.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var materias []*models.Materia
	for rows.Next() {
		var m models.Materia
		err := rows.Scan(
			&m.ID,
			&m.Materia,
		)
		if err != nil {
			return nil, err
		}

		materias = append(materias, &m)
	}

	aula.Materias = materias

	return &aula, err
}

func (a *PostgresDBRepo) EditarUmaAula(id int) (*models.Aula, []*models.Materia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, name, size, active, review, created_at, updated_at
		from
			aulas
		where id = $1`

	row := a.DB.QueryRowContext(ctx, query, id)

	var aula models.Aula

	err := row.Scan(
		&aula.ID,
		&aula.Name,
		&aula.Size,
		&aula.Active,
		&aula.Review,
		&aula.CreatedAt,
		&aula.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	//get materias
	query = `SELECT m.id, m.materia
			FROM aulas_materias am
			LEFT JOIN materias m ON am.materia_id = m.id
			WHERE am.aula_id = $1
			ORDER BY m.materia;`

	rows, err := a.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}
	defer rows.Close()

	var materias []*models.Materia
	var MateriasArray []int

	for rows.Next() {
		var m models.Materia
		err := rows.Scan(
			&m.ID,
			&m.Materia,
		)
		if err != nil {
			return nil, nil, err
		}

		materias = append(materias, &m)
		MateriasArray = append(MateriasArray, m.ID)
	}

	aula.Materias = materias
	aula.MateriasArray = MateriasArray

	var todasMaterias []*models.Materia

	query = "select id, materia from materias order by materia"
	mRows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer mRows.Close()

	for mRows.Next() {
		var m models.Materia
		err := mRows.Scan(
			&m.ID,
			&m.Materia,
		)
		if err != nil {
			return nil, nil, err
		}

		todasMaterias = append(todasMaterias, &m)
	}

	return &aula, todasMaterias, err
}

func (a *PostgresDBRepo) TodasMaterias() ([]*models.Materia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, materia, created_at, updated_at from materias order by materia`

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materias []*models.Materia

	for rows.Next() {
		var m models.Materia
		err := rows.Scan(
			&m.ID,
			&m.Materia,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		materias = append(materias, &m)
	}

	return materias, nil
}

func (a *PostgresDBRepo) InserirAula(aula models.Aula) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into aulas (name, size, active, review, created_at, updated_at )
			 values ($1, $2, $3, $4, $5, $6) returning id`

	var newID int

	err := a.DB.QueryRowContext(ctx, stmt,
		aula.Name,
		aula.Size,
		aula.Active,
		aula.Review,
		aula.CreatedAt,
		aula.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (a *PostgresDBRepo) AtualizarMateria(id int, materiasID []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	log.Printf("Deleting existing materias for aula_id = %d", id)
	stmt := `delete from aulas_materias where aula_id = $1`

	_, err := a.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Printf("Error deleting materias: %v", err)
		return err
	}

	log.Printf("Inserting new materias for aula_id = %d: %v", id, materiasID)
	for _, n := range materiasID {
		stmt := `insert into aulas_materias	(aula_id, materia_id) values ($1, $2)`
		_, err := a.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			log.Printf("Error inserting into aulas_materias: %v", err)
			return err
		}
	}
	return nil
}

func (a *PostgresDBRepo) AtualizarAula(aula models.Aula) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update aulas set name=$1, size=$2, active=$3, review=$4,
			updated_at=$5 where id = $6`

	_, err := a.DB.ExecContext(ctx, stmt,
		aula.Name,
		aula.Size,
		aula.Active,
		aula.Review,
		aula.UpdatedAt,
		aula.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgresDBRepo) DeleteAula(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from aulas where id = $1`

	_, err := a.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// uma turma
func (t *PostgresDBRepo) UmaTurma(id int) (*models.Turma, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, name, school, year, created_at, updated_at
		from
			turmas
		where id = $1`

	row := t.DB.QueryRowContext(ctx, query, id)

	var turma models.Turma

	err := row.Scan(
		&turma.ID,
		&turma.Name,
		&turma.School,
		&turma.Year,
		&turma.CreatedAt,
		&turma.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &turma, err

}

//Todas as turmas

func (t *PostgresDBRepo) TodaTurma() ([]*models.Turma, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, name, school, year, created_at, updated_at
		from
			turmas
		order by
			name
	`
	rows, err := t.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turmas []*models.Turma

	for rows.Next() {
		var turma models.Turma

		err := rows.Scan(
			&turma.ID,
			&turma.Name,
			&turma.School,
			&turma.Year,
			&turma.CreatedAt,
			&turma.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Append the scanned turma to the slice
		turmas = append(turmas, &turma)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return turmas, nil
}

// Inserir uma turma
func (t *PostgresDBRepo) InserirTurma(turma models.Turma) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into turmas (name, school, year, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id`

	var newID int

	err := t.DB.QueryRowContext(ctx, stmt,
		turma.Name,
		turma.School,
		turma.Year,
		turma.CreatedAt,
		turma.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// Atualizar uma turma
func (t *PostgresDBRepo) AtualizarTurma(turma models.Turma) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update turma set name = $1, school = $2, year = $3,
		 updated_at = $4 where id = $5`

	_, err := t.DB.ExecContext(ctx, stmt,
		turma.Name,
		turma.School,
		turma.Year,
		turma.UpdatedAt,
		turma.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// editar uma aula

func (t *PostgresDBRepo) EditarUmaTurma(id int) (*models.Turma, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, name, school, year, created_at, updated_at
		from
			turmas
		where id = $1`

	row := t.DB.QueryRowContext(ctx, query, id)

	var turma models.Turma

	err := row.Scan(
		&turma.ID,
		&turma.School,
		&turma.Year,
		&turma.CreatedAt,
		&turma.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &turma, err
}

//Deletar uma turma

func (t *PostgresDBRepo) DeleteTurma(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from turmas where id = $1`

	_, err := t.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
