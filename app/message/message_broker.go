package message

import (
	"api/app/config"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func Publish(w io.Writer, projectID, topicID, msg string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"

	jsonByte := []byte(fmt.Sprintf(`{
		"type": "%s",
		"project_id": "%s",
		"private_key_id": "%s",
		"private_key": "%s",
		"client_email": "%s",
		"client_id": "%s",
		"auth_uri": "%s",
		"token_uri": "%s",
		"auth_provider_x509_cert_url": "%s",
		"client_x509_cert_url": "%s"
	  }`,
		config.GOOGLE_TYPE,
		projectID,
		config.GOOGLE_PRIVATE_KEY_ID,
		config.GOOGLE_PRIVATE_KEY,
		config.GOOGLE_CLIENT_EMAIL,
		config.GOOGLE_CLIENT_ID,
		config.GOOGLE_AUTH_URI,
		config.GOOGLE_TOKEN_URI,
		config.GOOGLE_AUTH_PROVIDER_x509_CERT_URL,
		config.GOOGLE_CLIENT_x509_CERT_URL,
	))

	option := option.WithCredentialsJSON(jsonByte)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, errGet := result.Get(ctx)
	if errGet != nil {
		return fmt.Errorf("get: %v", err)
	}
	return nil
}
