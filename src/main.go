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

/*
func getRoot(c *gin.Context) {
	acceptEncoding := c.Request.Header.Get("Accept")
	// fmt.Println("Accept: ", acceptEncoding)
	if acceptEncoding == "application/json" {
		// fmt.Println("Accept header is JSON")
		c.JSON(http.StatusOK, gin.H{"message": "Port Forward Tester Built in GO"})
	} else {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "",
		})
	}
}*/

func handleCheckPorts(c *gin.Context) {
	host := c.PostForm("host")
	portsFormVal := c.PostForm("ports")
	// c.String(http.StatusOK, "host: %s, ports: %s", host, ports)

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
				fmt.Println("Connecting error:", err)
				errorMessage = err.Error()
				break
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
	// fmt.Println("Accept: ", acceptEncoding)
	if acceptEncoding == "application/json" {
		c.JSON(http.StatusOK, gin.H{"results": results})
	} else {
		c.HTML(http.StatusOK, "result.tmpl", gin.H{
			"Host":            host,
			"PortsFormVal":    portsFormVal,
			"ErrorMessage":    errorMessage,
			"PortScanResults": results,
		})
	}

	// c.JSON(200, gin.H{"results": results})
	// {"results": [{"port":22,"status":"closed"},{"port":80,"status":"open"},{"port":443,"status":"open"}]}
}

func main() {
	router := gin.Default()
	// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
	// srouter.SetTrustedProxies([]string{"127.0.0.1"})
	// router.TrustedPlatform = "X-CDN-IP"

	router.LoadHTMLGlob("templates/*.tmpl")

	// router.GET("/", getRoot)
	router.StaticFile("/", "static/index.html")
	router.StaticFile("/style.css", "static/style.css")
	router.POST("/", handleCheckPorts)

	router.Run("localhost:8080")
}

// Take a look at your page's HTML
// specially for a footer
