package main

import (
	"ataYemek/database"
	"ataYemek/models"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

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

	c.OnRequest(func(r *colly.Request) {

	})

	c.Visit("https://birimler.atauni.edu.tr/saglik-kultur-ve-spor-daire-baskanligi/")

	cnnt := database.Connect()
	clt := cnnt.Database("menuListe").Collection("menu")

	menuDB := models.MenuScrap{}

	anaMenu := models.Menu{}
	konukEviMenu := models.Menu{}
	menuData := []string{}
	cout := 0

	//Reversing slice
	for i, j := 0, len(tarihData)-1; i < j; i, j = i+1, j-1 {
		tarihData[i], tarihData[j] = tarihData[j], tarihData[i]
	}

	// yemek urlsinde yemekleri çeken döngü
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {

				if len(el.Text) == 0 || el.Text == "Gram" || el.Text == "KONUK EVİ 2 MENÜ" || el.Text == "ANA MENÜ" {
					return
				} else {
					if len(menuData) < 17 {

						menuData = append(menuData, el.Text)
						if len(menuData) == 16 { // eğer bir menu için 16 tane yemek varsa
							clt.FindOne(context.TODO(), bson.M{"tarih": tarihData[cout]}).Decode(&menuDB) // db'de bu tarihte kayıtlı bir menu varmı

							if len(menuDB.Tarih) == 0 { // eğer yoksa
								// menu yerlere göre ayırıyoruz
								for k, menu := range menuData {

									if k%4 == 1 {
										anaMenu.Gram = append(anaMenu.Gram, menu)
									} else if k%4 == 2 {
										anaMenu.Name = append(anaMenu.Name, menu)
									} else if k%4 == 3 {
										konukEviMenu.Gram = append(konukEviMenu.Gram, menu)
									} else if k%4 == 0 {
										konukEviMenu.Name = append(konukEviMenu.Name, menu)
									}

								}
								menuData = []string{}
								// database'e kaydetme
								mainData := models.MenuScrap{Tarih: tarihData[cout], Menuler: [2]models.Menu{anaMenu, konukEviMenu}}
								_, err := clt.InsertOne(context.TODO(), mainData)
								if err != nil {
									log.Fatal("Hata : " + err.Error())
								}
								menuDB = models.MenuScrap{}
							} else {
								fmt.Println("\nbu tarih daha önce kaydedilmiş")

							}
							cout++

						}
					}

				}
			})
		})
	})
	for _, url := range menurl {

		c.Visit(url)

	}

}

// belli periyotlarında yemek listesini internetten çekip mongodb'e kaydeten fonksiyon
func tickerForScraping() {
	ticker := time.NewTicker(10 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				getYemekListesi()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	<-quit
}
