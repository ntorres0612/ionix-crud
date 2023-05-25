package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/ntorres0612/ionix-crud/models"
)

var users = []models.User{
	{
		Name:     "Nelson Torres",
		Email:    "ntorres@gmail.com",
		Password: "password",
	},
	{
		Name:     "Test",
		Email:    "test@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Drug{}, &models.Vaccination{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Vaccination{}).AddForeignKey("drug_id", "drugs(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}
}
