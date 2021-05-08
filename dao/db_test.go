package dao

import (
	"fmt"
	"github.com/DAT4/backend-project/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

/*
 * THESE TEST FUNCTIONS REQUIRES DOCKER SETUP
 * WITH A COMMON NETWORK AND A MONGODB DOCKER
 * CONTAINER CALLED "testmongo"
 */

var (
	db MongoDB
)

var (
	port = "27017"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker:", err)
	}
	opts := dockertest.RunOptions{
		Hostname:     "dockertest",
		Name:         "dockertest",
		Repository:   "mongo",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {{HostIP: "0.0.0.0", HostPort: port}},
		},
	}

	fmt.Println("Creating resource")
	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatal("Could not start resource: ", err)
	}
	err = resource.Expire(60)
	if err != nil {
		log.Fatal("Could not set expire on resource: ", err)
	}

	if err = pool.Retry(func() error {
		fmt.Println("Trying")
		uri := fmt.Sprintf("mongo://localhost:%v", port)
		db, err = NewMongoDB(uri)
		return err
	}); err != nil {
		log.Fatal("Could not connect to docker: ", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge resource:", err)
	}

	fmt.Println("Done")

	os.Exit(code)
}

func TestMongoDB_Create(t *testing.T) {
	fmt.Println("Start Test Create User")
	t.Run("test create user in mongoDB", func(t *testing.T) {
		user := models.User{
			Username: "martin",
			Password: "T3stpass!",
			Email:    "mail@mama.sh",
		}
		user, err := db.Create(user)
		assert.NoError(t, err)
		fmt.Println("Create User OK")
	})
}

func TestMongoDB_UserFromName(t *testing.T) {
	fmt.Println("Start Test User From Name")
	user := models.User{
		Username: "martin",
		Password: "T3stpass!",
		Email:    "mail@mama.sh",
	}
	user, err := db.Create(user)
	assert.NoError(t, err)
	t.Run("Test authenticate user", func(t *testing.T) {
		_, err := db.UserFromName(string(user.Username))
		assert.NoError(t, err)
		fmt.Println("ok")
	})
	fmt.Println("User From DB OK")
}
