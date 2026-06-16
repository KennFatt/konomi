package forgejo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIError wraps a non-2xx Forgejo response.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

var _ error = (*APIError)(nil)

// Client is a lightweight Forgejo API client.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// New creates a new API client.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// doGet performs a single GET request and returns the raw body.
func (c *Client) doGet(path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + "/api/v1" + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := extractMessage(body)
		return nil, &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	return body, nil
}

// getArray fetches a paginated list endpoint and returns a single JSON array.
func (c *Client) getArray(path string, params map[string]string) ([]byte, error) {
	const pageSize = 50
	page := 1
	var all []json.RawMessage

	for {
		p := make(map[string]string, len(params)+2)
		for k, v := range params {
			p[k] = v
		}
		p["page"] = strconv.Itoa(page)
		p["limit"] = strconv.Itoa(pageSize)

		body, err := c.doGet(path, p)
		if err != nil {
			return nil, err
		}

		var items []json.RawMessage
		if err := json.Unmarshal(body, &items); err != nil {
			return nil, fmt.Errorf("decode page %d: %w", page, err)
		}

		if len(items) == 0 {
			break
		}

		all = append(all, items...)

		// Guards against endpoints that return the same data on every page.
		if len(items) < pageSize {
			break
		}

		page++
	}

	if len(all) == 0 {
		return []byte("[]"), nil
	}

	result, err := json.Marshal(all)
	if err != nil {
		return nil, fmt.Errorf("marshal combined: %w", err)
	}
	return result, nil
}

func extractMessage(body []byte) string {
	var tmp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &tmp); err == nil && tmp.Message != "" {
		return tmp.Message
	}
	if len(body) > 200 {
		return string(body[:200])
	}
	return string(body)
}
