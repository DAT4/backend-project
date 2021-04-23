package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	StatusCode int    `json:"-"`
}

func TestTokenHandler(t *testing.T) {
	//TODO look for dependency injection (testdatabase, timeout on jwt)
	var testDataList = []TestUser{
		{"martin", "T3stpass!", http.StatusOK},
		{"martin", "wrongpass!", http.StatusUnauthorized},
	}
	t.Run("Testing the login", func(t *testing.T) {
		for _, testData := range testDataList {
			body, _ := json.Marshal(testData)
			request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			response := httptest.NewRecorder()

			tokenHandler(response, request, CreateTestDB())

			if response.Code != testData.StatusCode {
				t.Errorf("Expected %v got %v", testData.StatusCode, response.Code)
			}
		}
	})
}
func TestCreateUser(t *testing.T) {
}
func TestRefreshToken(t *testing.T) {
}

func TestJoinWebsocketConnection(t *testing.T) {
	t.Run("Testing the login", func(t *testing.T) {
		user := TestUser{"martin", "T3stpass!", http.StatusOK}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		response := httptest.NewRecorder()

		tokenHandler(response, request, CreateTestDB())

		if response.Code != user.StatusCode {
			t.Errorf("Expected %v got %v", user.StatusCode, response.Code)
		}

		var token map[string]string

		_ = json.NewDecoder(response.Body).Decode(&token)

		request, _ = http.NewRequest(http.MethodGet, "/join", nil)
		tkn := fmt.Sprintf("Bearer %v", token["auth_token"])
		fmt.Println(tkn)
		response = httptest.NewRecorder()

		joinWebsocketConnection(response, request)

		if response.Code != user.StatusCode {
			t.Errorf("Expected %v got %v", user.StatusCode, response.Code)
		}
	})
}
func CreateTestDB() *dao.TestDB {
	db := &dao.TestDB{}
	users := []models.User{
		{
			Id:       primitive.NewObjectID(),
			PlayerID: 0,
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
		{
			Id:       primitive.NewObjectID(),
			PlayerID: 0,
			Username: "simon",
			Password: "hej",
			Email:    "simon@gmail.dk",
		},
	}
	for _, user := range users {
		_ = db.Create(&user)
	}
	return db
}
