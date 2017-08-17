package ddc

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"encoding/csv"
	"errors"
	"strconv"
	"strings"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
)

const apiErrorDDC = "Fehler beim Aufruf der DDC-API. Bitte erstelle ein Ticket bei DDC."

// RegisterRoutes registers the routes for OpenShift
func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/ddc/billing", getDDCBillingHandler)
}

func getDDCBillingHandler(c *gin.Context) {
	rows, err := calculateDDCBilling()
	result := createCSVReport(rows)

	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, err.Error())
	}
}

func createCSVReport(rows []common.DDCBillingRow) common.DDCBilling {
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.Comma = ';'

	for _, r := range rows {
		// Example:
		// TODO
		totalString := strconv.FormatFloat(r.Total, 'f', 6, 64)
		row := []string{"", r.Sender, "","","","","","", r.Art, "", totalString, "", "CHF", "",
			r.Assignment1, "", r.Assignment2, "", r.Assignment3, "", "", "", "", "", "1", "","ST", "", "", "", "LM1706 DDC ", r.Project + ": " + r.Host}
		wr.Write(row)
	}

	wr.Flush()

	return common.DDCBilling{
		CSV: b.String(),
		Rows: rows,
	}
}

func calculateDDCBilling() ([]common.DDCBillingRow, error) {
	client, req := getDDCClient()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling ddc api: ", err.Error())
		return nil, errors.New(apiErrorDDC)
	} else {
		defer resp.Body.Close()
	}

	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Println("Error parsing CSV.", err.Error())
		return nil, errors.New(apiErrorDDC)
	}

	const sender = "70029508"
	const art = "816753"
	const fee_server = 300.0
	const fee_client = 100.0
	const feeCpu = 30
	const feeMemory = 30
	const feeStorage = 1.0

	result := []common.DDCBillingRow{}
	for i, value := range records {
		if i > 0 {
			usedCpu, _ := strconv.ParseFloat(value[2], 64)
			usedMemory, _ := strconv.ParseFloat(value[3], 64)
			usedStorage, _ := strconv.ParseFloat(value[4], 64)

			var fee float64
			if strings.Contains(value[10], "Windows 7") {
				fee = fee_client
			} else {
				fee = fee_server
			}

			totalPrice := fee + usedCpu * feeCpu + usedMemory * feeMemory + usedStorage * feeStorage
			result = append(result, common.DDCBillingRow{
				Sender: sender,
				Art:    art,
				Assignment1: value[5],
				Assignment2: value[6],
				Assignment3: value[7],
				Total: totalPrice,
				TotalCPU: usedCpu * feeCpu,
				TotalMemory: usedMemory * feeMemory,
				TotalStorage: usedStorage * feeStorage,
				Project: value[1],
				Host: value[0],
			})
		}
	}

	return result, nil
}

func getDDCClient() (*http.Client, *http.Request) {
	api := os.Getenv("DDC_API")
	if len(api) == 0 {
		log.Fatal("Env variable 'DDC_API' must be specified")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", api, nil)

	if common.DebugMode() {
		log.Print("Calling ", req.URL.String())
	}

	return client, req
}
