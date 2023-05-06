package schema

// RequestJSON describes Request with URL in it.
type RequestJSON struct {
	URL string `json:"url"`
}

// ResponseJSON describes Response with URL in it.
type ResponseJSON struct {
	Result string `json:"result"`
}

// URL contains LongURL and ShortURL fields.
type URL struct {
	LongURL  string `json:"original_url"`
	ShortURL string `json:"short_url"`
}

// BatchURL contains Chars - chars that come after the slash in the URL,
// Original field - original URL.
type BatchURL struct {
	LongURL  string `json:"correlation_id"`
	ShortURL string `json:"original_url"`
}

// ResponseBatchURL describes Response that gives BatchURL Handler.
type ResponseBatchURL struct {
	Chars   string `json:"correlation_id"`
	Shorted string `json:"short_url"`
}

// StatsResponse describes Response that gives GetStats Handler.
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
