package services_test

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddResource(t *testing.T) {
	user, resource, id := "user", domain.DevType, "id"
	services.AddResource(user, resource, id)

	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("auth:users:%s:owned", user)
	cVal := fmt.Sprintf("%s:%s", resource, id)
	if exists, _ := redis.Bool(c.Do("SISMEMBER", cKey, cVal)); !exists {
		t.Errorf("case 1: expected resource to be bound")
	}
}
