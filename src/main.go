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

func checkPort(host string, port string) (result PortScanResult) {
	timeout := time.Second
	hostWithPort := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", hostWithPort, timeout)
	if err != nil {
		// fmt.Println("Connection error: ", err)
	}
	portAsInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Atoi error: ", err)
	}

	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", hostWithPort)

		return PortScanResult{portAsInt, true}
	} else {
		return PortScanResult{portAsInt, false}
	}
}

func getRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"ClientIP": c.ClientIP()})
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
			// chcek if it is a port range (8880-8888)
			if strings.Contains(port, "-") {
				portRange := strings.Split(port, "-")
				min, _ := strconv.Atoi(portRange[0])
				max, _ := strconv.Atoi(portRange[1])

				if (max-min)+1 > 30 {
					errorMessage = "The range is too big. (Max 30)"
					continue
				}

				if max < min {
					errorMessage = "The start value is smaller than the end value"
					continue
				}

				for portInRange := min; portInRange <= max; portInRange++ {
					portInRangeAsString := strconv.Itoa(portInRange)
					result := checkPort(host, portInRangeAsString)
					results = append(results, result)
				}
				continue
			}

			result := checkPort(host, port)
			results = append(results, result)
		}
	}

	acceptEncoding := c.Request.Header.Get("Accept")
	if acceptEncoding == "application/json" {
		c.JSON(http.StatusOK, gin.H{"host": host, "error": errorMessage, "results": results})
	} else {
		clientIP := c.ClientIP()
		c.HTML(http.StatusOK, "result.html", gin.H{
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

	router.LoadHTMLGlob("templates/*.html")

	router.StaticFile("/style.css", "static/style.css")
	router.GET("/", getRoot)
	router.POST("/", handleCheckPorts)

	fmt.Println("Listening on port 8080")
	router.Run(":8080")
}
