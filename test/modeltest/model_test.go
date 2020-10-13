package modeltest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/antonio91capa/go-apirest/api/controllers"
	"github.com/antonio91capa/go-apirest/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting the env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDbDriver")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func RefreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed")
	return nil
}

// Test insertart un usuario
func SeedOneUser() (models.User, error) {
	RefreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@mail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Cannot seed user table: %v", err)
	}
	return user, nil
}

// Test insertar varios usuarios
func SeedUsers() error {
	users := []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@mail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Carolina Pliego",
			Email:    "carolinap@mail.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Test Reiniciar table User y Post
func RefreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

// Test guardar un usuario con un post
func SeedOneUserAndOnePost() (models.Post, error) {
	err := RefreshUserAndPostTable()
	if err != nil {
		return models.Post{}, err
	}

	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@mail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}

	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

// Test save Users and Posts
func SeedUsersAndPosts() ([]models.User, []models.Post, error) {
	var err error
	if err != nil {
		return []models.User{}, []models.Post{}, err
	}

	var users = []models.User{
		models.User{
			Nickname: "Johan Bill",
			Email:    "johan@mail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Magu Fintech",
			Email:    "magu@mail.com",
			Password: "password",
		},
	}

	var posts = []models.Post{
		models.Post{
			Title:   "Title 1",
			Content: "Content 1",
		},
		models.Post{
			Title:   "Title 2",
			Content: "Content 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed user table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed post table: %v", err)
		}
	}
	return users, posts, nil
}
