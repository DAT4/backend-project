package middle

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"log"
)

func AddUsersToTestDb(users []models.User, db *dao.TestDB) {
	for _, user := range users {
		_, err := CreateUser(user, db)
		if err != nil {
			log.Fatalf("Failed in function AddUserToTestDb to testDB %v", err)
		}
	}
}
