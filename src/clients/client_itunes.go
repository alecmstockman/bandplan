package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const iTunesBaseURL = "https://itunes.apple.com/search"
const iTunesLookupURL = "https://itunes.apple.com/lookup"

func ITunesSearchByArtist(term string) ([]byte, error) {
	log.Println("- ITunesSearchByArtist")

	req, err := http.NewRequest("GET", iTunesBaseURL, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("term", term)
	q.Set("media", "music")
	q.Set("entity", "song")
	q.Set("limit", "10")
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

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

	return body, nil
}

func ITunesSearchByArtistAndSong(term string) ([]byte, error) {

	req, err := http.NewRequest("GET", iTunesBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create iTunes artist and song search request: %w", err)
	}

	q := url.Values{}
	q.Set("term", term)
	q.Set("media", "music")
	q.Set("entity", "song")
	q.Set("limit", "25")
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send iTunes artist and song search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes artist and song search returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read iTunes artist and song search response: %w", err)
	}

	return body, nil
}

func ITunesSearchByITunesID(itunesID string) ([]byte, error) {

	req, err := http.NewRequest("GET", iTunesLookupURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create iTunes ID lookup request: %w", err)
	}

	q := url.Values{}
	q.Set("id", itunesID)
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send iTunes ID lookup request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes ID lookup returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read iTunes ID lookup response: %w", err)
	}

	return body, nil
}
