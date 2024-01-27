package spapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
)

type Client struct {
	ClientID     string
	ClientSecret string
	SellerID     string
	Token        *oauth2.Token
	HTTPClient   *http.Client
	Marketplace  *Marketplace
}

func NewClient() *Client {
	return &Client{}
}

var (
	EndpointNorthAmerica = "sellingpartnerapi-na.amazon.com"
	EndpointEurope       = "sellingpartnerapi-eu.amazon.com"
	EndpointFarEast      = "sellingpartnerapi-fe.amazon.com"
)

type amazonToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (s *Client) refreshToken(ctx context.Context) (*oauth2.Token, error) {
	if !s.Token.Expiry.IsZero() && s.Token.Expiry.After(time.Now()) {
		return s.Token, nil
	}

	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", s.Token.RefreshToken)
	body.Set("client_id", s.ClientID)
	body.Set("client_secret", s.ClientSecret)

	req, err := http.NewRequest(http.MethodPost, amazon.Endpoint.TokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error refreshing tokens: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var spapiErr Error
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading spapi response body: %w", err)
		}

		_ = json.Unmarshal(b, &spapiErr)
		spapiErr.StatusCode = res.StatusCode
		spapiErr.URL = req.URL
		return nil, spapiErr
	}

	var tr amazonToken
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("error decoding token response: %w", err)
	}

	s.Token = &oauth2.Token{
		AccessToken:  tr.AccessToken,
		TokenType:    tr.TokenType,
		RefreshToken: tr.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}

	return s.Token, nil
}

type retryOptions struct {
	retryLimit    int
	sleepDuration time.Duration
}

func (s *Client) retry(ctx context.Context, req request, opts retryOptions) (*http.Response, error) {
	token, err := s.refreshToken(ctx)
	if err != nil {
		return nil, err
	}

	if opts.sleepDuration == 0 {
		opts.sleepDuration = 1 * time.Second
	}

	var spapiErr Error

	for i := 0; i < opts.retryLimit; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		body := bytes.NewReader(req.Body)
		request, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), body)
		if err != nil {
			return nil, err
		}
		request.Header.Set("x-amz-access-token", token.AccessToken)

		res, err := s.HTTPClient.Do(request)
		if err != nil {
			return nil, err
		}

		if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
			return res, nil
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading spapi response body: %w", err)
		}

		_ = json.Unmarshal(b, &spapiErr)
		spapiErr.StatusCode = res.StatusCode
		spapiErr.URL = req.URL

		switch res.StatusCode {
		case http.StatusUnauthorized:
			// only one retry for unauthorized because the new token is valid since it was just refreshed.
			if spapiErr.StatusCode == http.StatusUnauthorized {
				return nil, spapiErr
			}

			creds, err := s.refreshToken(ctx)
			if err != nil {
				return nil, err
			}
			req.Token = creds
		case http.StatusTooManyRequests:
			if spapiErr.Errors[0].Code == "QuotaExceeded" {
				i = 0
				time.Sleep(30 * time.Second)
				continue
			}
			time.Sleep(opts.sleepDuration)
			continue
		default:
			return nil, spapiErr
		}
	}

	return nil, RetryError{
		RetryCount:    opts.retryLimit,
		Err:           spapiErr,
		SleepDuration: opts.sleepDuration,
	}
}

type request struct {
	Token  *oauth2.Token
	Method string
	URL    *url.URL
	Body   []byte
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Error struct {
	URL         *url.URL
	Body        []byte
	StatusCode  int
	Errors      []ResponseError `json:"errors"`
	Description string          `json:"error_description"`
	Msg         string          `json:"error"`
}

func (e Error) Error() string {
	errMsg := ""
	for _, err := range e.Errors {
		errMsg += fmt.Sprintf("%s: %s - %s\n", err.Code, err.Message, err.Details)
	}
	if errMsg == "" {
		errMsg = fmt.Sprintf("%s %s", e.Msg, e.Description)
	}
	return fmt.Sprintf("SPAPI API Error (Endpoint: %s - Status Code: %v) - Body: %s -- Message: %s", e.URL.String(), e.StatusCode, string(e.Body), errMsg)
}

type RetryError struct {
	RetryCount    int
	SleepDuration time.Duration
	Err           error
}

func (e RetryError) Error() string {
	return fmt.Sprintf("SPAPI Retry Error (Retry Count: %v, Sleep Duration: %v): %s", e.RetryCount, e.SleepDuration, e.Err.Error())
}
