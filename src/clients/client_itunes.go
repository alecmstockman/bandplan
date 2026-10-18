package clients

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func ClientITunesGetArtwork(artworkURL string) ([]byte, error) {

	parsedURL, err := url.Parse(artworkURL)
	if err != nil {
		return nil, fmt.Errorf("parse iTunes artwork URL: %w", err)
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) == 0 {
		return nil, errors.New("invalid iTunes artwork URL path")
	}
	fileName := parts[len(parts)-1]

	newFileName := strings.Replace(fileName, "100x100", "1000x1000", 1)

	parts[len(parts)-1] = newFileName
	parsedURL.Path = strings.Join(parts, "/")

	newURL := parsedURL.String()

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create iTunes artwork request: %w", err)
	}

	req.Header.Set("User-Agent", "BandPlan/0.1")
	req.Header.Set("Accept", "image/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request iTunes artwork: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes artwork request returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read iTunes artwork response: %w", err)
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
