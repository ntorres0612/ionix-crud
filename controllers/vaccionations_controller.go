package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ntorres0612/ionix-crud/api/auth"
	"github.com/ntorres0612/ionix-crud/api/utils/formaterror"
	"github.com/ntorres0612/ionix-crud/models"
	"github.com/ntorres0612/ionix-crud/responses"
)

func (server *Server) CreateVaccination(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	vaccination := models.Vaccination{}
	err = json.Unmarshal(body, &vaccination)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	vaccination.Prepare()
	err = vaccination.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postCreated, err := vaccination.SaveVaccination(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetVaccinations(w http.ResponseWriter, r *http.Request) {

	vaccination := models.Vaccination{}

	vaccinations, err := vaccination.FindAllSaveVaccinations(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, vaccinations)
}

func (server *Server) GetVaccination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	vaccination := models.Vaccination{}

	postReceived, err := vaccination.FindVaccinationByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdateVaccination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	vaccination := models.Vaccination{}
	err = server.DB.Debug().Model(models.Vaccination{}).Where("id = ?", id).Take(&vaccination).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("vaccination not found "))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	vaccionationUpdate := models.Vaccination{}
	err = json.Unmarshal(body, &vaccionationUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	vaccionationUpdate.Prepare()
	err = vaccionationUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	vaccionationUpdate.ID = vaccination.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := vaccionationUpdate.UpdateAVaccination(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeleteVaccination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	vaccination := models.Vaccination{}
	err = server.DB.Debug().Model(models.Vaccination{}).Where("id = ?", id).Take(&vaccination).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = vaccination.DeleteAVaccination(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}
