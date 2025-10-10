package request

import "github.com/cockroachdb/errors"

type GoogleOAuthRequestBody struct {
	AccessToken string `json:"access_token"`
}

func (r GoogleOAuthRequestBody) Validate() error {
	if r.AccessToken == "" {
		return errors.New("access_token cannot be empty")
	}

	return nil
}
