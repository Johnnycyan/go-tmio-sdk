package trackmania

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Campaign struct {
	Official  bool   `json:"official"`
	ID        int    `json:"id"`
	ClubID    int    `json:"clubid"`
	ClubName  string `json:"clubname"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	MapCount  int    `json:"mapcount"`
	Tracked   bool   `json:"tracked"`
	// Ignoring other fields for simplicity
}

type Response struct {
	Page      int        `json:"page"`
	PageCount int        `json:"pageCount"`
	Campaigns []Campaign `json:"campaigns"`
}

type OfficialCampaign struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Media         string     `json:"media"`
	CreationTime  int        `json:"creationtime"`
	PublishTime   int        `json:"publishtime"`
	ClubID        int        `json:"clubid"`
	LeaderboardID string     `json:"leaderboarduid"`
	Playlist      []MapInfo  `json:"playlist"`
	Mediae        MediaeInfo `json:"mediae"`
	Tracked       bool       `json:"tracked"`
}

type MapInfo struct {
	Author          string         `json:"author"`
	Name            string         `json:"name"`
	MapType         string         `json:"mapType"`
	MapStyle        string         `json:"mapStyle"`
	AuthorScore     int            `json:"authorScore"`
	GoldScore       int            `json:"goldScore"`
	SilverScore     int            `json:"silverScore"`
	BronzeScore     int            `json:"bronzeScore"`
	CollectionName  string         `json:"collectionName"`
	Filename        string         `json:"filename"`
	IsPlayable      bool           `json:"isPlayable"`
	MapID           string         `json:"mapId"`
	MapUID          string         `json:"mapUid"`
	Submitter       string         `json:"submitter"`
	Timestamp       string         `json:"timestamp"`
	FileURL         string         `json:"fileUrl"`
	ThumbnailURL    string         `json:"thumbnailUrl"`
	AuthorPlayer    PlayerMetaInfo `json:"authorplayer"`
	SubmitterPlayer PlayerMetaInfo `json:"submitterplayer"`
	ExchangeID      int            `json:"exchangeid"`
}

type PlayerMetaInfo struct {
	Name string                 `json:"name"`
	ID   string                 `json:"id"`
	Meta map[string]interface{} `json:"meta"`
}

type MediaeInfo struct {
	ButtonBackground     string `json:"buttonbackground"`
	ButtonForeground     string `json:"buttonforeground"`
	Decal                string `json:"decal"`
	PopupBackground      string `json:"popupbackground"`
	Popup                string `json:"popup"`
	LiveButtonBackground string `json:"livebuttonbackground"`
	LiveButtonForeground string `json:"livebuttonforeground"`
}

var (
	campaignLeaderboardCache = make(map[int]string)
	campaignsCache           []Campaign
	campaignCacheDuration    = 5 * time.Minute // Adjust the cache duration as needed
)

func FetchCampaignLeaderboardID(campaignID int) (string, error) {
	// Check if the leaderboard ID is already in the cache
	if leaderboardID, ok := campaignLeaderboardCache[campaignID]; ok {
		return leaderboardID, nil
	}

	url := "https://trackmania.io/api/officialcampaign/" + strconv.Itoa(campaignID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("For questions about this project, contact me on Discord: %s", name))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var officialCampaign OfficialCampaign
	err = json.Unmarshal(body, &officialCampaign)
	if err != nil {
		return "", err
	}

	// Store the leaderboard ID in the cache
	campaignLeaderboardCache[campaignID] = officialCampaign.LeaderboardID

	// Start a goroutine to remove the leaderboard ID from the cache after the specified duration
	go func() {
		time.Sleep(campaignCacheDuration)
		delete(campaignLeaderboardCache, campaignID)
	}()

	return officialCampaign.LeaderboardID, nil
}

func FetchCampaigns() ([]Campaign, error) {
	// Check if the campaigns are already in the cache
	if len(campaignsCache) > 0 {
		return campaignsCache, nil
	}

	url := "https://trackmania.io/api/campaigns/0"
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

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var officialCampaigns []Campaign
	for _, campaign := range response.Campaigns {
		if campaign.Official {
			officialCampaigns = append(officialCampaigns, campaign)
		}
	}

	// Store the fetched campaigns in the cache
	campaignsCache = officialCampaigns

	// Start a goroutine to clear the campaigns cache after the specified duration
	go func() {
		time.Sleep(campaignCacheDuration)
		campaignsCache = nil
	}()

	return officialCampaigns, nil
}

func SearchCampaigns(search string) ([]Campaign, error) {
	campaigns, err := FetchCampaigns()
	if err != nil {
		return nil, err
	}

	var results []Campaign
	for _, campaign := range campaigns {
		if strings.EqualFold(campaign.Name, search) {
			results = append(results, campaign)
		}
	}

	return results, nil
}
