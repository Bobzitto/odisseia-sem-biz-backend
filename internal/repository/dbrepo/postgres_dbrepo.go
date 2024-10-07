package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
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
			id, name, size, active, review, created_at, updated_at
		from
			aulas
		order by
			name
		`
	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aulas []*models.Aula

	for rows.Next() {
		var aula models.Aula
		err := rows.Scan(
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

		aulas = append(aulas, &aula)
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

func (a *PostgresDBRepo) InserirAula(aula models.Aula) (int, error){
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