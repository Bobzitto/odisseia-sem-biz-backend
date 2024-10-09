package main

import (
	"backend/internal/graph"
	"backend/internal/models"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func (app *application) authenticate(c echo.Context) error {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	log.Println("Authenticate function hit")
	err := app.ReadJSON(c, &requestPayload)
	if err != nil {
		return app.WriteJSON(c, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}

	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		return app.errorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)

	}

	//validate pswd
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		return app.errorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)

	}

	//create jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	//generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		return app.errorJSON(c, err)

	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(c.Response().Writer, refreshCookie)

	return app.WriteJSON(c, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(c echo.Context) error {
	cookies := c.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			//parse token for claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})

			if err != nil {
				return app.errorJSON(c, errors.New("unauthorized"), http.StatusUnauthorized)
			}

			//get user id
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				return app.errorJSON(c, errors.New("unknown user"), http.StatusUnauthorized)
			}

			//retrieve user
			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				return app.errorJSON(c, errors.New("unknown user"), http.StatusUnauthorized)
			}

			//create a new jwtUser
			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				return app.errorJSON(c, errors.New("error generating token"), http.StatusUnauthorized)
			}

			c.SetCookie(app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			return app.WriteJSON(c, http.StatusOK, tokenPairs)
		}
	}
	return app.errorJSON(c, errors.New("refresh token not found"), http.StatusUnauthorized)
}

func (app *application) logOut(c echo.Context) error {
	expiredCookie := app.auth.GetExpiredRefreshCookie()
	c.SetCookie(expiredCookie)

	return c.NoContent(http.StatusAccepted)
}

func (app *application) TodasAulas(c echo.Context) error {
	aulas, err := app.DB.TodaAula()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error fetching auals"})

	}
	return c.JSON(http.StatusOK, aulas)
}

func (app *application) ListaAulas(c echo.Context) error {
	aulas, err := app.DB.TodaAula()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}

	return c.JSON(http.StatusOK, aulas)
}

func (app *application) PegarAula(c echo.Context) error {
	id := c.Param("id")

	//convert id to int
	aulaID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de aula inválido"})
	}

	//fetch aula from db
	aula, err := app.DB.UmaAula(aulaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, aula)
}

func (app *application) EditarAula(c echo.Context) error {
	id := c.Param("id")

	aulaID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de aula invalido"})
	}

	aula, materias, err := app.DB.EditarUmaAula(aulaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var payload = struct {
		Aula     *models.Aula      `json:"aula"`
		Materias []*models.Materia `json:"materias"`
	}{
		aula,
		materias,
	}

	return c.JSON(http.StatusOK, payload)
}

func (app *application) TodasMaterias(c echo.Context) error {
	materias, err := app.DB.TodasMaterias()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Matéria invalida"})
	}

	return c.JSON(http.StatusOK, materias)
}

func (app *application) InserirAula(c echo.Context) error {
	log.Println("inseriraula endpoint hit")
	var aula models.Aula

	err := c.Bind(&aula)
	if err != nil {
		log.Printf("bind error: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid"})
	}

	var materiasID []int
	for _, materia := range aula.Materias {
		materiasID = append(materiasID, materia.ID)
	}
	log.Printf("Materias IDs to insert: %v", materiasID)

	//materia
	newID, err := app.DB.InserirAula(aula)
	if err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "problem making request"})
	}

	err = app.DB.AtualizarMateria(newID, materiasID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error"})
	}

	resp := JSONResponse{
		Error:   false,
		Message: "matéria atualizada",
	}

	return c.JSON(http.StatusAccepted, resp)
}

func (app *application) AtualizarAula(c echo.Context) error {
	log.Printf("endpoint hit")
	var payload models.Aula

	err := c.Bind(&payload)
	log.Printf("Attempting to fetch aula with ID: %d\n", payload.ID)

	if err != nil {
		log.Printf("bind error: %v/n", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	var materiasID []int
	for _, materia := range payload.Materias {
		materiasID = append(materiasID, materia.ID)
	}

	aula, err := app.DB.UmaAula(payload.ID)
	if err != nil {
		log.Printf("db error: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	aula.Name = payload.Name
	aula.Size = payload.Size
	aula.Active = payload.Active
	aula.Review = payload.Review
	aula.UpdatedAt = time.Now()

	err = app.DB.AtualizarAula(*aula)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	err = app.DB.AtualizarMateria(aula.ID, materiasID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	resp := JSONResponse{
		Error:   false,
		Message: "aula atualizada",
	}
	return c.JSON(http.StatusAccepted, resp)
}

func (app *application) DeletarAula(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	err = app.DB.DeleteAula(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	resp := JSONResponse{
		Error:   false,
		Message: "movie deleted",
	}

	return c.JSON(http.StatusAccepted, resp)
}

func (app *application) AulasGraphQL(c echo.Context) error {
	//need to populate graph type with aulas
	aulas, _ := app.DB.TodaAula()

	//get query from request
	q, _ := io.ReadAll(c.Request().Body)
	query := string(q)
	//create new variable of type *graph.graph
	g := graph.New(aulas)
	//set query string on the variable
	g.QueryString = query
	//perform query
	resp, err := g.Query()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error"})
	}
	//send response

	return c.JSON(http.StatusOK, resp)
}
