package models

type Yemek struct {
	Name string
	Gram string
}

type MenuScrap struct {
	Tarih   string
	Menuler []Yemek
}
