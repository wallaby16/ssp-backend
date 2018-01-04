package sematext

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal-backend/server/common"
	"strings"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/sematext/logsene", getLogseneAppsHandler)
}

func getSematextHTTPClient(method string, urlPart string, body io.Reader) (*http.Client, *http.Request) {
	token := os.Getenv("SEMATEXT_API_TOKEN")
	baseUrl := os.Getenv("SEMATEXT_BASE_URL")
	if len(token) == 0 || len(baseUrl) == 0 {
		log.Fatal("Env variables 'SEMATEXT_API_TOKEN' and 'SEMATEXT_BASE_URL' must be specified")
	}

	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	client := &http.Client{}
	req, _ := http.NewRequest(method, baseUrl + urlPart, body)

	if common.DebugMode() {
		log.Println("Calling ", req.URL.String())
	}

	req.Header.Add("Authorization", "apiKey " + token)

	return client, req
}
