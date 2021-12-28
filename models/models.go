package models

type Menu struct {
	Name []string
	Gram []string
}

type MenuScrap struct {
	Tarih   string
	Menuler [2]Menu
}
