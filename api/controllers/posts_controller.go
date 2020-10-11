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

// ************** Create New Post
func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := exception.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.ResponseJSON(w, http.StatusCreated, postCreated)
}

// ************************* Get All Posts
func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	responses.ResponseJSON(w, http.StatusOK, posts)
}

// ************************** Get Post By ID
func (server *Server) GetPostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}
	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.ResponseJSON(w, http.StatusOK, postReceived)
}

// **************************** Update Post
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// Check if the auth token is valid and get the user id
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id=?", pid).Take(&post).Error
	if err != nil {
		responses.Error(w, http.StatusNotFound, errors.New("Post Not Found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdatePost(server.DB)
	if err != nil {
		formattedError := exception.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.ResponseJSON(w, http.StatusOK, postUpdated)
}

// ********************************* Delete Post
func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Is a valid post id given to us
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id=?", pid).Take(&post).Error
	if err != nil {
		responses.Error(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = post.DeletePost(server.DB, pid, uid)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.ResponseJSON(w, http.StatusNoContent, "")
}
