package lyrics

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

func GetLyrics() {
	artistName := "スリーズブーケ"
	numTracks := 5 // 取得したい曲名の数を指定
	tracks, err := GetArtistTopTracksBySpotifyAPI(artistName, numTracks)
	if err != nil {
		log.Fatal(err)
	}

	// 結果を出力
	for _, track := range tracks {
		// URLエンコードされた曲名を使用
		webPage := "https://www.uta-net.com/search/?Keyword=" + url.QueryEscape(track)
		resp, err := http.Get(webPage)
		if err != nil {
			log.Printf("failed to get html: %s", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
			continue
		}

		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			log.Printf("failed to load html: %s", err)
			continue
		}

		// タイトルを取得
		title := htmlquery.FindOne(doc, "//title")
		if title != nil {
			fmt.Println("Title:", htmlquery.InnerText(title))
		} else {
			fmt.Println("Title not found")
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
					lyrics := []string{}
					for _, part := range lyricsParts {
						// HTMLタグを除去してテキスト部分だけを取得
						partNode, err := htmlquery.Parse(strings.NewReader(part))
						if err != nil {
							log.Printf("failed to parse lyrics part: %s", err)
							continue
						}
						textPart := htmlquery.InnerText(partNode)
						lyrics = append(lyrics, textPart)
					}
					fmt.Println("Lyrics Parts[0]:", lyrics[0])
				} else {
					fmt.Println("Linked Page Lyrics Element not found")
				}
			}
		} else {
			fmt.Println("Element not found")
		}
	}
}
