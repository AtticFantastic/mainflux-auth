package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// CheckPermissions determines whether or not a platform request can be
// performed given the provided key.
func CheckPermissions(key string, req domain.AccessRequest) error {
	if valid := req.Validate(); !valid {
		fmt.Println("Invalid :(")
		return &domain.AuthError{Code: http.StatusForbidden}
	}

	subject, err := domain.SubjectOf(key)
	if err != nil {
		return err
	}

	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("auth:%s:%s:master", domain.UserType, subject)
	masterKey, _ := redis.String(c.Do("GET", cKey))
	if req.Restricted() && key != masterKey {
		return &domain.AuthError{Code: http.StatusForbidden}
	}

	if key == masterKey {
		if req.Restricted() {
			return nil
		}

		cKey = fmt.Sprintf("auth:%s:%s:owned", domain.UserType, subject)
		res := fmt.Sprintf("%s:%s", req.Type, req.Id)
		if owner, _ := redis.Bool(c.Do("SISMEMBER", cKey, res)); !owner {
			return &domain.AuthError{Code: http.StatusForbidden}
		}

		return nil
	}

	wList := fmt.Sprintf("auth:%s:%s:%s", req.Type, req.Id, req.Action)
	if allowed, _ := redis.Bool(c.Do("SISMEMBER", wList, key)); !allowed {
		return &domain.AuthError{Code: http.StatusForbidden}
	}

	return nil
}
