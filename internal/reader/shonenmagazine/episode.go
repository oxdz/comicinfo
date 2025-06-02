package shonenmagazine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func EpisodeInfo(ctx context.Context, episodeID int, cookies []*http.Cookie) (*Episode, error) {
	url := APIURL + fmt.Sprintf("?platform=3&episode_id=%d", episodeID)
	mangaHash := MangeHash(episodeID, 3)

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

	var info Episode
	if err := json.Unmarshal(buf, &info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &info, nil
}

type Episode struct {
	Status            string         `json:"status"`
	TitleID           int            `json:"title_id"`
	EpisodeID         int            `json:"episode_id"`
	PageStartPosition int            `json:"page_start_position"`
	ScrambleSeed      int            `json:"scramble_seed"`
	PageList          []string       `json:"page_list"`
	PreviousEpisode   map[string]int `json:"previous_episode"`
	NextEpisode       map[string]int `json:"next_episode"`
}

func MangeHash(episode_id int, platform int) string {
	e := strconv.FormatInt(int64(episode_id), 10)
	p := strconv.FormatInt(int64(platform), 10)
	v1 := decode.Sum256("platform")
	v2 := decode.Sum512(p)
	v3 := decode.Sum256("episode_id")
	v4 := decode.Sum512(e)

	v6 := decode.Sum256("")
	v7 := decode.Sum512("")

	v5 := decode.Sum256(
		fmt.Sprintf("%s_%s,%s_%s", v3, v4, v1, v2),
	)
	return decode.Sum512(v5 + v6 + "_" + v7)
}
