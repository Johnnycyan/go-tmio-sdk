package trackmania

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	client = &http.Client{
		Timeout: time.Second * 30,
	}
	name string
)

type PlayerInfo struct {
	Player      Player      `json:"player"`
	Matchmaking []MatchType `json:"matchmaking"`
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Tag  string `json:"tag,omitempty"`
	Zone Zone   `json:"zone"`
	Meta Meta   `json:"meta,omitempty"`
}

type Zone struct {
	Name   string `json:"name"`
	Flag   string `json:"flag"`
	Parent *Zone  `json:"parent,omitempty"`
}

type Meta struct {
	Vanity  string `json:"vanity,omitempty"`
	Twitch  string `json:"twitch,omitempty"`
	Twitter string `json:"twitter,omitempty"`
}

type MatchType struct {
	TypeName     string   `json:"typename"`
	TypeID       int      `json:"typeid"`
	AccountID    string   `json:"accountid"`
	Rank         int      `json:"rank"`
	Score        int      `json:"score"`
	Progression  int      `json:"progression"`
	Division     Division `json:"division"`
	DivisionNext Division `json:"division_next,omitempty"`
}

type Division struct {
	Position  int    `json:"position"`
	Rule      string `json:"rule"`
	MinPoints int    `json:"minpoints"`
	MaxPoints int    `json:"maxpoints"`
	MinWins   int    `json:"minwins,omitempty"`
	MaxWins   int    `json:"maxwins,omitempty"`
}

type cacheEntry struct {
	data      PlayerInfo
	timestamp time.Time
}

var cache = make(map[string]*cacheEntry)
var mutex = &sync.Mutex{}

const cacheDuration = time.Hour * 48

func getCachedPlayerInfo(key string) (*PlayerInfo, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getCachedPlayerInfo", r)
		}
	}()
	mutex.Lock()
	defer mutex.Unlock()

	entry, exists := cache[key]
	if !exists {
		return nil, false
	}
	if time.Since(entry.timestamp) > cacheDuration {
		delete(cache, key)
		return nil, false
	}
	return &entry.data, true
}

func fetchPlayerInfo(playerName string) ([]PlayerInfo, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in fetchPlayerInfo", r)
		}
	}()
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	name = os.Getenv("NAME")
	url := fmt.Sprintf("https://trackmania.io/api/players/find?search=%s", playerName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("Displays number of players in COTD using a Twitch command. For questions about this project, contact me on Discord: %s", name))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []PlayerInfo
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func GetPlayerID(playerName string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetPlayerID", r)
		}
	}()
	if data, found := getCachedPlayerInfo(playerName); found {
		return data.Player.ID, nil
	}

	results, err := fetchPlayerInfo(playerName)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("player not found")
	}

	mutex.Lock()
	cache[playerName] = &cacheEntry{data: results[0], timestamp: time.Now()}
	mutex.Unlock()

	return results[0].Player.ID, nil
}

func GetFormattedName(playerName string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetFormattedName", r)
		}
	}()
	if data, found := getCachedPlayerInfo(playerName); found {
		return data.Player.Name, nil
	}

	results, err := fetchPlayerInfo(playerName)
	if err != nil {
		return "", err
	}

	mutex.Lock()
	cache[playerName] = &cacheEntry{data: results[0], timestamp: time.Now()}
	mutex.Unlock()

	return results[0].Player.Name, nil
}

func GetPlayerCampaignRank(playerName string, campaign string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetPlayerCampaignRank", r)
		}
	}()
	playerID, err := GetPlayerID(playerName)
	if err != nil {
		return 0, err
	}

	campaignResults, err := SearchCampaigns(campaign)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	leaderboardID, err := FetchCampaignLeaderboardID(campaignResults[0].ID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	for i := 0; i < 5; i++ {
		offset := i * 100
		leaderboards, err := FetchLeaderboardData(leaderboardID, offset)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		for _, top := range leaderboards.Tops {
			if top.Player.ID == playerID {
				return top.Position, nil
			}
		}
	}

	return 0, fmt.Errorf("player not found in leaderboard top 500")
}

func GetPlayerCampaignPoints(playerName string, campaign string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetPlayerCampaignPoints", r)
		}
	}()
	playerID, err := GetPlayerID(playerName)
	if err != nil {
		return 0, err
	}

	campaignResults, err := SearchCampaigns(campaign)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	leaderboardID, err := FetchCampaignLeaderboardID(campaignResults[0].ID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	for i := 0; i < 5; i++ {
		offset := i * 100
		leaderboards, err := FetchLeaderboardData(leaderboardID, offset)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		for _, top := range leaderboards.Tops {
			if top.Player.ID == playerID {
				return top.Points, nil
			}
		}
	}

	return 0, fmt.Errorf("player not found in leaderboard top 500")
}
