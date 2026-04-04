package runner

import (
	"bufio"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// InputProcessor reads URLs from a single URL string or a file.
type InputProcessor struct {
	urls chan string
}

// Regex patterns for URL extraction.
// Compiled once at package initialization for performance.
var (
	// Matches full URLs starting with http:// or https://
	// Excludes: whitespace, |, <, >, ", ', [, ], (, ), comma, semicolon, backtick
	fullURLPattern = regexp.MustCompile("https?://[^\\s|<>\"'\\[\\](),;`]+")

	// Matches domain names without protocol (e.g., example.com, sub.example.co.uk)
	// Also matches localhost with optional port and path
	// Excludes: whitespace, |, <, >, ", ', [, ], (, ), comma, semicolon, backtick
	domainPattern = regexp.MustCompile("(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)*\\.[a-zA-Z]{2,}|localhost)(?::\\d+)?(?:/[^\\s|<>\"'\\[\\](),;`]*)?")

	// Matches IPv4 addresses with optional port and path
	// Excludes: whitespace, |, <, >, ", ', [, ], (, ), comma, semicolon, backtick
	ipPattern = regexp.MustCompile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}(?::\\d+)?(?:/[^\\s|<>\"'\\[\\](),;`]*)?")
)

// extractURLs finds all valid URLs and domains in a line of text.
// Returns a deduplicated list of normalized URLs (with http:// prefix if needed).
func extractURLs(line string) []string {
	var results []string
	seen := make(map[string]bool)

	// Find full URLs (http/https) and track their positions
	fullMatches := fullURLPattern.FindAllString(line, -1)
	fullPositions := fullURLPattern.FindAllStringIndex(line, -1)

	// Find IP addresses and track their positions
	ipMatches := ipPattern.FindAllString(line, -1)
	ipPositions := ipPattern.FindAllStringIndex(line, -1)

	// Combine excluded positions from full URLs and IPs
	excluded := append(fullPositions, ipPositions...)

	// Find domain names, excluding positions already consumed by full URLs or IPs
	domainMatches := findNonOverlapping(line, domainPattern, excluded)

	// Process full URLs
	for _, match := range fullMatches {
		cleaned := cleanMatch(match)
		if normalized := validateAndNormalize(cleaned); normalized != "" && !seen[normalized] {
			results = append(results, normalized)
			seen[normalized] = true
		}
	}

	// Process IP addresses
	for _, match := range ipMatches {
		cleaned := cleanMatch(match)
		if normalized := validateAndNormalize(cleaned); normalized != "" && !seen[normalized] {
			results = append(results, normalized)
			seen[normalized] = true
		}
	}

	// Process domain names
	for _, match := range domainMatches {
		cleaned := cleanMatch(match)
		if normalized := validateAndNormalize(cleaned); normalized != "" && !seen[normalized] {
			results = append(results, normalized)
			seen[normalized] = true
		}
	}

	return results
}

// findNonOverlapping finds regex matches that don't overlap with excluded positions.
func findNonOverlapping(line string, pattern *regexp.Regexp, excluded [][]int) []string {
	var results []string
	matches := pattern.FindAllStringIndex(line, -1)

	for _, pos := range matches {
		if !overlaps(pos, excluded) {
			results = append(results, line[pos[0]:pos[1]])
		}
	}

	return results
}

// overlaps checks if a position range overlaps with any excluded ranges.
func overlaps(pos []int, excluded [][]int) bool {
	for _, ex := range excluded {
		if pos[0] < ex[1] && pos[1] > ex[0] {
			return true
		}
	}
	return false
}

// cleanMatch removes surrounding characters that are not part of the URL.
// Handles: parentheses, brackets, angle brackets, quotes, backticks, trailing punctuation.
func cleanMatch(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Remove surrounding parentheses, brackets, angle brackets, quotes, backticks
	for {
		trimmed := s
		if len(s) >= 2 {
			first, last := s[0], s[len(s)-1]
			switch {
			case first == '(' && last == ')':
				trimmed = s[1 : len(s)-1]
			case first == '[' && last == ']':
				trimmed = s[1 : len(s)-1]
			case first == '<' && last == '>':
				trimmed = s[1 : len(s)-1]
			case first == '"' && last == '"':
				trimmed = s[1 : len(s)-1]
			case first == '\'' && last == '\'':
				trimmed = s[1 : len(s)-1]
			case first == '`' && last == '`':
				trimmed = s[1 : len(s)-1]
			}
		}
		if trimmed == s {
			break
		}
		s = trimmed
	}

	// Remove trailing punctuation that is unlikely to be part of a URL
	s = strings.TrimRight(s, ".,;:!?")

	return strings.TrimSpace(s)
}

// validateAndNormalize checks if a string is a valid URL and returns it normalized.
// Adds http:// prefix if missing. Returns empty string for invalid URLs.
func validateAndNormalize(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	// Determine if it already has a protocol
	hasProtocol := strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")

	testURL := raw
	if !hasProtocol {
		testURL = "http://" + raw
	}

	// Parse the URL
	u, err := url.Parse(testURL)
	if err != nil {
		return ""
	}

	// Must have a host
	if u.Host == "" {
		return ""
	}

	host := u.Hostname()
	if !isValidHost(host) {
		return ""
	}

	return testURL
}

// isValidHost checks if a hostname is valid.
// Valid hosts: domains with dots, valid IP addresses, or localhost.
func isValidHost(host string) bool {
	if host == "" {
		return false
	}

	// localhost is valid
	if host == "localhost" {
		return true
	}

	// Check if it looks like an IP address (only digits and dots)
	if looksLikeIP(host) {
		// Validate as actual IP
		return net.ParseIP(host) != nil
	}

	// Must contain at least one dot for a domain
	if strings.Contains(host, ".") {
		return true
	}

	return false
}

// looksLikeIP checks if a string contains only digits and dots.
func looksLikeIP(s string) bool {
	for _, r := range s {
		if r != '.' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// NewInputProcessor creates a new InputProcessor that reads URLs from a single URL string or a file.
// When reading from a file, it extracts URLs from each line using regex patterns,
// ignoring lines that contain no valid URLs.
func NewInputProcessor(singleURL, urlFile string) (*InputProcessor, error) {
	ip := &InputProcessor{
		urls: make(chan string),
	}

	go func() {
		defer close(ip.urls)

		// Process single URL if provided
		if singleURL != "" {
			ip.urls <- singleURL
		}

		// Process file if provided
		if urlFile == "" {
			return
		}

		file, err := os.Open(urlFile)
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// Increase buffer size for very long lines
		scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			// Extract all URLs from this line
			urls := extractURLs(line)

			// Send each URL to the channel
			for _, u := range urls {
				ip.urls <- u
			}
		}
	}()

	return ip, nil
}

func (ip *InputProcessor) URLs() <-chan string {
	return ip.urls
}
