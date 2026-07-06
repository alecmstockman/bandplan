package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const iTunesBaseURL = "https://itunes.apple.com/search"

func ITunesSearchByArtist(term string) ([]byte, error) {
	log.Println("- ITunesSearchSongs")

	req, err := http.NewRequest("GET", iTunesBaseURL, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("term", term)
	q.Set("media", "music")
	q.Set("entity", "song")
	q.Set("limit", "20")
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

	fmt.Println("Request URL:", req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes api returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("================== BODY ==================")
	// fmt.Println(string(body))

	return body, nil
}
