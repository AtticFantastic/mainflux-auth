package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth/services"

	dockertest "gopkg.in/ory-am/dockertest.v2"
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRedis(5, time.Second, func(url string) bool {
		services.StartCaching(url)
		return true
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	result := m.Run()

	services.StopCaching()
	c.KillRemove()
	os.Exit(result)
}
