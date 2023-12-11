package ym

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var loginUrl = "https://login.yandex.ru/info?format=json"
var ymLikeIdsUrl = "https://api.music.yandex.net/users/%s/likes/tracks"
var ymLikesUrl = "https://api.music.yandex.net/tracks"

func getMyYmId() (*string, error) {
	req, err := http.NewRequest(http.MethodGet, loginUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN_YM"))
	res, err := (&(http.Client{})).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode not OK, but %s", res.Status)
	}
	var body map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%v", body["id"])
	return &id, nil
}

func getMyTrackIds(myId string) ([]TrackWithAlbum, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ymLikeIdsUrl, myId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN_YM"))
	res, err := (&(http.Client{})).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode not OK, but %s", res.Status)
	}
	var body TrackIdRes
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body.Result.Library.Tracks, nil
}

func GetMyTracks() ([]Track, error) {
	id, err := getMyYmId()
	if err != nil {
		return nil, err
	}
	trackIds, err := getMyTrackIds(*id)
	if err != nil {
		return nil, err
	}
	var trackIdsStr string
	for i, t := range trackIds {
		trackIdsStr += t.Id + ":" + t.AlbumId
		if i != len(trackIds)-1 {
			trackIdsStr += ","
		}
	}
	data := url.Values{}
	data.Set("track-ids", trackIdsStr)
	req, err := http.NewRequest(http.MethodPost, ymLikesUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN_YM"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := (&(http.Client{})).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode not OK, but %s", res.Status)
	}
	var body TrackRes
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body.Result, nil
}
