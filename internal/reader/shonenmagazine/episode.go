package shonenmagazine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oxdz/comicinfo/pkg/decode"
)

const (
	APIURL = "https://api.pocket.shonenmagazine.com/web/episode/viewer?platform=3&episode_id=416173"
)

var (
	client = &http.Client{
		Transport: &http.Transport{},
	}
)

func EpisodeInfo(ctx context.Context, episodeID int, cookies []*http.Cookie) (*decode.EpisodeInfo, error) {
	url := APIURL + fmt.Sprintf("?platform=3&episode_id=%d", episodeID)
	mangaHash := decode.MangeHash(episodeID, 3)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-manga-hash", mangaHash)
	req.Header.Set("x-manga-is-crawler", "false")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response: status: %s, body: %s", resp.Status, string(buf))
	}

	var info decode.EpisodeInfo
	if err := json.Unmarshal(buf, &info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &info, nil
}
