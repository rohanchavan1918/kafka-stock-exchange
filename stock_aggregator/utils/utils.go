package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/rohanchavan1918/stock_aggregator/conf"
)

type Payload struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Types string `json:"type"`
	Text  Text   `json:"text,omitempty"`
}

type Text struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Verbatim bool   `json:"verbatim"`
}

func AlertAndPanic(err error) {
	// Sends an alert before panicking
	fmt.Println("conf.AppConfig.SlackUrl > ", conf.AppConfig.SlackUrl)
	if conf.AppConfig.SlackUrl == "" {
		fmt.Println("SLackurl not found, panicking without alert.")
		panic(err)
	}

	SendSlackAlert(err.Error())
	panic(err)
}

func SendSlackAlert(msg string) error {
	AlertTitle := fmt.Sprintf("[%s]", conf.AppConfig.ServiceName)
	jsonmap := Payload{
		Blocks: []Block{
			{
				Types: "section",
				Text: Text{
					Type: "mrkdwn",
					Text: AlertTitle,
				},
			},
			{
				Types: "section",
				Text: Text{
					Type: "mrkdwn",
					Text: msg,
				},
			},
		},
	}
	jsonBody, err := json.Marshal(jsonmap)
	if err != nil {
		log.Printf("error marshalling json map %+v \n", err)
		return err
	}
	responseBody, err := DoPostReq(conf.AppConfig.SlackUrl, []byte(jsonBody))
	if err != nil {
		return err
	}
	fmt.Println("SLACK RESP > ", responseBody)

	return nil

}

func DoPostReq(url string, requestBody []byte) ([]byte, error) {
	// Create a request with POST method and request body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set the content type header (adjust as needed)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func LogInfo(msg string, args ...interface{}) {
	conf.AppConnections.Logger.Infof(msg, args...)
}

func LogError(msg string, args ...interface{}) {
	conf.AppConnections.Logger.Errorf(msg, args...)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
