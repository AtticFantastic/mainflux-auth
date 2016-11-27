package services

import (
	"encoding/json"
	"strings"

	"github.com/mainflux/mainflux-auth/domain"
	nats "github.com/nats-io/go-nats"
)

// Subscribe subscribes to the events published on a given topic.
func Subscribe(url, topic string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}

	nc.Subscribe(topic, func(m *nats.Msg) {
		go func() {
			event := struct {
				Owner string `json:"owner"`
				Type  string `json:"type"`
				Id    string `json:"id"`
			}{}

			if err := json.Unmarshal(m.Data, &event); err != nil {
				return
			}

			if parts := strings.Split(event.Owner, " "); len(parts) == 2 {
				id, err := domain.SubjectOf(parts[1])
				if err != nil {
					return
				}

				AddResource(id, event.Type, event.Id)
			}
		}()
	})

	return nil
}
