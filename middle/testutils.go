package middle

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/dto"
	"log"
)

func AddUsersToTestDb(users []dto.User, db *dao.TestDB) {
	for _, user := range users {
		err := db.Insert(&user)
		if err != nil {
			log.Fatalf("Failed in function AddUserToTestDb to testDB %v", err)
		}
	}
}
