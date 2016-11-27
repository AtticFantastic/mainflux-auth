package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/domain"
)

const maxIdle int = 10

var cache *redis.Pool

// StartCaching creates new redis connection pool with fixed number of allowed
// inactive connections.
func StartCaching(url string) {
	cache = &redis.Pool{
		MaxIdle: maxIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

// StopCaching terminates the redis pool.
func StopCaching() error {
	return cache.Close()
}

// Subscribe subscribes to the events published on a given topic.
func Subscribe(topic string) error {
	c := cache.Get()
	psc := redis.PubSubConn{c}
	defer psc.Close()

	if err := psc.Subscribe(topic); err != nil {
		return err
	}

	go func() {
		for {
			switch e := psc.Receive().(type) {
			case redis.Message:
				go storeEvent(e.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", e.Channel, e.Kind, e.Count)
			case error:
				fmt.Printf("an error has occurred: %s\n", e)
				return
			}
		}
	}()

	return nil
}

func storeEvent(data []byte) {
	event := struct {
		Owner string `json:"owner"`
		Type  string `json:"type"`
		Id    string `json:"id"`
	}{}

	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	if parts := strings.Split(event.Owner, " "); len(parts) == 2 {
		id, err := domain.SubjectOf(parts[1])
		if err != nil {
			return
		}

		AddResource(id, event.Type, event.Id)
	}
}
