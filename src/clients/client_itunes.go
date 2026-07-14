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
	log.Println("- ITunesSearchByArtistAndSong")

	req, err := http.NewRequest("GET", iTunesBaseURL, nil)
	if err != nil {
		log.Println("   Error making http request: ", err)
		return nil, err
	}

	q := url.Values{}
	q.Set("term", term)
	q.Set("media", "music")
	q.Set("entity", "song")
	q.Set("limit", "15")
	q.Set("country", "US")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("   Error sending reqeust: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes api returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("   Error reading response body: ", body)
		return nil, err
	}

	return body, nil
}

func ClientITunesGetArtwork(artworkURL string) ([]byte, error) {
	log.Println("- ClientITunesGetArtwork")

	parsedURL, err := url.Parse(artworkURL)
	if err != nil {
		log.Println("   Unable to parse artwork url: ", err)
		return nil, nil
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid artwork URL path")
	}
	fileName := parts[len(parts)-1]

	newFileName := strings.Replace(fileName, "100x100", "1000x1000", 1)

	parts[len(parts)-1] = newFileName
	parsedURL.Path = strings.Join(parts, "/")

	newURL := parsedURL.String()

	log.Println("   Artwork URL: ", newURL)

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		log.Println("   Error creating http.NewRequest with newURL: ", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "BandPlan/0.1")
	req.Header.Set("Accept", "image/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("   Unable to complete http request for iTunes Art", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("   iTunes Artwork request response: ", resp.StatusCode)
		return nil, fmt.Errorf("   itunes artwork requests returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("   Unable to read iTunes Artwork response body: ", err)
		return nil, err
	}

	return body, nil
}

func ITunesSearchByITunesID(itunesID string) ([]byte, error) {
	log.Println("- ITunesSearchByITunesID")

	req, err := http.NewRequest("GET", iTunesLookupURL, nil)
	if err != nil {
		log.Println("   Unable to make http request: ", err)
		return nil, err
	}

	q := url.Values{}
	q.Set("id", itunesID)
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
