package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/dto"
	"github.com/DAT4/backend-project/middle"
)

type TestUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	StatusCode int    `json:"-"`
}

func TestTokenHandler(t *testing.T) {
	users := []dto.User{
		{
			Id:       primitive.NewObjectID().Hex(),
			PlayerID: 0,
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	}

	testDb := dao.NewTestDB()
	server := API{
		Game: middle.NewGame(testDb),
	}
	middle.AddUsersToTestDb(users, testDb)

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

			server.TokenHandler(response, request)

			if response.Code != testData.StatusCode {
				t.Errorf("Expected %v got %v", testData.StatusCode, response.Code)
			}
		}
	})
}

func TestCreateUser(t *testing.T) {
	server := API{middle.NewGame(dao.NewTestDB())}
	t.Run("Testing creating a user", func(t *testing.T) {

		newUser := dto.User{
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		}

		body, _ := json.Marshal(newUser)
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		response := httptest.NewRecorder()

		server.InsertUser(response, request)

		if response.Code != http.StatusCreated {
			t.Errorf("Expected %v got %v", http.StatusCreated, response.Code)
		}
	})

}

func TestRefreshToken(t *testing.T) {
	db := dao.NewTestDB()
	server := API{middle.NewGame(db)}

	users := []dto.User{
		{
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	}

	middle.AddUsersToTestDb(users, db)

	t.Run("Testing creating a user", func(t *testing.T) {
		ok, tokens := assertLogin(dto.User{Username: "martin", Password: "T3stpass!"}, server)
		if !ok {
			t.Error("could not login")
		}
		request, _ := http.NewRequest(http.MethodGet, "/register", nil)
		request.Header.Add("RefreshToken", tokens["refresh_token"])
		response := httptest.NewRecorder()

		server.RefreshToken(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Expected %v got %v", http.StatusOK, response.Code)
		}
	})
}

func assertLogin(user dto.User, server API) (ok bool, token map[string]string) {
	body, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	response := httptest.NewRecorder()

	server.TokenHandler(response, request)

	if response.Code != http.StatusOK {
		return false, nil
	}
	_ = json.NewDecoder(response.Body).Decode(&token)
	return true, token
}

func TestJoinWebsocketConnection(t *testing.T) {

	testDb := dao.NewTestDB()
	id := primitive.NewObjectID().Hex()

	users := []dto.User{
		{
			Id:       id,
			PlayerID: 0,
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		},
	}

	middle.AddUsersToTestDb(users, testDb)

	server := API{
		Game: middle.NewGame(testDb),
	}

	go server.Game.Run()

	t.Run("Testing the websocket", func(t *testing.T) {

		s := httptest.NewServer(http.HandlerFunc(server.JoinWebsocketConnection))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Error(err)
		}

		defer ws.Close()

		ok, token := assertLogin(users[0], server)

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
