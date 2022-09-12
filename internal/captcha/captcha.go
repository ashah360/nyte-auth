package captcha

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type RecaptchaService interface {
	ValidateRecaptcha(response string) (bool, error)
}

type GoogleRecaptchaResponse struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	ErrorCodes         []string `json:"error-codes"`
}

type recaptchaService struct {
	VerifyURL string
	ClientKey string
	ServerKey string
}

func (rs *recaptchaService) ValidateRecaptcha(response string) (bool, error) {
	res, err := http.PostForm(rs.VerifyURL, url.Values{
		"secret":   {rs.ServerKey},
		"response": {response},
	})
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var grc GoogleRecaptchaResponse

	err = json.Unmarshal(body, &grc)
	if err != nil {
		return false, err
	}

	return grc.Success, nil
}

func NewRecaptchaService(verifyURL, clientKey, serverKey string) RecaptchaService {
	return &recaptchaService{
		verifyURL, clientKey, serverKey,
	}
}
