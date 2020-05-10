package unsplash

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const apiUrl = "https://api.unsplash.com/photos/random?collections=%s&orientation=landscape&count=1"

type Client struct {
	http        *http.Client
	collections string
	token       string
}

func New(token string, collections []int) *Client {
	var colls []string
	for _, c := range collections {
		colls = append(colls, fmt.Sprintf("%d", c))
	}
	return &Client{
		token:       token,
		http:        &http.Client{},
		collections: strings.Join(colls, ","),
	}
}

type imageData struct {
	Links struct {
		Html string `json:"html"`
	} `json:"links"`
	Urls struct {
		Raw string `json:"raw"`
	} `json:"urls"`
}

func (client *Client) Random() (im []byte, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf(apiUrl, client.collections)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", client.token))

	r, err := client.http.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	var imgs []imageData

	err = json.NewDecoder(r.Body).Decode(&imgs)
	if err != nil {
		return
	}

	if len(imgs) < 1 {
		err = errors.New("API returned no data")
		return
	}

	log.Printf("using %v", imgs[0].Links.Html)

	if imgs[0].Urls.Raw == "" {
		err = errors.New("API returned no raw url")
		return
	}

	rr, err := client.http.Get(imgs[0].Urls.Raw)
	if err != nil {
		return
	}

	defer func() {
		if cerr := rr.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	im, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		return
	}

	return
}
