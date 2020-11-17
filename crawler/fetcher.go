package crawler

// Fetcher defines the Fetch method.
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (item Item, err error)
}
