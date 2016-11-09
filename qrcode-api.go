package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const ()

var (
	router *gin.Engine
)

type qrData struct {
	URL  string `json:"url" binding:"required"`
	Size int    `json:"size"`
}

// init ...
func init() {

} // func

// Cors ...
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
} // func

// SetupRouter ...
func SetupRouter() *gin.Engine {
	// Set up the router.
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// Global middleware
	//router.Use(gin.Logger())
	//router.Use(gin.Recovery())
	router.Use(Cors())

	// Now set up the routes.
	router.POST("/qrcode", MakeQR)
	router.GET("/ping", PingTheAPI)

	return router
}

// MakeQR ...
func MakeQR(c *gin.Context) {
	var parms qrData

	pr, pw := io.Pipe()
	x := new(bytes.Buffer)

	if c.BindJSON(&parms) == nil {
		qrcode, err := qr.Encode(parms.URL, qr.L, qr.Auto)
		if err != nil {
			fmt.Println(err)
		} else {
			if parms.Size < 25 || parms.Size > 300 {
				parms.Size = 200
			}
			qrcode, err = barcode.Scale(qrcode, parms.Size, parms.Size)
			if err != nil {
				fmt.Println(err)
			} else {
				go func() {
					defer pw.Close()
					err = png.Encode(pw, qrcode)
				}()
			}
		} // else

		//
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			x.ReadFrom(pr)
			c.Data(http.StatusOK, "image/png", x.Bytes())
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
	}
} // func

// PingTheAPI ...
func PingTheAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"payload": "PONG"})
} // func

//
func main() {
	// Start the server.
	router := SetupRouter()
	fmt.Printf("Starting...\n")
	router.Run(":3022")
} // main
