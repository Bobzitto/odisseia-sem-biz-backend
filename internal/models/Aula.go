package models

import "time"

type Aula struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Size          string     `json:"size"`
	Active        bool       `json:"active"`
	Review        *float64   `json:"review,omitempty"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
	Materias      []*Materia `json:"materias,omitempty"`
	MateriasArray []int      `json:"materias_array,omitempty"`
}

type Materia struct {
	ID        int       `json:"id"`
	Materia   string    `json:"materia"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
