package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
	"github.com/qor/media"
	"github.com/qor/media/media_library"
	"github.com/qor/media/oss"
	"github.com/qor/qor/utils"
)

type Grid struct {
	gorm.Model
	Name  string
	Image oss.OSS
}
type Swiper struct {
	gorm.Model
	Name  string
	Image oss.OSS
}
type Color struct {
	gorm.Model
	Name string
}
type Size struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Name  string
	Image oss.OSS
}

// Create another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Pric        float64
	Description string
	Image       oss.OSS
	//	Images          media_library.MediaBox
	Category        Category
	CategoryID      uint
	ColorVariations []ColorVariation
}
type ColorVariation struct {
	gorm.Model
	Name           string
	ProductID      uint
	Color          Color
	ColorID        uint
	SizeVariations []SizeVariation
}
type SizeVariation struct {
	gorm.Model
	Name             string
	ColorVariationID uint
	Size             Size
	SizeID           uint
	Num              uint
}
type Category struct {
	gorm.Model
	Name string
}

func main() {
	DB, _ := gorm.Open("sqlite3", "tdwl.db")
	//DB.DropTableIfExists(&Color{}, &Product{})
	DB.AutoMigrate(&Swiper{}, &Grid{}, &Category{}, &Product{}, &ColorVariation{}, &SizeVariation{}, &Color{}, &Size{}, &media_library.MediaLibrary{})
	media.RegisterCallbacks(DB)

	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Allow to use Admin to manage User, Product
	// add resource
	Admin.AddResource(&Category{})
	Admin.AddResource(&Product{})
	Admin.AddResource(&ColorVariation{})
	Admin.AddResource(&SizeVariation{})
	Admin.AddResource(&Color{})
	Admin.AddResource(&Size{})
	//fields
	//	product.IndexAttrs("CategoryID", "Category", "Name")
	// scopes
	// initalize an HTTP request multiplexer
	//======================
	mux := http.NewServeMux()
	for _, path := range []string{"system", "javascripts", "stylesheets", "images"} {
		mux.Handle(fmt.Sprintf("/%s/", path), utils.FileServer(http.Dir("public")))
	}

	// Mount admin interface to mux
	Admin.MountTo("/admin", mux)
	fmt.Println("Listening on: 80")
	http.ListenAndServe(":80", mux)
}
