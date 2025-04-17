// announcement_parser.go
package mexc

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Announcement represents a parsed MEXC announcement.
type Announcement struct {
	Symbol      string
	ListingTime time.Time
	URL         string
	Title       string
	RawHTML     string // Optional: for debugging or structure change detection
	ParsedAt    time.Time
}

// AnnouncementParserConfig holds configuration for the parser.
type AnnouncementParserConfig struct {
	AnnouncementsURL string
	TitleSelector    string
	SymbolSelector   string
	TimeSelector     string
	PollInterval     time.Duration
	MaxRetries       int
	RetryDelay       time.Duration
}

// AnnouncementParser is responsible for polling and parsing MEXC announcements.
type AnnouncementParser struct {
	config     AnnouncementParserConfig
	logger     *log.Logger
	httpClient *http.Client
	stopCh     chan struct{}
}

// NewAnnouncementParser creates a new parser with the given config and logger.
func NewAnnouncementParser(cfg AnnouncementParserConfig, logger *log.Logger, httpClient *http.Client) *AnnouncementParser {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &AnnouncementParser{
		config:     cfg,
		logger:     logger,
		httpClient: httpClient,
		stopCh:     make(chan struct{}),
	}
}

// StartPolling begins polling for new announcements and sends them to the provided channel.
func (p *AnnouncementParser) StartPolling(out chan<- Announcement) {
	// TODO: Implement polling loop using config.PollInterval
	// - Fetch HTML with retries
	// - Parse announcements using goquery and selectors from config
	// - Extract symbol, listing time (UTC), URL, title, etc.
	// - Handle errors with logging and alerts
	// - Send parsed Announcement to 'out' channel
}

// StopPolling signals the polling loop to stop.
func (p *AnnouncementParser) StopPolling() {
	close(p.stopCh)
}

// fetchHTML fetches the HTML from the configured URL with retries.
func (p *AnnouncementParser) fetchHTML() (string, error) {
	var lastErr error
	for attempt := 0; attempt < p.config.MaxRetries; attempt++ {
		resp, err := p.httpClient.Get(p.config.AnnouncementsURL)
		if err != nil {
			p.logger.Printf("fetchHTML attempt %d: %v", attempt+1, err)
			lastErr = err
			time.Sleep(p.config.RetryDelay)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			p.logger.Printf("fetchHTML attempt %d: non-200 status %d", attempt+1, resp.StatusCode)
			lastErr = err
			time.Sleep(p.config.RetryDelay)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			p.logger.Printf("fetchHTML attempt %d: read error: %v", attempt+1, err)
			lastErr = err
			time.Sleep(p.config.RetryDelay)
			continue
		}
		return string(body), nil
	}
	return "", lastErr
}

// parseAnnouncements parses the HTML and extracts announcements using goquery and config selectors.
func (p *AnnouncementParser) parseAnnouncements(html string) ([]Announcement, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var announcements []Announcement
	// For demo, assume each announcement is a div. In real use, selector should be in config.
	doc.Find("div.announcement").Each(func(i int, s *goquery.Selection) {
		titleSel := p.config.TitleSelector
		symbolSel := p.config.SymbolSelector
		timeSel := p.config.TimeSelector

		title := s.Find(titleSel).Text()
		url, _ := s.Find(titleSel).Attr("href")
		symbol := s.Find(symbolSel).Text()
		timeStr := s.Find(timeSel).Text()

		// Parse time (assume format "2006-01-02 15:04:05 UTC")
		var listingTime time.Time
		if timeStr != "" {
			listingTime, _ = time.Parse("2006-01-02 15:04:05 UTC", timeStr)
		}

		announcements = append(announcements, Announcement{
			Symbol:      symbol,
			ListingTime: listingTime,
			URL:         url,
			Title:       title,
			RawHTML:     "", // Optionally store raw HTML
			ParsedAt:    time.Now().UTC(),
		})
	})

	return announcements, nil
}
