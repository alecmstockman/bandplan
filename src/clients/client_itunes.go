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

	return body, nil
}

func ITunesSearchByArtistAndSong(term string) ([]byte, error) {
	log.Println("- ITunesSearchByArtistAndSong")

	req, err := http.NewRequest("GET", iTunesBaseURL, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("term", term)
	q.Set("media", "music")
	q.Set("entity", "song")
	q.Set("limit", "1")
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

	return body, nil
}

func ClientITunesGetArtwork(artworkURL string) ([]byte, error) {
	log.Println("- ITunesGetArtwork")

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

	newFileName := strings.Replace(fileName, "100x100", "1024x1024", 1)

	parts[len(parts)-1] = newFileName
	parsedURL.Path = strings.Join(parts, "/")

	newURL := parsedURL.String()

	// log.Println("   newURL:", newURL)

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		log.Println("   Error creating http.NewRequest with newURL: ", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("   Unable to complete http request for iTunes Art", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("   iTunes Artwork request response: ", resp.StatusCode)
		return nil, fmt.Errorf("   itunes artwork requests returned status: ", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("   Unable to read iTunes Artwork response body: ", err)
		return nil, err
	}

	return body, nil
}
