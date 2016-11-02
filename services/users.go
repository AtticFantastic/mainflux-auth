package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// RegisterUser invokes creation of new user account based on provided username
// and password.
func RegisterUser(username, password string) (domain.User, error) {
	var user domain.User

	if username == "" || password == "" {
		return user, &domain.ServiceError{Code: http.StatusBadRequest}
	}

	c := cache.Connection()
	defer c.Close()

	cVal, err := redis.Int64(c.Do("SADD", "users", username))
	if err != nil {
		return user, &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	if cVal == 0 {
		return user, &domain.ServiceError{Code: http.StatusConflict}
	}

	user, err = domain.CreateUser(username, password)
	if err != nil {
		return user, err
	}

	//
	// NOTE: consider using MULTI to ensure consistency
	//
	cKey := fmt.Sprintf("users:%s", user.Id)
	_, err = c.Do("HMSET", cKey, "username", user.Username, "password", user.Password, "masterKey", user.MasterKey)
	if err != nil {
		return user, &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return user, nil
}

// AddUserKey adds secondary user key. Bear in mind that any additional keys
// can be created only when identified as "master", i.e. by providing a master
// key.
func AddUserKey(userId, key string, access domain.AccessSpec) (string, error) {
	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("users:%s", userId)
	mKey, _ := redis.String(c.Do("HGET", cKey, "masterKey"))

	if mKey == "" {
		return "", &domain.ServiceError{Code: http.StatusNotFound}
	}

	if key != mKey {
		return "", &domain.ServiceError{Code: http.StatusForbidden}
	}

	if valid := access.Validate(); !valid {
		return "", &domain.ServiceError{Code: http.StatusBadRequest}
	}

	newKey, err := domain.CreateKey(userId, &access)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	cKey = fmt.Sprintf("users:%s:keys", userId)
	_, err = c.Do("SADD", cKey, newKey)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return newKey, nil
}
