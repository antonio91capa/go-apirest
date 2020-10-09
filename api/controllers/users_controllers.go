package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/antonio91capa/go-apirest/api/auth"
	"github.com/antonio91capa/go-apirest/api/exception"
	"github.com/antonio91capa/go-apirest/api/models"
	"github.com/antonio91capa/go-apirest/api/responses"
	"github.com/gorilla/mux"
)

// ---------------------- Create a new User
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		formattedError := exception.FormatError(err.Error())

		responses.Error(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.ResponseJSON(w, http.StatusCreated, userCreated)
}

// ---------------------- Get all users
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	responses.ResponseJSON(w, http.StatusOK, users)
}

// -------------------------- Get user by ID
func (server *Server) GEtUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	getUser, err := user.FindUserById(server.DB, uint32(uid))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	responses.ResponseJSON(w, http.StatusOK, getUser)
}

// ------------------------- Update User
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := exception.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.ResponseJSON(w, http.StatusOK, updatedUser)
}

// -------------------------- Delete User
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != 0 && tokenID != uint32(uid) {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteUser(server.DB, uint32(uid))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.ResponseJSON(w, http.StatusNoContent, "")

}
