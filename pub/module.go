package pub

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Module struct {
	ok    bool
	ctx   context.Context
	topic *pubsub.Topic
}

func (m *Module) Load() error {
	credentialJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialJSON == "" {
		return nil
	}
	jwtConfig, err := google.JWTConfigFromJSON([]byte(credentialJSON), pubsub.ScopePubSub)
	if err != nil {
		return nil
	}

	m.ctx = context.Background()
	c, err := pubsub.NewClient(m.ctx, "sable-home", option.WithTokenSource(jwtConfig.TokenSource(m.ctx)))
	if err != nil {
		return nil
	}
	m.topic = c.Topic("projects/sable-home/topics/game-with-updated-price")
	m.ok = true
	return nil
}

func (m *Module) Pub(game string, price float64) {
	if !m.ok {
		return
	}
	m.topic.Publish(m.ctx, &pubsub.Message{
		Attributes: map[string]string{
			"name":  game,
			"price": fmt.Sprintf("%.2f", price),
		},
	})
}
