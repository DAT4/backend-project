package main

import (
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
		input         Password
		expectedError bool
	}{
		{"", true},
		{"hej", true},
		{"M!rTin123HEJHEJHEJEHEJLL!!!###", true},
		{"Hej123!#", false},
	}

	for _, test := range tests {
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestUsernameValidation(t *testing.T) {
	var tests = []struct {
		input         Username
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"M!rTin123HEJHEJHEJEHEJLL!!!###", true},
		{"Martin", true},
		{"martin", false},
	}

	for _, test := range tests {
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestEmailValidation(t *testing.T) {
	var tests = []struct {
		input         Email
		expectedError bool
	}{
		{"", true},
		{"123", true},
		{"mail@google.dk", true},
		{"bhsi@dtu.dk", true},
		{"s195469@student.dtu.dk", false},
	}

	for _, test := range tests {
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestIpValidation(t *testing.T) {
	var tests = []struct {
		input         Ip
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
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestMacValidation(t *testing.T) {
	var tests = []struct {
		input         Mac
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
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}

func TestUserValidation(t *testing.T) {
	var tests = []struct {
		input         User
		expectedError bool
	}{
		{
			input: User{
				Id:       primitive.ObjectID{},
				Username: "martin",
				Password: "teSt123!#",
				Email:    "s123123@student.dtu.dk",
				Macs: []Mac{
					"00:01:e6:57:8b:68",
					"00:04:27:6a:5d:a1",
					"00:30:c1:5e:58:7d",
					"00:02:b3:bb:66:98",
				},
				Ips: []Ip{
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
			input: User{
				Id:       primitive.ObjectID{},
				Username: "",
				Password: "teSt123!#",
				Email:    "s123123@student.dtu.dk",
				Macs: []Mac{
					"00:01:e6:57:8b:68",
					"00:04:27:6a:5d:a1",
					"00:30:c1:5e:58:7d",
					"00:02:b3:bb:66:98",
				},
				Ips: []Ip{
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

	for _, test := range tests {
		err := test.input.validate()
		if (err != nil) != test.expectedError {
			t.Errorf("Expected %s got %s", exp(test.expectedError), got(err))
		}
	}
}
