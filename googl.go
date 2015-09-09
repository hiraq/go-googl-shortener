package googl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

// GoogleShorten as string type, so we can put our target url in construction method
type GoogleShorten struct {
	Target     string
	BaseURL    string
	Version    string
	ShortenURL string
}

// NewGoogleShorten acts as constructor to create initial values
func (gs *GoogleShorten) NewGoogleShorten(key, target, version string) GoogleShorten {
	gs.Target = target
	gs.Version = version
	gs.BaseURL = "https://www.googleapis.com/urlshortener/v" + version + "/url?key=" + key

	return *gs
}

// ShortIt create shorten url send request to google api
func (gs *GoogleShorten) ShortIt() error {

	jsonParam := []byte(fmt.Sprintf(`{"longUrl": "%v"}`, gs.Target))
	req, err := http.NewRequest("POST", gs.BaseURL, bytes.NewBuffer(jsonParam))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content, err := simplejson.NewJson(body)
	if err != nil {
		return err
	}

	jsonError := content.Get("error").Interface()
	if jsonError != nil {
		errorMap := jsonError.(map[string]interface{})
		errorStr := errorMap["message"]
		return errors.New(errorStr.(string))
	}

	gs.ShortenURL = content.Get("id").MustString()
	return nil
}

// String used return the value put in the construction type
func (gs *GoogleShorten) String() string {
	return gs.ShortenURL
}

func catchTheError(err error) {
	log.Fatalf("We have an error : %v", err)
}

// ShortIt is a main function to shorten given url
func ShortIt(key, url string) (string, error) {

	shortener := GoogleShorten{}
	googl := shortener.NewGoogleShorten(key, url, "1")
	err := googl.ShortIt()
	if err != nil {
		return "", err
	}

	return googl.String(), nil
}
