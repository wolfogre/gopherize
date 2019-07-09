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

var (
	artwork *Artwork
	verbose bool
)

func SetVerbose(v bool) {
	verbose = v
}

func RandomImage() ([]byte, string, error) {
	artwork, err := GetArtwork()
	if err != nil {
		return nil, "", err
	}
	options := RandomOptions(artwork)
	imageId, err := GetImageId(options)
	if err != nil {
		return nil, "", err
	}
	image, err := GetImage(imageId)
	if err != nil {
		return nil, "", err
	}
	return image, imageId + ".png", nil
}

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
	if verbose {
		fmt.Println("request:", u)
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
	if verbose {
		fmt.Println("options:", ret)
	}
	return ret
}

func IsRequiredOption(name string) bool {
	return name == "Body" || name == "Eyes"
}

func GetArtwork() (*Artwork, error) {
	if artwork != nil {
		return artwork, nil
	}

	body, err := httpGet("https://gopherize.me/api/artwork/")
	if err != nil {
		return nil, err
	}
	var temp Artwork
	if err := json.Unmarshal(body, &temp); err != nil {
		return nil, err
	} else {
		artwork = &temp
	}
	return artwork, nil
}

func httpGet(u string) ([]byte, error) {
	if verbose {
		fmt.Println("request:", u)
	}
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
