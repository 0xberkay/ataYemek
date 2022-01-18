package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/0xberkay/ataYemek/database"
	"github.com/0xberkay/ataYemek/models"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
)

// Yemek listesini internetten çekip mongodb'e kaydeten fonksiyon
func getYemekListesi() {
	c := colly.NewCollector()
	menurl := []string{}
	tarihData := []string{}

	// Yemek listesi sayfasının url'si arıyan döngü
	c.OnHTML(".post-title", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			tarih := strings.TrimSpace(el.Attr("title"))

			// regexp date
			match := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)
			if match.MatchString(tarih) {
				tarihData = append(tarihData, tarih)
				menurl = append(menurl, el.Attr("href"))
			}

		})

	})

	c.Visit("https://birimler.atauni.edu.tr/saglik-kultur-ve-spor-daire-baskanligi/")

	cnnt := database.Connect()
	clt := cnnt.Database("menuListe").Collection("menu")

	menuDB := models.MenuScrap{}
	menuData := models.MenuScrap{}
	yemek := models.Yemek{}

	d := colly.NewCollector()

	for i := 0; i < len(menurl); i++ {

		cout := 0
		d.OnHTML("table", func(e *colly.HTMLElement) {

			if cout == 0 {

				e.ForEach("tr", func(_ int, e *colly.HTMLElement) {

					e.ForEach("td", func(_ int, el *colly.HTMLElement) {

						if el.Text != "ANA MENÜ" && el.Text != "Gram" {

							if el.Attr("colspan") != "" {
								yemek.Name = el.Text
							} else {
								yemek.Gram = el.Text
								menuDB.Menuler = append(menuDB.Menuler, yemek)
								yemek = models.Yemek{}

							}

						}

					})
				})

				menuDB.Tarih = tarihData[i]
				clt.FindOne(context.TODO(), bson.M{"tarih": menuDB.Tarih}).Decode(&menuData)

				if menuData.Tarih != "" {
					fmt.Println("Tarih daha önce kaydedildi")
				} else {
					_, err := clt.InsertOne(context.TODO(), menuDB)
					menuDB = models.MenuScrap{}

					if err != nil {
						panic(err)
					}
				}

				cout++
			}

		})

		d.Visit(menurl[i])

	}
}

// belli periyotlarında yemek listesini internetten çekip mongodb'e kaydeten fonksiyon
func tickerForScraping() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		fmt.Println("SCRAPER STARTED")
		for {
			select {
			case <-ticker.C:
				getYemekListesi()
				fmt.Println("fonksiyon çalıştı")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	<-quit
}
