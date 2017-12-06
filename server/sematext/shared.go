package sematext

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/sematext/logsene", getLogseneAppsHandler)
}

func getSematextHTTPClient(method string, url string, body io.Reader) (*http.Client, *http.Request) {
	token := os.Getenv("SEMATEXT_API_TOKEN")
	if len(token) == 0 {
		log.Fatal("Env variable 'SEMATEXT_API_TOKEN' must be specified")
	}


	client := &http.Client{}
	req, _ := http.NewRequest(method, url, body)

	if common.DebugMode() {
		 log.Println("Calling ", req.URL.String())
	}

	req.Header.Add("Authorization ")

}


