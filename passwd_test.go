package five9

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	pgen "github.com/sethvargo/go-password/password"
)

func TestPasswd(t *testing.T) {
	LoggingEnabled = true

	ctx := context.Background()

	api, err := NewFive9APIClient(
		ctx,
		os.Getenv("FIVE9_USERNAME"),
		os.Getenv("FIVE9_PASSWORD"),
	)
	if err != nil {
		t.Errorf("Failed to perform login: %s", err)
		return
	}

	pwd, err := pgen.Generate(18, 5, 0, false, false)
	if err != nil {
		t.Errorf("Failed to generate password: %s", err)
		return
	}
	pwd = fmt.Sprintf("%s__", pwd)

	log.Printf("Trying with password: %s", pwd)

	err = api.ChangePassword(
		ctx,
		os.Getenv("FIVE9_PASSWORD"),
		pwd,
	)
	if err != nil {
		t.Errorf("Failed to change password: %s", err)
		return
	}
}
