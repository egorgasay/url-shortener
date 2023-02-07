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

type BatchURL struct {
	Chars    string `json:"correlation_id"`
	Original string `json:"original_url"`
}

type ResponseBatchURL struct {
	Chars   string `json:"correlation_id"`
	Shorted string `json:"short_url"`
}
