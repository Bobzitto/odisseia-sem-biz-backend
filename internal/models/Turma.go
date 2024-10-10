package models

import "time"

type Turma struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	School    string `json:"school"`
	Year      string `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}