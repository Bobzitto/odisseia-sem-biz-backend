package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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
