package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var NoApiTokenError error = errors.New("error: client object has no api token")

type TinderClient struct {
	FbToken  string `json:"facebook_token"`
	FbId     int64  `json:"facebook_id"`
	apiToken string
}

type Recommendation struct {
	Id         string `json:"_id"`
	Bio        string `json:"bio"`
	BirthDate  string `json:"birth_date"`
	DistanceMi int    `json:"distance_mi"`
	Gender     int    `json:"gender"`
	Name       string `json:"name"`
	Photos     []struct {
		Id             string `json:"id"`
		Url            string `json:"url"`
		ProcessedFiles []struct {
			Height int    `json:"height"`
			Width  int    `json:"width"`
			Url    string `json:"url"`
		} `json:"processedFiles"`
	} `json:"photos"`
}

const apiRoot = "https://api.gotinder.com"

// NewTinderClient create new Tinder client.
func NewTinderClient(fbToken string, fbId int64) TinderClient {
	return TinderClient{FbToken: fbToken, FbId: fbId}
}

// Connect to Facebook to get a Tinder API key.
func (tc *TinderClient) Connect() error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(*tc)
	res, err := http.Post(fmt.Sprintf("%s/auth", apiRoot), "application/json", buf)
	if err != nil {
		return err
	}
	var response struct {
		User struct {
			ApiToken string `json:"api_token"`
		} `json:"user"`
	}
	json.NewDecoder(res.Body).Decode(&response)
	if response.User.ApiToken == "" {
		return NoApiTokenError
	}
	tc.apiToken = response.User.ApiToken
	return nil
}

var httpClient http.Client = http.Client{}

// Do an http get with auth headers.
func (tc *TinderClient) httpGet(url string) (*http.Response, error) {
	if tc.apiToken == "" {
		return nil, NoApiTokenError
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", tc.apiToken)
	return httpClient.Do(req)
}

// GetRecs gets recommendations.
func (tc *TinderClient) GetRecs() ([]Recommendation, error) {
	res, err := tc.httpGet(fmt.Sprintf("%s/user/recs", apiRoot))
	if err != nil {
		return nil, err
	}
	var response struct {
		Results []Recommendation `json:"results"`
	}
	json.NewDecoder(res.Body).Decode(&response)
	return response.Results, nil
}

// SwipeRight swipes a user right.
func (tc *TinderClient) SwipeRight(user *Recommendation) error {
	_, err := tc.httpGet(fmt.Sprintf("%s/like/%s", apiRoot, user.Id))
	return err
}

// SwipeLeft swipes a user right.
func (tc *TinderClient) SwipeLeft(user *Recommendation) error {
	_, err := tc.httpGet(fmt.Sprintf("%s/pass/%s", apiRoot, user.Id))
	return err
}
