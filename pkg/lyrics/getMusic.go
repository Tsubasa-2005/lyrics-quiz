package lyrics

import (
	"fmt"
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/message"
	"math/rand"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

func GetArtistTopTracksBySpotifyAPI(ctx *gin.Context, repo *rdb.Queries, artistName string, quizManager rdb.QuizManager, bot *linebot.Client, event *linebot.Event) ([]string, error) {
	numTracks := int(quizManager.TheNumberOfQuestions) * 4
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

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

	albums, err := client.GetArtistAlbums(ctx, artistID, nil)
	if err != nil {
		return nil, fmt.Errorf("no artist album found")
	}

	albumIDs := make([]spotify.ID, 0)
	for _, album := range albums.Albums {
		albumIDs = append(albumIDs, album.ID)
	}
	fullTracks := make([]string, 0)
	for _, albumID := range albumIDs {
		albums, err := client.GetAlbumTracks(ctx, albumID)
		if err != nil {
			return nil, fmt.Errorf("no artist album found")
		}
		for _, track := range albums.Tracks {
			fullTracks = append(fullTracks, track.Name)
		}
	}

	// "- Off Vocal"と書かれている曲を除外
	filteredTracks := make([]string, 0)
	for _, track := range fullTracks {
		if !strings.Contains(track, " - Off Vocal") && !strings.Contains(track, " - 104期 Ver.") {
			filteredTracks = append(filteredTracks, track)
		}
	}

	if len(filteredTracks) == 0 {
		return nil, fmt.Errorf("no valid tracks found for artist")
	}

	if numTracks > len(filteredTracks) {
		numTracks = len(filteredTracks)
	}
	// Shuffle indices
	for i := 0; i < 3; i++ {
		rand.Shuffle(len(filteredTracks), func(i, j int) { filteredTracks[i], filteredTracks[j] = filteredTracks[j], filteredTracks[i] })
	}

	// Select the first numTracks indices
	selectedTracks := filteredTracks[:numTracks]
	trackNames := make([]string, quizManager.TheNumberOfQuestions)
	selectedMusic := 0
	for len(selectedTracks) < int(quizManager.TheNumberOfQuestions)*4 {
		quizManager.TheNumberOfQuestions--
	}
	for i := 0; i < int(quizManager.TheNumberOfQuestions); i++ {
		choice1 := selectedTracks[i*4]
		choice2 := selectedTracks[i*4+1]
		choice3 := selectedTracks[i*4+2]
		choice4 := selectedTracks[i*4+3]

		err = repo.CreateChoices(ctx, rdb.CreateChoicesParams{
			QuizManagerID:  quizManager.UserID,
			QuestionNumber: int64(i + 1),
			Choice1:        choice1,
			Choice2:        choice2,
			Choice3:        choice3,
			Choice4:        choice4,
		})
		if err != nil {
			message.ErrorCreatingProblems(ctx, bot, event)
			break // エラーが発生したらループを終了
		}

		// 各グループからランダムに1曲ずつ選んでtrackNamesに追加
		trackNames[selectedMusic] = selectedTracks[i*4]
		err = repo.CreateAnswer(ctx, rdb.CreateAnswerParams{
			QuizManagerID:  quizManager.UserID,
			QuestionNumber: int64(i + 1),
			MusicName:      selectedTracks[i*4],
		})
		if err != nil {
			message.ErrorCreatingProblems(ctx, bot, event)
			break // エラーが発生したらループを終了
		}
		selectedMusic++
	}

	return trackNames, nil
}
