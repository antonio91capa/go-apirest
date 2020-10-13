package modeltest

import (
	"log"
	"testing"

	"github.com/antonio91capa/go-apirest/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {
	err := RefreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = SeedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := RefreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID:       1,
		Email:    "test@mail.com",
		Nickname: "test",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
}

func TestGetUserByID(t *testing.T) {
	err := RefreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := SeedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v\n", err)
	}

	foundUser, err := userInstance.FindUserById(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}

func TestUpdateUser(t *testing.T) {
	err := RefreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := SeedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	userUpdate := models.User{
		ID:       1,
		Nickname: "modiUpdate",
		Email:    "modiupdate@mail.com",
		Password: "password",
	}

	updatedUser, err := userUpdate.UpdateUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("This is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

func TestDeleteUser(t *testing.T) {
	err := RefreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := SeedOneUser()
	if err != nil {
		log.Fatalf("cannot seed user: %v\n", err)
	}

	isDeleted, err := userInstance.DeleteUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("This is the error updating the user: %v\n", err)
		return
	}

	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}