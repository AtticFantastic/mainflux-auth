package api_test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth/api"
	"github.com/mainflux/mainflux-auth/services"

	"gopkg.in/ory-am/dockertest.v2"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRedis(5, time.Second, func(url string) bool {
		services.StartCaching(url)
		return true
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	ts = httptest.NewServer(api.Server())
	defer ts.Close()

	result := m.Run()

	services.StopCaching()
	c.KillRemove()
	os.Exit(result)
}
