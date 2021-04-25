package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type TestUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	StatusCode int    `json:"-"`
}

func TestTokenHandler(t *testing.T) {
	users := []models.User{
		{
			Id:       primitive.NewObjectID(),
			PlayerID: 0,
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	}

	testDb := CreateTestDB(users)

	//TODO look for dependency injection (timeout on jwt)
	var testDataList = []TestUser{
		{"martin", "T3stpass!", http.StatusOK},
		{"martin", "wrongpass!", http.StatusUnauthorized},
	}
	t.Run("Testing the login", func(t *testing.T) {
		for _, testData := range testDataList {

			body, _ := json.Marshal(testData)
			request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			response := httptest.NewRecorder()

			tokenHandler(response, request, testDb)

			if response.Code != testData.StatusCode {
				t.Errorf("Expected %v got %v", testData.StatusCode, response.Code)
			}
		}
	})
}
func TestCreateUser(t *testing.T) {
	testDb := CreateTestDB([]models.User{})
	t.Run("Testing creating a user", func(t *testing.T) {

		newUser := models.User{
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		}

		body, _ := json.Marshal(newUser)
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		response := httptest.NewRecorder()

		createUser(response, request, testDb)

		if response.Code != http.StatusCreated {
			t.Errorf("Expected %v got %v", http.StatusCreated, response.Code)
		}
	})

}

func TestRefreshToken(t *testing.T) {
	testDb := CreateTestDB([]models.User{
		{
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	})

	t.Run("Testing creating a user", func(t *testing.T) {
		ok, tokens := assertLogin(models.User{Username: "martin", Password: "T3stpass!"}, testDb)
		if !ok {
			t.Error("could not login")
		}
		request, _ := http.NewRequest(http.MethodGet, "/register", nil)
		request.Header.Add("RefreshToken", tokens["refresh_token"])
		response := httptest.NewRecorder()

		refreshToken(response, request, testDb)

		if response.Code != http.StatusOK {
			t.Errorf("Expected %v got %v", http.StatusOK, response.Code)
		}
	})
}

func assertLogin(user models.User, db dao.DBase) (ok bool, token map[string]string) {
	body, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	response := httptest.NewRecorder()

	tokenHandler(response, request, db)

	if response.Code != http.StatusOK {
		return false, nil
	}
	_ = json.NewDecoder(response.Body).Decode(&token)
	return true, token
}

func TestJoinWebsocketConnection(t *testing.T) {

	id := primitive.NewObjectID()

	users := []models.User{
		{
			Id:       id,
			PlayerID: 0,
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	}

	testDb := CreateTestDB(users)

	go middle.G.Run(testDb)

	t.Run("Testing the websocket", func(t *testing.T) {

		s := httptest.NewServer(http.HandlerFunc(joinWebsocketConnection))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Error(err)
		}

		defer ws.Close()

		ok, token := assertLogin(users[0], testDb)

		if !ok {
			t.Error("login failed")
		}
		tkn := fmt.Sprintf("Bearer %v", token["auth_token"])
		msg := append([]byte{0, 0, 0, 0}, tkn...)
		err = ws.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			t.Error(err)
		}
		_, got, err := ws.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		expected := []byte{0, 1, 1, 1}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %v got %v", expected, got)
		}
	})
}

func CreateTestDB(users []models.User) *dao.TestDB {
	db := &dao.TestDB{}
	for _, user := range users {
		_ = db.Create(&user)
	}
	return db
}
