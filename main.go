package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/medivhzhan/weapp"
	"github.com/qor/admin"
	"github.com/qor/media"
	"github.com/qor/media/media_library"
	"github.com/qor/media/oss"
	"golang.org/x/crypto/acme/autocert"
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
	Admin.AddResource(&Swiper{})
	Admin.AddResource(&Grid{})
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
	r := gin.Default()
	r.GET("/login", Login)
	r.GET("/grids", GetGrids)
	r.GET("/swipers", GetSwipers)
	r.GET("/products", GetProducts)
	mux := http.NewServeMux()
	for _, path := range []string{"system", "javascripts", "stylesheets", "images"} {
		r.StaticFS(fmt.Sprintf("/%s", path), http.Dir(fmt.Sprintf("public/%s", path)))
		//	mux.Handle(fmt.Sprintf("/%s/", path), utils.FileServer(http.Dir("public")))
	}

	// Mount admin interface to mux
	Admin.MountTo("/admin", mux)
	//	fmt.Println("Listening on: 80")
	//	http.Handle("/", mux)
	//	http.ListenAndServe(":8080", mux)

	r.Any("/admin/*filepath", gin.WrapH(mux))
	//	r.Run(":80")

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("wcqt.site"),
		Cache:      autocert.DirCache("./.cache"),
	}

	log.Fatal(autotls.RunWithManager(r, &m))

}
func GetProducts(c *gin.Context) {
	DB, _ := gorm.Open("sqlite3", "tdwl.db")
	defer DB.Close()
	res := map[uint][]Product{}
	m := []Product{}
	DB.Find(&m)
	for _, p := range m {
		res[p.CategoryID] = append(res[p.CategoryID], p)
	}
	c.JSON(200, res)
}
func GetSwipers(c *gin.Context) {
	DB, _ := gorm.Open("sqlite3", "tdwl.db")
	defer DB.Close()
	m := []Swiper{}
	DB.Find(&m)

	c.JSON(200, m)
}
func GetGrids(c *gin.Context) {
	DB, _ := gorm.Open("sqlite3", "tdwl.db")
	defer DB.Close()
	grids := []Grid{}
	DB.Find(&grids)

	c.JSON(200, grids)
}
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pongxxx",
	})
}
func Login(c *gin.Context) {
	const (
		appID  = "wx5032f0d783147d67"
		secret = "6a194cc97598e54d76e43ac2fd632c3d"
	)
	//	firstname := c.DefaultQuery("firstname", "Guest")
	code := c.Query("code") // shortcut for c.Request.URL.Query().Get("lastname")
	res, err := weapp.Login(appID, secret, code)
	fmt.Println(res, err)
	//	c.String(http.StatusOK,"abc")
}
