package sheets

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jaredallard/mm/pkg/config"
	log "github.com/jaredallard/mm/pkg/logger"
)

// SheetContents represents, to the best of my ability, the contents of a sheet.
type SheetContents struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

var apiKey string

// Init the client
func Init(cfg *config.ConfigurationFile) {
	apiKey = cfg.Google.APIKey
}

// flatten a two dimensional array into one dimension
func flatten(contents SheetContents) []string {
	alloc := make([]string, len(contents.Values))

	for i := range contents.Values {
		alloc[i] = contents.Values[i][0]
	}

	return alloc
}

// GetRange is a friendlier frontend to GetSheet, this is for single value ranges ONLY.
func GetRange(spreadsheetID string, readRange string) ([]string, error) {
	var response []string

	contents, err := GetSheet(spreadsheetID, readRange)
	if err != nil {
		return response, err
	}

	return flatten(contents), nil
}

// GetSheet gets a spread sheet and range
func GetSheet(spreadsheetID string, readRange string) (SheetContents, error) {
	contents := SheetContents{}

	url := "https://sheets.googleapis.com/v4/spreadsheets/" + spreadsheetID + "/values/" + readRange + "?key=" + apiKey
	log.Debug("Hit", url)

	resp, err := http.Get(url)
	if err != nil {
		return contents, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &contents)

	return contents, nil
}
