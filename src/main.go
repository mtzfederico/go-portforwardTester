package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PortScanResult struct {
	Port int
	Open bool
}

func getRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"ClientIP": c.ClientIP()})
}

func handleCheckPorts(c *gin.Context) {
	host := c.PostForm("host")
	host = strings.ReplaceAll(host, " ", "")
	portsFormVal := c.PostForm("ports")
	portsFormVal = strings.ReplaceAll(portsFormVal, " ", "")

	errorMessage := ""

	if host == "" { // make sure it is a valid hostname
		errorMessage = "No host specified"
	}

	if portsFormVal == "" { // make sure that it is a valid port or multiple ports
		errorMessage = "No ports specified"
	}

	ports := strings.Split(portsFormVal, ",")

	results := []PortScanResult{}

	if errorMessage == "" {
		for _, port := range ports {
			timeout := time.Second
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
			if err != nil {
				// fmt.Println("Connection error: ", err)
			}
			portAsInt, err := strconv.Atoi(port)
			if err != nil {
				fmt.Println("Atoi error: ", err)
			}

			if conn != nil {
				defer conn.Close()
				fmt.Println("Opened", net.JoinHostPort(host, port))

				results = append(results, PortScanResult{portAsInt, true})
			} else {
				results = append(results, PortScanResult{portAsInt, false})
			}
		}
	}

	acceptEncoding := c.Request.Header.Get("Accept")
	if acceptEncoding == "application/json" {
		c.JSON(http.StatusOK, gin.H{"host": host, "results": results})
	} else {
		clientIP := c.ClientIP()
		c.HTML(http.StatusOK, "result.tmpl", gin.H{
			"ClientIP":        clientIP,
			"Host":            host,
			"PortsFormVal":    portsFormVal,
			"ErrorMessage":    errorMessage,
			"PortScanResults": results,
		})
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	fmt.Println("Starting server")
	router := gin.Default()
	// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
	// srouter.SetTrustedProxies([]string{"127.0.0.1"})
	// router.TrustedPlatform = "X-CDN-IP"

	router.LoadHTMLGlob("templates/*.tmpl")

	router.StaticFile("/style.css", "static/style.css")
	router.GET("/", getRoot)
	router.POST("/", handleCheckPorts)

	fmt.Println("Listening on port 8080")
	router.Run(":8080")
}
