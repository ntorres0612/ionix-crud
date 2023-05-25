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

func (server *Server) CreateDrug(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	drug := &models.Drug{}
	err = json.Unmarshal(body, &drug)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	drug.Prepare()
	err = drug.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	drugCreated, err := drug.SaveDrug(server.DB)
	if err != nil {
		fmt.Println(err.Error())
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, drugCreated.ID))
	responses.JSON(w, http.StatusCreated, drugCreated)
}

func (server *Server) GetDrugs(w http.ResponseWriter, r *http.Request) {

	drug := models.Drug{}

	posts, err := drug.FindAllDrugs(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetDrug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	drug := models.Drug{}

	drugReceived, err := drug.FindDrugByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, drugReceived)
}

func (server *Server) UpdateDrug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	drug := models.Drug{}
	err = server.DB.Debug().Model(models.Drug{}).Where("id = ?", pid).Take(&drug).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("drug not found"))
		return
	}

	// Read the data posted
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	drugUpdate := models.Drug{}
	err = json.Unmarshal(body, &drugUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	drugUpdate.Prepare()
	err = drugUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	drugUpdate.ID = drug.ID

	postUpdated, err := drugUpdate.UpdateADrug(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeleteDrug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	drug := models.Drug{}
	err = server.DB.Debug().Model(models.Drug{}).Where("id = ?", id).Take(&drug).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = drug.DeleteADrug(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}
