package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/fegroders/vineosAPI/api/auth"
	"github.com/fegroders/vineosAPI/api/models"
	"github.com/fegroders/vineosAPI/api/responses"
	"github.com/fegroders/vineosAPI/api/utils/formaterror"
)

func (server *Server) CreateWine(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	wine := models.Wine{}
	err = json.Unmarshal(body, &wine)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	wine.Prepare()
	err = wine.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	wineCreated, err := wine.SaveWine(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, wineCreated.ID))
	responses.JSON(w, http.StatusCreated, wineCreated)
}

func (server *Server) GetWines(w http.ResponseWriter, r *http.Request) {
	wine := models.Wine{}

	wines, err := wine.FindAllWines(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, wines)
}

func (server *Server) GetWine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	wine := models.Wine{}

	wineReceived, err := wine.FindWineByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, wineReceived)
}

func (server *Server) UpdateWine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the wine id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if the auth token is valid and  get the user id from it
	_, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the wine exist
	wine := models.Wine{}
	err = server.DB.Debug().Model(models.Wine{}).Where("id = ?", pid).Take(&wine).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Wine not found"))
		return
	}

	// Read the data wineed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	wineUpdate := models.Wine{}
	err = json.Unmarshal(body, &wineUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	wineUpdate.Prepare()
	err = wineUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	wineUpdate.ID = wine.ID //this is important to tell the model the wine id to update, the other update field are set above

	wineUpdated, err := wineUpdate.UpdateWine(server.DB)

	// if err != nil {
	// 	formattedError := formaterror.FormatError(err.Error())
	// 	responses.ERROR(w, http.StatusInternalServerError, formattedError)
	// 	return
	// }
	responses.JSON(w, http.StatusOK, wineUpdated)
}

func (server *Server) DeleteWine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Is a valid wine id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	_, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the wine exist
	wine := models.Wine{}
	err = server.DB.Debug().Model(models.Wine{}).Where("id = ?", pid).Take(&wine).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = wine.DeleteWine(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}