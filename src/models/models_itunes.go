package models

type ITunesSearchResponse struct {
	ResultCount int          `json:"resultCount"`
	Results     []ITunesSong `json:"results"`
}

type ITunesSong struct {
	TrackID          int    `json:"trackId"`
	ArtistName       string `json:"artistName"`
	TrackName        string `json:"trackName"`
	CollectionName   string `json:"collectionName"`
	ArtworkURL100    string `json:"artworkUrl100"`
	TrackViewURL     string `json:"trackViewUrl"`
	PreviewURL       string `json:"previewUrl"`
	TrackTimeMillis  int    `json:"trackTimeMillis"`
	PrimaryGenreName string `json:"primaryGenreName"`
	ReleaseDate      string `json:"releaseDate"`
}
