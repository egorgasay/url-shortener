package schema

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

type URL struct {
	LongURL  string `json:"original_url"`
	ShortURL string `json:"short_url"`
}
