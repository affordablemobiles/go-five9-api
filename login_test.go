package five9

import (
	"context"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	LoggingEnabled = true

	_, err := NewFive9APIClient(
		context.Background(),
		os.Getenv("FIVE9_USERNAME"),
		os.Getenv("FIVE9_PASSWORD"),
	)
	if err != nil {
		t.Errorf("Failed to perform login: %s", err)
		return
	}
}
