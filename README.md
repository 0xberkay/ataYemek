# Atauni Yemek listesi
Atatürk üniversitesi Sağlık Kültür Ve Spor Daire Başkanlığına bağlı yemekhanelerin yemek listesi apisi

# Api

app.Get("berkay.one/api/bugun", bugunkiYemekler) //Bugunki Yemekleri döndürür
app.Get("berkay.one/api/yarin", yarinkiYemekler) //Yarınki Yemekleri döndürür
app.Get("berkay.one/api/tum", tumYemekler)       //Tüm Yemekleri döndürür


# Scrapper.go
https://birimler.atauni.edu.tr/saglik-kultur-ve-spor-daire-baskanligi/ sitesine gider <br>
eğer tarihle başlık atılmış gönderi varsa onu databasele karşılaştırıp kendi veri setine ekler<br>
