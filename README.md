# Atauni Yemek listesi
Atatürk üniversitesi Sağlık Kültür Ve Spor Daire Başkanlığına bağlı yemekhanelerin yemek listesi apisi <br>

# Api

app.Get("/api/bugun", bugunkiYemekler) //Bugunki Yemekleri döndürür <br>
app.Get("/api/yarin", yarinkiYemekler) //Yarınki Yemekleri döndürür <br>
app.Get("/api/tum", tumYemekler)       //Tüm Yemekleri döndürür <br>


# Scrapper.go
https://birimler.atauni.edu.tr/saglik-kultur-ve-spor-daire-baskanligi/ sitesine gider <br>
eğer tarihle başlık atılmış gönderi varsa onu databasele karşılaştırıp kendi veri setine ekler<br>
