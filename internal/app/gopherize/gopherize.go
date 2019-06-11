package gopherize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func GetImage(imageId string) ([]byte, error) {
	return httpGet(fmt.Sprintf("https://storage.googleapis.com/gopherizeme.appspot.com/gophers/%s.png", imageId))
}

func GetImageId(options []string) (string, error) {
	u := "https://gopherize.me/save?images=" + url.QueryEscape(strings.Join(options, "|"))
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPermanentRedirect {
		return "", fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}

	return filepath.Base(resp.Header.Get("location")), nil
}

func RandomOptions(artwork *Artwork) []string {
	var ret []string
	rand.Seed(time.Now().UnixNano())
	for _, v := range artwork.Categories {
		if len(v.Images) > 0 {
			if IsRequiredOption(v.Name) || rand.Intn(3) > 0 {
				ret = append(ret, v.Images[rand.Intn(len(v.Images))].Id)
			}
		}
	}
	return ret
}

func IsRequiredOption(name string) bool {
	return name == "Body" || name == "Eyes"
}

func GetArtwork() (*Artwork, error) {
	var artwork Artwork
	body, err := httpGet("https://gopherize.me/api/artwork/")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &artwork); err != nil {
		return nil, err
	}
	return &artwork, nil
}

func httpGet(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
