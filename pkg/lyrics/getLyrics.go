package lyrics

import (
	"fmt"
	"log"
	"lyrics-quiz/pkg/infra/rdb"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
)

func GetLyrics(c *gin.Context, repo *rdb.Queries, tracks []string, quizManager rdb.QuizManager) ([]string, error) {
	for i, track := range tracks {
		webPage := "https://www.uta-net.com/search/?Keyword=" + url.QueryEscape(track)
		resp, err := http.Get(webPage)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
		}

		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			return nil, err
		}

		// 指定のXPathの要素を取得し、そのリンクをクリック
		element := htmlquery.FindOne(doc, "/html/body/div[2]/div[2]/div[1]/div[2]/div[2]/table/tbody/tr/td[1]/a")
		if element != nil {
			link := htmlquery.SelectAttr(element, "href")
			if link != "" {
				// リンク先のページを取得
				if !strings.HasPrefix(link, "http") {
					link = "https://www.uta-net.com" + link
				}
				resp, err := http.Get(link)
				if err != nil {
					log.Printf("failed to get linked page: %s", err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Printf("failed to fetch linked page data: %d %s", resp.StatusCode, resp.Status)
					continue
				}

				doc, err := htmlquery.Parse(resp.Body)
				if err != nil {
					log.Printf("failed to load linked page html: %s", err)
					continue
				}

				// リンク先ページの歌詞要素を取得
				lyricsNode := htmlquery.FindOne(doc, "/html/body/div[2]/div[1]/div[1]/div[4]/div/div[1]")
				if lyricsNode != nil {
					lyricsHTML := htmlquery.OutputHTML(lyricsNode, true)
					lyricsParts := strings.Split(lyricsHTML, "<br/><br/>")
					for j, part := range lyricsParts {
						// HTMLタグを除去してテキスト部分だけを取得
						partNode, err := htmlquery.Parse(strings.NewReader(part))
						if err != nil {
							log.Printf("failed to parse lyrics part: %s", err)
							continue
						}
						textPart := htmlquery.InnerText(partNode)
						if quizManager.Type == "hard" {
							textPart = replaceNonHiraganaWithSquare(textPart)
						}
						err = repo.CreateLyrics(c, rdb.CreateLyricsParams{
							QuizManagerID:  quizManager.UserID,
							QuestionNumber: int64(i + 1),
							Count:          int64(j + 1),
							Lyrics:         textPart,
						})
						if err != nil {
							return nil, err
						}
					}
				} else {
					fmt.Println("Linked Page Lyrics Element not found")
				}
			}
		} else {
			return nil, fmt.Errorf("element not found")
		}
	}
	return nil, nil
}

func replaceNonHiraganaWithSquare(input string) string {
	var result strings.Builder
	for _, r := range input {
		if unicode.In(r, unicode.Hiragana) {
			result.WriteRune(r)
		} else {
			result.WriteRune('□')
		}
	}
	return result.String()
}
