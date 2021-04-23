package middle

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func exp(err bool) string {
	if err {
		return "error"
	} else {
		return "nil"
	}
}

func got(err error) string {
	if err != nil {
		return err.Error()
	} else {
		return "nil"
	}
}

func TestPasswordValidation(t *testing.T) {
	var tests = []struct {
		input         models.Password
		expectedError bool
	}{
		{"", true},
		{"hej", true},
		{"M!rTin123HEJHEJHEJEHEJLL!!!###", true},
		{"Hej123!#", false},
	}

	for _, test := range tests {
		err := validatePassword(test.input)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestUsernameValidation(t *testing.T) {
	var tests = []struct {
		input         models.Username
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"M!rTin123HEJHEJHEJEHEJLL!!!###", true},
		{"Martin", true},
		{"martin", false},
	}

	for _, test := range tests {
		err := validateUsername(test.input)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestEmailValidation(t *testing.T) {
	var tests = []struct {
		input         models.Email
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"mail@google.dk", false},
		{"bhsi@dtu.dk", false},
		{"s195469@studenteqqw1231$#@!....dtu.dk", true},
	}

	for _, test := range tests {
		err := validateEmail(test.input)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestIpValidation(t *testing.T) {
	var tests = []struct {
		input         models.Ip
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"mail@google.dk", true},
		{"192.168.8", true},
		{"s195469@student.dtu.dk", true},
		{"192.168.0.1", false},
		{"192.168.0.3", false},
		{"192.168.0.5", false},
		{"192.168.0.6", false},
		{"192.168.0.12", false},
		{"192.168.0.13", false},
		{"192.168.0.87", false},
		{"192.168.0.90", false},
		{"192.168.0.105", false},
		{"192.168.0.153", false},
		{"192.168.0.191", false},
		{"192.168.0.251", false},
		{"192.168.0.196", false},
		{"185.107.12.169", false},
	}

	for _, test := range tests {
		err := validateIp(test.input)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestMacValidation(t *testing.T) {
	var tests = []struct {
		input         models.Mac
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"mail@google.dk", true},
		{"192.168.8", true},
		{"s195469@student.dtu.dk", true},
		{"00:c0:9f :09:b8:db", true},
		{"00:02:a5:90;c3:e6", true},
		{"00:c0:9f:0b:9$:d1", true},
		{"00:02:b3:46:0d:4c", false},
		{"00:02:a5:de:c2:17", false},
		{"00:0b:db:b2:fa:60", false},
		{"00:02:b3:06:d7:9b", false},
		{"00:13:72:09:ad:76", false},
		{"00:10:db:26:4d:52", false},
		{"00:01:e6:57:8b:68", false},
		{"00:04:27:6a:5d:a1", false},
		{"00:30:c1:5e:58:7d", false},
		{"00:02:b3:bb:66:98", false},
	}

	for _, test := range tests {
		err := validateMac(test.input)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestUserValidation(t *testing.T) {
	var tests = []struct {
		input         models.User
		expectedError bool
	}{
		{
			input: models.User{
				Id:       primitive.ObjectID{},
				Username: "martini",
				Password: "teSt123!#",
				Email:    "s123123@studentdtu.dk",
				Macs: []models.Mac{
					"00:01:e6:57:8b:68",
					"00:04:27:6a:5d:a1",
					"00:30:c1:5e:58:7d",
					"00:02:b3:bb:66:98",
				},
				Ips: []models.Ip{
					"192.168.0.105",
					"192.168.0.153",
					"192.168.0.191",
					"192.168.0.251",
					"192.168.0.196",
				},
			},
			expectedError: false,
		},
		{
			input: models.User{
				Id:       primitive.ObjectID{},
				Username: "",
				Password: "teSt123!#",
				Email:    "s123123@studedtu.dk",
				Macs: []models.Mac{
					"00:01:e6:57:8b:68",
					"00:04:27:6a:5d:a1",
					"00:30:c1:5e:58:7d",
					"00:02:b3:bb:66:98",
				},
				Ips: []models.Ip{
					"192.168.0.105",
					"192.168.0.153",
					"192.168.0.191",
					"192.168.0.251",
					"192.168.0.196",
				},
			},
			expectedError: true,
		},
	}

	db := CreateTestDB()

	for _, test := range tests {
		err := Validate(test.input, &db)
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func CreateTestDB() dao.TestDB {
	db := dao.TestDB{}
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
