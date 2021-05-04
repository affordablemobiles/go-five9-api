package five9

import (
	"context"
	"fmt"
	"log"
	"strings"

	resty "github.com/go-resty/resty/v2"
)

func (api *Five9APIClient) performLogin(ctx context.Context, username, password string) (*Token, error) {
	client := resty.New()

	request, result := &Auth{
		Credentials: &PasswordCredentials{
			Username: username,
			Password: string(password),
		},
		AppKey: "web-ui",
		Policy: "ForceIn",
	}, &Token{}

	resp, err := client.SetCookieJar(api.cjar).R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		ForceContentType("application/json").
		SetResult(result).
		SetContext(ctx).
		Post("https://app.five9.com/appsvcs/rs/svc/auth/login")
	if err != nil {
		return nil, fmt.Errorf("Login call error: %s", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("Login error (%s): %s", resp.Status(), resp.String())
	}

	return resp.Result().(*Token), nil
}

func (api *Five9APIClient) handleStateChange(ctx context.Context) error {
	for i := 0; i < 5; i++ {
		if LoggingEnabled {
			log.Printf("Checking login state...")
		}
		state, err := api.getLoginState(ctx)
		if err != nil {
			return err
		}

		if LoggingEnabled {
			log.Printf("Handling login state, got %s", state)
		}
		switch state {
		case "WORKING":
			return nil
		case "SELECT_STATION":
			err := api.startSession(ctx)
			if err != nil {
				return err
			}
		case "ACCEPT_NOTICE":
			err := api.acceptNotices(ctx)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Invalid Login State: %s", state)
		}
	}

	return fmt.Errorf("Invalid Login State after 5 tries.. Exiting")
}

func (api *Five9APIClient) getLoginState(ctx context.Context) (string, error) {
	client := resty.New()

	resp, err := client.SetCookieJar(api.cjar).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer-%s", api.token.TokenID)).
		SetHeaderVerbatim("farmId", api.token.Context.FarmID).
		SetContext(ctx).
		Get(fmt.Sprintf("https://app.five9.com/appsvcs/rs/svc/agents/%s/login_state", api.token.UserID))
	if err != nil {
		return "", fmt.Errorf("Login State call error: %s", err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("Login State error (%s): %s", resp.Status(), resp.String())
	}

	return strings.Trim(resp.String(), "\""), nil
}

func (api *Five9APIClient) startSession(ctx context.Context) error {
	client := resty.New()

	request := &StationInfo{
		StationType: "EMPTY",
	}

	resp, err := client.SetCookieJar(api.cjar).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer-%s", api.token.TokenID)).
		SetHeaderVerbatim("farmId", api.token.Context.FarmID).
		SetBody(request).
		SetContext(ctx).
		Put(fmt.Sprintf("https://app.five9.com/appsvcs/rs/svc/agents/%s/session_start", api.token.UserID))
	if err != nil {
		return fmt.Errorf("Session Start call error: %s", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("Session Start error (%s): %s", resp.Status(), resp.String())
	}

	return nil
}

func (api *Five9APIClient) acceptNotices(ctx context.Context) error {
	notices, err := api.getNotices(ctx)
	if err != nil {
		return err
	}

	for _, notice := range notices {

		client := resty.New()

		resp, err := client.SetCookieJar(api.cjar).R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Bearer-%s", api.token.TokenID)).
			SetHeaderVerbatim("farmId", api.token.Context.FarmID).
			SetContext(ctx).
			Put(fmt.Sprintf("https://app.five9.com/appsvcs/rs/svc/agents/%s/maintenance_notices/%s/accept", api.token.UserID, notice.ID))
		if err != nil {
			return fmt.Errorf("Accept Notice call error: %s", err)
		}

		if !resp.IsSuccess() {
			return fmt.Errorf("Accept Notice error (%s): %s", resp.Status(), resp.String())
		}

	}

	return nil
}

func (api *Five9APIClient) getNotices(ctx context.Context) ([]*MaintenanceNoticeInfo, error) {
	client := resty.New()

	var nList []*MaintenanceNoticeInfo = []*MaintenanceNoticeInfo{}

	resp, err := client.SetCookieJar(api.cjar).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer-%s", api.token.TokenID)).
		SetHeaderVerbatim("farmId", api.token.Context.FarmID).
		ForceContentType("application/json").
		SetResult(nList).
		SetContext(ctx).
		Get(fmt.Sprintf("https://app.five9.com/appsvcs/rs/svc/agents/%s/maintenance_notices", api.token.UserID))
	if err != nil {
		return nil, fmt.Errorf("Error fetching list of notices (call): %s", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("Error fetching list of notices (%s): %s", resp.Status(), resp.String())
	}

	return *resp.Result().(*[]*MaintenanceNoticeInfo), nil
}
