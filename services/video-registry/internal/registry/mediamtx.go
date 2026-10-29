package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PathItem is one entry of MediaMTX GET /v3/paths/list (v1.21.3).
type PathItem struct {
	Name      string     `json:"name"`
	ConfName  string     `json:"confName"`
	Ready     bool       `json:"ready"`
	ReadyTime *time.Time `json:"readyTime"`
	Source    *struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"source"`
	Tracks  []string `json:"tracks"`
	Tracks2 []struct {
		Codec      string         `json:"codec"`
		CodecProps map[string]any `json:"codecProps"`
	} `json:"tracks2"`
	BytesReceived uint64 `json:"bytesReceived"`
}

// PathLister lists MediaMTX paths. The service depends on this, so tests can fake MediaMTX.
type PathLister interface {
	ListPaths(ctx context.Context) ([]PathItem, error)
}

// MediaMTXClient reads the MediaMTX control API with the video-registry user.
type MediaMTXClient struct {
	BaseURL  string // e.g. http://onemind-media.onetak.svc.cluster.local:9997
	User     string
	Password string
	HTTP     *http.Client
}

type pathList struct {
	PageCount int        `json:"pageCount"`
	Items     []PathItem `json:"items"`
}

// ListPaths returns every path, following pagination.
func (c *MediaMTXClient) ListPaths(ctx context.Context) ([]PathItem, error) {
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 5 * time.Second}
	}
	var out []PathItem
	for page := 0; page < 1000; page++ {
		q := url.Values{"page": {strconv.Itoa(page)}, "itemsPerPage": {"200"}}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/v3/paths/list?"+q.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(c.User, c.Password)
		resp, err := hc.Do(req)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("mediamtx paths/list: HTTP %d", resp.StatusCode)
		}
		var pl pathList
		if err := json.Unmarshal(body, &pl); err != nil {
			return nil, fmt.Errorf("mediamtx paths/list: %w", err)
		}
		out = append(out, pl.Items...)
		if page+1 >= pl.PageCount {
			break
		}
	}
	return out, nil
}
