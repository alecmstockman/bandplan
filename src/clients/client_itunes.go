package clients

import (
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
		return nil, fmt.Errorf("ITunesSearchByArtistAndSong - error making request: %w", err)
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
		return nil, fmt.Errorf("ITunesSearchByArtistAndSong - error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ITunesSearchByArtistAndSong - itunes api returned status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ITunesSearchByArtistAndSong - rror reading response body: %v", err)
	}

	return body, nil
}

func ClientITunesGetArtwork(artworkURL string) ([]byte, error) {

	parsedURL, err := url.Parse(artworkURL)
	if err != nil {
		return nil, fmt.Errorf("ClientITunesGetArtwork - unable to parse artwork url: %v", err)
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("ClientITunesGetArtwork - invalid artwork URL path")
	}
	fileName := parts[len(parts)-1]

	newFileName := strings.Replace(fileName, "100x100", "1000x1000", 1)

	parts[len(parts)-1] = newFileName
	parsedURL.Path = strings.Join(parts, "/")

	newURL := parsedURL.String()

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ClientITunesGetArtwork - error creating http.NewRequest with newURL: %v", err)
	}

	req.Header.Set("User-Agent", "BandPlan/0.1")
	req.Header.Set("Accept", "image/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ClientITunesGetArtwork - unable to complete http request for iTunes Art: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ClientITunesGetArtwork - itunes artwork requests returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ClientITunesGetArtwork - unable to read itunes artwork response body: %v", err)
	}

	return body, nil
}

func ITunesSearchByITunesID(itunesID string) ([]byte, error) {

	req, err := http.NewRequest("GET", iTunesLookupURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ITunesSearchByITunesID - unable to make http request: %v", err)
	}

	q := url.Values{}
	q.Set("id", itunesID)
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ITunesSearchByITunesID - err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ITunesSearchByITunesID - itunes api returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ITunesSearchByITunesID - unable to read resp body: %s", resp.Status)
	}

	return body, nil
}
