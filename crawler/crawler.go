package crawler

import (
	"github.com/google/uuid"
)

// Item is the result of fetching data from one url.
type Item struct {
	Body  string   `json:"body"`
	Links []string `json:"links"`
}

// Crawler defines the Crawl method.
type Crawler interface {
	// Crawl takes an URL and a depth, and
	// returns a slice of items as a result
	// of fetching all the URLs until a depth.
	ID() string
	Get() <-chan Item
	Depth() int
	Close() error
}

type crawl struct {
	id      uuid.UUID
	url     string
	depth   int
	items   chan Item
	closing chan chan error
}

func (c *crawl) ID() string {
	return c.id.String()
}

func (c *crawl) Get() <-chan Item {
	return c.items
}

func (c *crawl) Depth() int {
	return c.depth
}

func (c *crawl) Close() error {
	errc := make(chan error)
	c.closing <- errc
	return <-errc
}

func (c *crawl) loop() {

	type fetchResult struct {
		body  string
		links []string
		err   error
	}

	const maxQueued = 5

	// List of pending links to visit. Start with the parameter url
	pending := []string{c.url}
	visited := map[string]bool{c.url: true}
	sendq := []Item{}

	var err error
	var fetchDone chan fetchResult
	next := make(chan string, 1)
	for {
		// Enable fetching if nothing is being fetch and queue is not full
		var willFetch string
		var willFetchChannel chan string
		if fetchDone == nil && len(pending) > 0 && len(sendq) < maxQueued {
			willFetch = pending[0]
			willFetchChannel = make(chan string, 1)
		}

		// Send the result items one by one.
		var first Item
		var items chan Item
		if len(sendq) > 0 {
			first = sendq[0]
			items = c.items
		}

		// Finish condition for the crawler:
		// - There's no pending url to visit
		// - There's no pending result to send to the api
		// - There's no URL that will be fetch in the next steps
		// - There's no Item that will be fetch in the next steps
		// - There's no fetching that's being done
		if len(pending) == 0 && len(sendq) == 0 && willFetchChannel == nil && items == nil && fetchDone == nil {
			go func() { c.Close() }()
		}

		select {
		case willFetchChannel <- willFetch: // Pick an URL and send it on the next channel.
			pending = pending[1:]
			fetchDone = make(chan fetchResult, 1)
			go func() { next <- willFetch }()

		case url := <-next: // Fetch the url from the next channel
			go func() {
				body, links, err := FakeFetcher.Fetch(url)
				fetchDone <- fetchResult{body, links, err}
			}()

		case res := <-fetchDone: // Process the result from fetching.
			fetchDone = nil
			item := Item{res.body, res.links}
			err := res.err

			if err != nil {
				break //TODO: Set a retry policy on recoverable errors.
			}

			sendq = append(sendq, item)
			for _, url := range item.Links {
				if !visited[url] {
					pending = append(pending, url)
					visited[url] = true
				}
			}

		case errc := <-c.closing: // Process cancelation request.
			errc <- err
			close(c.items)
			return

		case items <- first: // Remove sent items from the sending queue.
			sendq = sendq[1:]
		}
	}
}
