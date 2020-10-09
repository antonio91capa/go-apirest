package controllers

import (
	"net/http"

	"github.com/antonio91capa/go-apirest/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.ResponseJSON(w, http.StatusOK, "Welcome To This API")
}
