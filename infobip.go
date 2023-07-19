package go_sms_sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

type InfobipClient struct {
	baseUrl string
	sender  string
	apiKey  string
}

type InfobipConfigService struct {
	baseUrl string
	sender  string
	apiKey  string
}

type SmsService struct {
	configService InfobipConfigService
}

type MessageData struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	From         string        `json:"from"`
	Destinations []Destination `json:"destinations"`
	Text         string        `json:"text"`
}

type Destination struct {
	To string `json:"to"`
}

func GetInfobipClient(sender string, appKey string, other []string) (*InfobipClient, error) {
	fmt.Println("GetInfobipClient")
	return &InfobipClient{
		baseUrl:other[0],
		sender:sender,
		apiKey:appKey,
	}, nil
}

func (c *InfobipClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: msg code")
	}

	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missin parer: trgetPhoneNumber")
	}
	
	mobile := targetPhoneNumber[0]

	if strings.HasPrefix(mobile, "0") {
		mobile = "886" + mobile[1:]
	}
	if strings.HasPrefix(mobile, "+") {
		mobile = mobile[1:]
	}

	fmt.Println("mbile:"+mobile )

	
	endpoint := fmt.Sprintf("%s/sms/2/text/advanced", c.baseUrl)
	fmt.Println("endpoint: "+ endpoint)

	messageData := MessageData{
		Messages: []Message{
			{
				From: c.sender,
				Destinations: []Destination{
					{
						To: mobile,
					},
				},
				Text: code,
			},
		},
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("App %s", c.apiKey),
		"Content-Type":  "application/json",
	}

	fmt.Println("headers: ", headers)

	

	messageDataBytes, _ := json.Marshal(messageData)
	fmt.Println("messageDataBytes: ", messageDataBytes)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(messageDataBytes))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	dump, err := httputil.DumpRequestOut(req, true)
  if err != nil {
    return nil
  }



  fmt.Printf("%s\n", dump)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("sms error -> ", err)
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("resp:" + resp.Status)
	return  nil
}

