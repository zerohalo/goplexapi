package goplexapi

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PlexClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewPlexClient(baseURL, token string) *PlexClient {
	return &PlexClient{
		BaseURL: baseURL,
		Token:   token,
		Client:  &http.Client{},
	}
}

type MediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Size    int      `xml:"size,attr"`
	Tracks  []Track  `xml:"Track"`
}

type Track struct {
	XMLName              xml.Name `xml:"Track"`
	AddedAt              int      `xml:"addedAt,attr"`
	Art                  string   `xml:"art,attr"`
	Duration             int      `xml:"duration,attr"`
	GrandparentArt       string   `xml:"grandparentArt,attr"`
	GrandparentGuid      string   `xml:"grandparentGuid,attr"`
	GrandparentKey       string   `xml:"grandparentKey,attr"`
	GrandparentRatingKey string   `xml:"grandparentRatingKey,attr"`
	GrandparentThumb     string   `xml:"grandparentThumb,attr"`
	GrandparentTitle     string   `xml:"grandparentTitle,attr"`
	Guid                 string   `xml:"guid,attr"`
	Index                int      `xml:"index,attr"`
	Key                  string   `xml:"key,attr"`
	LibrarySectionID     string   `xml:"librarySectionID,attr"`
	LibrarySectionKey    string   `xml:"librarySectionKey,attr"`
	LibrarySectionTitle  string   `xml:"librarySectionTitle,attr"`
	MusicAnalysisVersion int      `xml:"musicAnalysisVersion,attr"`
	ParentGuid           string   `xml:"parentGuid,attr"`
	ParentIndex          int      `xml:"parentIndex,attr"`
	ParentKey            string   `xml:"parentKey,attr"`
	ParentRatingKey      string   `xml:"parentRatingKey,attr"`
	ParentStudio         string   `xml:"parentStudio,attr"`
	ParentThumb          string   `xml:"parentThumb,attr"`
	ParentTitle          string   `xml:"parentTitle,attr"`
	ParentYear           int      `xml:"parentYear,attr"`
	RatingCount          int      `xml:"ratingCount,attr"`
	RatingKey            string   `xml:"ratingKey,attr"`
	SessionKey           string   `xml:"sessionKey,attr"`
	Thumb                string   `xml:"thumb,attr"`
	Title                string   `xml:"title,attr"`
	Type                 string   `xml:"type,attr"`
	UpdatedAt            int      `xml:"updatedAt,attr"`
	ViewOffset           int      `xml:"viewOffset,attr"`
	Media                []Media  `xml:"Media"`
	Player               Player   `xml:"Player"`
}

type Media struct {
	XMLName       xml.Name `xml:"Media"`
	AudioChannels int      `xml:"audioChannels,attr"`
	AudioCodec    string   `xml:"audioCodec,attr"`
	Bitrate       int      `xml:"bitrate,attr"`
	Duration      int      `xml:"duration,attr"`
	Parts         []Part   `xml:"Part"`
}

type Part struct {
	XMLName xml.Name `xml:"Part"`
	File    string   `xml:"file,attr"`
}

type Stream struct {
	XMLName xml.Name `xml:"Stream"`
	Codec   string   `xml:"codec,attr"`
}

type User struct {
	XMLName xml.Name `xml:"User"`
	ID      string   `xml:"id,attr"`
	Thumb   string   `xml:"thumb,attr"`
	Title   string   `xml:"title,attr"`
}

type Player struct {
	XMLName             xml.Name `xml:"Player"`
	Title               string   `xml:"title,attr"`
	Address             string   `xml:"address,attr"`
	Device              string   `xml:"device,attr"`
	MachineIdentifier   string   `xml:"machineIdentifier,attr"`
	Platform            string   `xml:"platform,attr"`
	PlatformVersion     string   `xml:"platformVersion,attr"`
	Product             string   `xml:"product,attr"`
	RemotePublicAddress string   `xml:"remotePublicAddress,attr"`
	State               string   `xml:"state,attr"`
	Version             string   `xml:"version,attr"`
	Local               string   `xml:"local,attr"`
	Relayed             string   `xml:"relayed,attr"`
	Secure              string   `xml:"secure,attr"`
	UserID              string   `xml:"userID,attr"`
}

type TrackInfo struct {
	Artist string
	Album  string
	Title  string
	Thumb  string
}

func (p *PlexClient) makeRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", p.BaseURL, endpoint)
	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequest(method, url, strings.NewReader(payload.(string)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Plex-Token", p.Token)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *PlexClient) GetCurrentPlayingSong(clientName, userID string) (*TrackInfo, error) {
	data, err := p.makeRequest("GET", "/status/sessions", nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	var mediaContainer MediaContainer
	if err := xml.Unmarshal(data, &mediaContainer); err != nil {
		return nil, err
	}

	for _, track := range mediaContainer.Tracks {
		if track.Player.Product == clientName && track.Player.UserID == userID {
			return &TrackInfo{
				Artist: track.GrandparentTitle,
				Album:  track.ParentTitle,
				Title:  track.Title,
				Thumb:  track.ParentThumb,
			}, nil
		}
	}

	return nil, fmt.Errorf("No song currently playing on %s", clientName)
}

func (p *PlexClient) GetAlbumArt(albumArtURL string) ([]byte, error) {
	data, err := p.makeRequest("GET", albumArtURL, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
