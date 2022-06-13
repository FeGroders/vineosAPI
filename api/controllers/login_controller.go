package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fegroders/vineosAPI/api/auth"
	"github.com/fegroders/vineosAPI/api/models"
	"github.com/fegroders/vineosAPI/api/responses"
	"github.com/fegroders/vineosAPI/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, admin, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, struct {
		Token    string `json:"token"`
		Admin    bool   `json:"admin"`
	}{Token: token, Admin: admin})
}

func (server *Server) SignIn(email, password string) (string, bool, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", false, err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", false, err
	}
	// return auth.CreateToken(user.ID)
	token, err := auth.CreateToken(user.ID)

	return token, user.Admin, err
}

	