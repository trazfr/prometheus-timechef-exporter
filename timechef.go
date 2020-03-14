package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	oauthAuthenticateURL = "https://timechef.elior.com/api/oauth/?scope=timechef" // POST
	oauthMeURL           = "https://timechef.elior.com/api/oauth/me"              // GET
	oauthRefreshURL      = "https://timechef.elior.com/api/oauth/refresh"         // POST
	soldeURL             = "https://timechef.elior.com/api/convive/%s/solde"      // GET
)

type TimechefResults struct {
	Site  string
	Solde float64
}

type soldeResponse struct {
	SiteName string  `json:"siteName"`
	Solde    float64 `json:"solde"`
}

type meResponse struct {
	Sites []meResponseSite `json:"sites"`
}

type meResponseSite struct {
	Name string `json:"name"`
}

type oauthRequest struct {
	AuthType string `json:"authType"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type oauthResponse struct {
	AccessToken  string `json:"accessToken"`
	Expires      string `json:"expires"`
	RefreshToken string `json:"refreshToken"`
}

type oauthRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type oauthData struct {
	accessToken  string
	refreshToken string
	expires      time.Time
}

type TimechefFetcher struct {
	client             *http.Client
	refreshTokenMargin time.Duration
	fetchURL           string
	oauthData          *oauthData
}

func NewTimecheFetcher(client *http.Client, username, password string) (*TimechefFetcher, error) {
	result := &TimechefFetcher{
		client:             client,
		refreshTokenMargin: client.Timeout,
	}

	return result, result.callAuthenticate(username, password)
}

func (t *TimechefFetcher) Fetch() (TimechefResults, error) {
	result := TimechefResults{}
	if t.oauthData.expires.Before(time.Now()) {
		if err := t.callRefreshToken(); err != nil {
			return result, err
		}
	}

	response := soldeResponse{}
	if err := t.get(t.fetchURL, &response, t.getBearerToken()); err != nil {
		return result, err
	}

	result.Site = response.SiteName
	result.Solde = response.Solde
	return result, nil
}

func (t *TimechefFetcher) callAuthenticate(username, password string) error {
	log.Printf("Authenticate using %s\n", username)

	request := oauthRequest{
		Username: username,
		Password: password,
	}
	response := oauthResponse{}
	if err := t.post(oauthAuthenticateURL, &request, &response, ""); err != nil {
		return err
	}

	if err := t.updateTokenFromResponse(&response); err != nil {
		return err
	}
	return t.getSite()
}

func (t *TimechefFetcher) callRefreshToken() error {
	log.Print("Refresh token\n")

	request := oauthRefreshRequest{
		RefreshToken: t.oauthData.refreshToken,
	}
	response := oauthResponse{}
	if err := t.post(oauthRefreshURL, &request, &response, "null"); err != nil {
		return err
	}
	return t.updateTokenFromResponse(&response)
}

func (t *TimechefFetcher) getSite() error {
	response := meResponse{}
	if err := t.get(oauthMeURL, &response, t.getBearerToken()); err != nil {
		return err
	}

	if len(response.Sites) == 0 {
		return errors.New("No site")
	}
	site := response.Sites[0].Name
	log.Printf("Got site %s\n", site)
	t.fetchURL = fmt.Sprintf(soldeURL, site)
	return nil
}

func (t *TimechefFetcher) updateTokenFromResponse(response *oauthResponse) error {
	expires, err := time.Parse("2006-01-02T15:04:05.9999999Z", response.Expires)
	if err != nil {
		return err
	}
	t.oauthData = &oauthData{
		accessToken:  response.AccessToken,
		refreshToken: response.RefreshToken,
		expires:      expires.Add(-t.refreshTokenMargin),
	}
	log.Printf("Got token. Refresh after: %v\n", expires)
	return nil
}

func (t *TimechefFetcher) getBearerToken() string {
	return fmt.Sprintf("Bearer %s", t.oauthData.accessToken)
}

func (t *TimechefFetcher) get(url string, response interface{}, authorization string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return t.do(req, response, authorization)
}

func (t *TimechefFetcher) post(url string, request, response interface{}, authorization string) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return t.do(req, response, authorization)
}

func (t *TimechefFetcher) do(req *http.Request, response interface{}, authorization string) error {
	req.Header.Set("Accept", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}
