package pubsub

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/mitchellh/mapstructure"
	"go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/pubsub".
func init() {
	modules.Register("k6/x/pubsub", new(PubSub))
}

// PubSub is the k6 extension for a Google Pub/Sub client.
// See https://cloud.google.com/pubsub/docs/overview
type PubSub struct{}

type publisherConf struct {
	ProjectID string
}

func (ps *PubSub) Publisher(config map[string]interface{}) *pubsub.Client {

	cnf := &publisherConf{}
	err := mapstructure.Decode(config, cnf)
	if err != nil {
		log.Fatalf("xk6-pubsub: unable to read publisher config: %v", err)
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cnf.ProjectID)

	if err != nil {
		log.Fatalf("xk6-pubsub: unable to initialise publisher")
	}

	return client
}

func (ps *PubSub) Publish(p *pubsub.Client, topic, msg string) error {
	ctx := context.Background()
	t := p.Topic(topic)
	r := t.Publish(
		ctx,
		&pubsub.Message{
			Data: []byte(msg),
		},
	)

	_, err := r.Get(ctx)
	if err != nil {
		ReportError(err, fmt.Sprintf("xk6-pubsub: unable to publish message: message was '%s', topic was '%s'", msg, topic))
		return err
	}

	return nil
}
