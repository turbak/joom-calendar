package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/auth"
	"net/http"
)

type Client struct {
	client *http.Client

	clientID     string
	clientSecret string
}

func NewClient(clientID, clientSecret string) *Client {
	c := &http.Client{}
	return &Client{
		client:       c,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type GetAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (c *Client) GetAccessToken(code string) (string, error) {
	body, _ := json.Marshal(GetAccessTokenRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Code:         code,
	})

	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenResp := &GetAccessTokenResponse{}
	if err = json.NewDecoder(resp.Body).Decode(tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

type GetUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *Client) GetUser(accessToken string) (*auth.User, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user")
	}

	userInfo := &GetUserResponse{}
	if err = json.NewDecoder(resp.Body).Decode(userInfo); err != nil {
		return nil, err
	}

	return &auth.User{
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}, nil
}
