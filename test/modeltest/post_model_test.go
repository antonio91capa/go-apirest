package modeltest

import (
	"log"
	"testing"

	"github.com/antonio91capa/go-apirest/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllPosts(t *testing.T) {
	err := RefreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post Table: %v\n", err)
	}

	_, _, err = SeedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post table: %v\n", err)
	}

	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the posts: %v\n", err)
		return
	}

	assert.Equal(t, len(*posts), 2)
}

func TestSavePost(t *testing.T) {
	err := RefreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table: %v\n", err)
	}

	user, err := SeedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newPost := models.Post{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}

	savedPost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the post: %v\n", err)
		return
	}

	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)
}

func TestGetPostByID(t *testing.T) {
	err := RefreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refresing user and post table %v\n", err)
	}

	post, err := SeedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding table")
	}

	foundPost, err := postInstance.FindPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("This is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)
}

func TestUpdatePost(t *testing.T) {
	err := RefreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}

	post, err := SeedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding table")
	}

	postUpdate := models.Post{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@mail.com",
		AuthorID: post.AuthorID,
	}

	updatedPost, err := postUpdate.UpdatePost(server.DB)
	if err != nil {
		t.Errorf("This is the error updating the post: %v\n", err)
		return
	}

	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)
}

func TestDeletePost(t *testing.T) {
	err := RefreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refresing user and post table: %v\v", err)
	}

	post, err := SeedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding tables")
	}

	isDeleted, err := postInstance.DeletePost(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("This is the error updating the user: %v\n", err)
		return
	}

	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))

}
