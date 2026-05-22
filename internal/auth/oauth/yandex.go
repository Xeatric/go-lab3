package oauth

import (
	"encoding/json"

	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// --- Константы эндпоинтов Yandex API ---
const (
	yandexAuthURL  = "https://oauth.yandex.ru/authorize"
	yandexTokenURL = "https://oauth.yandex.ru/token"
	yandexUserURL  = "https://login.yandex.ru/info"
)

type YandexOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// YandexTokenResponse ответ от сервера Яндекса с токенами
type YandexTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// YandexUserInfo структура данных пользователя от Яндекса
type YandexUserInfo struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"default_email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RealName  string `json:"real_name"`
}

func NewYandexOAuth(clientID, clientSecret, redirectURL string) *YandexOAuth {
	return &YandexOAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

// GetAuthURL формирует URL для редиректа пользователя на Яндекс
func (y *YandexOAuth) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", y.ClientID)
	params.Add("redirect_uri", y.RedirectURL)
	params.Add("state", state)

	return yandexAuthURL + "?" + params.Encode()
}

// ExchangeCode обменивает временный код на токен доступа
func (y *YandexOAuth) ExchangeCode(code string) (*YandexTokenResponse, error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("client_id", y.ClientID)
	data.Add("client_secret", y.ClientSecret)

	req, err := http.NewRequest("POST", yandexTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp YandexTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetUserInfo получает информацию о пользователе по Access Token'у
func (y *YandexOAuth) GetUserInfo(accessToken string) (*YandexUserInfo, error) {
	req, err := http.NewRequest("GET", yandexUserURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo YandexUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
