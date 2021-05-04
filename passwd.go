package five9

import (
	"context"
	"fmt"
	"log"

	resty "github.com/go-resty/resty/v2"
)

func (api *Five9APIClient) ChangePassword(ctx context.Context, oldPassword, newPassword string) error {
	client := resty.New()

	request := &ChangePasswordInfo{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	if LoggingEnabled {
		log.Printf("Changing password...")
	}

	resp, err := client.SetCookieJar(api.cjar).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer-%s", api.token.TokenID)).
		SetHeaderVerbatim("farmId", api.token.Context.FarmID).
		SetBody(request).
		SetContext(ctx).
		Put(fmt.Sprintf("https://app.five9.com/appsvcs/rs/svc/agents/%s/password", api.token.UserID))
	if err != nil {
		return fmt.Errorf("Password change call error: %s", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("Password change error (%s): %s", resp.Status(), resp.String())
	}

	return nil
}
