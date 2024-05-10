package trackmania

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LeaderboardData struct {
	Tops []struct {
		Player struct {
			Name string `json:"name"`
			Tag  string `json:"tag"`
			ID   string `json:"id"`
			Zone struct {
				Name   string `json:"name"`
				Flag   string `json:"flag"`
				Parent struct {
					Name   string `json:"name"`
					Flag   string `json:"flag"`
					Parent struct {
						Name   string `json:"name"`
						Flag   string `json:"flag"`
						Parent struct {
							Name string `json:"name"`
							Flag string `json:"flag"`
						} `json:"parent"`
					} `json:"parent"`
				} `json:"parent"`
			} `json:"zone"`
			Meta struct {
				Twitch  string `json:"twitch"`
				Youtube string `json:"youtube"`
				Twitter string `json:"twitter"`
			} `json:"meta"`
		} `json:"player"`
		Position int `json:"position"`
		Time     int `json:"time"`
		Points   int `json:"points"`
	} `json:"tops"`
}

var (
	leaderboardCache         = make(map[string]*LeaderboardData)
	leaderboardCacheDuration = 5 * time.Minute // Adjust the cache duration as needed
)

func FetchLeaderboardData(campaignID string, offset int) (*LeaderboardData, error) {
	// Create a cache key that includes both campaignID and offset
	cacheKey := fmt.Sprintf("%s_%d", campaignID, offset)

	// Check if the data is already in the cache
	if data, ok := leaderboardCache[cacheKey]; ok {
		return data, nil
	}

	url := fmt.Sprintf("https://trackmania.io/api/leaderboard/%s?offset=%d&length=100", campaignID, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("For questions about this project, contact me on Discord: %s", name))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data LeaderboardData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Store the fetched data in the cache
	leaderboardCache[campaignID] = &data

	// Start a goroutine to remove the data from the cache after the specified duration
	go func() {
		time.Sleep(leaderboardCacheDuration)
		delete(leaderboardCache, campaignID)
	}()

	return &data, nil
}
