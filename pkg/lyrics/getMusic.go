package lyrics

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

func GetArtistTopTracksBySpotifyAPI(artistName string, numTracks int) ([]string, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	// アーティストのIDを取得
	searchResult, err := client.Search(ctx, artistName, spotify.SearchTypeArtist, spotify.Limit(1))
	if err != nil {
		return nil, err
	}
	if len(searchResult.Artists.Artists) == 0 {
		return nil, fmt.Errorf("no artist found")
	}
	artistID := searchResult.Artists.Artists[0].ID

	// アーティストのトップトラックを取得
	topTracks, err := client.GetArtistsTopTracks(ctx, artistID, "JP")
	if err != nil {
		return nil, err
	}

	// "- Off Vocal"と書かれている曲を除外
	filteredTracks := make([]spotify.FullTrack, 0)
	for _, track := range topTracks {
		if !strings.Contains(track.Name, " - Off Vocal") {
			filteredTracks = append(filteredTracks, track)
		}
	}

	if len(filteredTracks) == 0 {
		return nil, fmt.Errorf("no valid tracks found for artist")
	}

	// 曲をランダムにシャッフル
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(filteredTracks), func(i, j int) { filteredTracks[i], filteredTracks[j] = filteredTracks[j], filteredTracks[i] })

	// 曲を指定された数だけ取得
	if numTracks > len(filteredTracks) {
		numTracks = len(filteredTracks)
	}
	selectedTracks := filteredTracks[:numTracks]

	// 曲名を表示
	trackNames := make([]string, numTracks)
	for i, track := range selectedTracks {
		trackNames[i] = track.Name
	}

	return trackNames, nil
}
