package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

const ()

var (
	configFile string
	router     *gin.Engine
	err        error
	apiv       string
	pkg        Package
)

// Package config holds the application configuration settings.
type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Repository  string `json:"repository"`
	License     string `json:"license"`
	IsPrivate   bool   `json:"private"`
}

type qrData struct {
	URL  string `json:"url" binding:"required"`
	Size int    `json:"size"`
}

// init ...
func init() {
	pkgFile, err := os.Open("./package.json")
	// If there is a problem with the file, err on the side of caution and
	// bail out.
	if err != nil {
		fmt.Printf("error: Unable to open package.json/%s\n", configFile, err.Error())
		os.Exit(1)
	}
	defer pkgFile.Close()

	// Decode the json into something we can process.  The JSON is set up to load
	// into a map.  We could also do an array and move it to a map, but why?
	decoder := json.NewDecoder(pkgFile)
	err = decoder.Decode(&pkg)
	if err != nil {
		fmt.Printf("error: Could not decode JSON format in package.json/%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("info: Loaded package.json from disk\n")

	// Set the API version.
	apiv = pkg.Version
} // func

func cors() gin.HandlerFunc {
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
	router.Use(cors())

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
	log.Printf("info: Starting Centric QR Code Generator version %s...\n", apiv)
	router.Run(":3022")
} // main
