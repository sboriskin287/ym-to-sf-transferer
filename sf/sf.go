package sf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"ym-to-spotify-transferer/ym"
)

var searchUrl = "https://api.spotify.com/v1/search"
var likeUrl = "https://api.spotify.com/v1/me/tracks"

func search(track ym.Track) (*string, error) {
	req, err := http.NewRequest(http.MethodGet, searchUrl, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	q := "remaster%20track:" + track.Title
	if len(track.Albums) > 0 {
		q += "%20album:" + track.Albums[0].Title
	}
	if len(track.Artists) > 0 {
		q += "%20artist:" + track.Artists[0].Name
	}
	query.Set("q", q)
	query.Set("type", "track")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN_SF")))
	res, err := (&(http.Client{})).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode not OK, but %s", res.Status)
	}
	var body SearchRes
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	var id = ""
	if len(body.Tracks.Items) > 0 {
		id = body.Tracks.Items[0].Id
	}
	return &id, nil
}

func Like(tracks []ym.Track) error {
	ids := make([]string, 50)
	for i := range tracks {
		track := tracks[len(tracks)-i-1]
		id, err := search(track)
		if err != nil {
			return err
		}
		ids[i%50] = *id
		if i != 0 && i%50 == 0 || i == len(tracks)-1 {
			if err = like(ids); err != nil {
				return err
			}
		}
	}
	return nil
}

func like(ids []string) error {
	reqBody := LikeTracks{Ids: ids}
	reqBodyMarshall, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, likeUrl, bytes.NewReader(reqBodyMarshall))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN_SF")))
	req.Header.Set("Content-Type", "application/json")
	res, err := (&(http.Client{})).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("StatusCode not OK, but %s", res.Status)
	}
	return nil
}
