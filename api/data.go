package dummy

import (
	"log"

	"github.com/antonio91capa/go-apirest/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "Steve victor",
		Email:    "steven@mail.com",
		Password: "p@ssw0rd",
	},
	models.User{
		Nickname: "Alex Morgan",
		Email:    "alex@mail.com",
		Password: "p@ssM",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Title Number 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Title Number 2",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("autho_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot dummy users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot dummy posts table: %v", err)
		}
	}
}
