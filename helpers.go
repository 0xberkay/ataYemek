package main

import (
	"context"
	"time"

	"github.com/0xberkay/ataYemek/database"
	"github.com/0xberkay/ataYemek/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
)

// api ana fonksiyonu
func runner() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	Setup(app)

	app.Listen("127.0.0.1:3000")
}

// bugünki yemekleri çeken fonksiyon
func bugunkiYemekler(ctx *fiber.Ctx) error {
	cnnt := database.Connect()
	clt := cnnt.Database("menuListe").Collection("menu")

	bugunTarih := time.Now()
	tarihKiyas := bugunTarih.Format("02.01.2006")

	var menuDB models.MenuScrap

	clt.FindOne(context.TODO(), bson.M{"tarih": tarihKiyas}).Decode(&menuDB) // tarih kıyası yapılıyor

	if len(menuDB.Tarih) == 0 { // eğer tarih kıyası yapılmamışsa
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Bu tarihte menu yok",
		})

	} else { // eğer tarih kıyası yapılmışsa
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Menu listesi",
			"menu":    menuDB,
		})
	}

}

// yarinki yemekleri çeken fonksiyon
func yarinkiYemekler(ctx *fiber.Ctx) error {
	cnnt := database.Connect()
	clt := cnnt.Database("menuListe").Collection("menu")

	// yarinki tarih
	yarinkitarih := time.Now().AddDate(0, 0, 1)
	tarihKiyas := yarinkitarih.Format("02.01.2006")

	var menuDB models.MenuScrap

	clt.FindOne(context.TODO(), bson.M{"tarih": tarihKiyas}).Decode(&menuDB) // tarih kıyası yapılıyor

	if len(menuDB.Tarih) == 0 { // eğer tarih kıyası yapılmamışsa
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Bu tarihte menu yok",
		})

	} else { // eğer tarih kıyası yapılmışsa
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Menu listesi",
			"menu":    menuDB,
		})
	}

}

// Tüm yemekleri çeken fonksiyon
func tumYemekler(ctx *fiber.Ctx) error {
	cnnt := database.Connect()
	clt := cnnt.Database("menuListe").Collection("menu")

	var tumMenu []models.MenuScrap

	// tüm yemekleri çekiyoruz
	cur, _ := clt.Find(context.TODO(), bson.M{})

	cur.All(context.TODO(), &tumMenu)

	if len(tumMenu) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Menu yok",
		})

	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tum Menuler",
			"menu":    tumMenu,
		})
	}
}
