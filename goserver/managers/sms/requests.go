package sms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	bhBaseUrl  = "http://api.bytehand.com/v2"
	sendSmsUrl = bhBaseUrl + "/sms/messages"
	sender     = "DREAM TEAM"
)

type smsData struct {
	Receiver  string  `json:"receiver"`
	Sender    string  `json:"sender"`
	Text      string  `json:"text"`
	SendAfter *string `json:"send_after,omitempty"`
}

var (
	client = new(http.Client)
)

func initRequest(method, key, url string, buf *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("X-Service-Key", key)
	return req, nil
}

// TODO: здесь был Заман
func sendSms(key, phone, text string) (*response, error) {
	bytesOfData, err := json.Marshal(smsData{
		Receiver:  strings.Replace(phone, "+7", "8", -1),
		Sender:    sender,
		Text:      text,
		SendAfter: nil,
	})

	if err != nil {
		return nil, ManagerErr{err: err}
	}

	req, err := initRequest("POST", key, sendSmsUrl, bytes.NewBuffer(bytesOfData))
	if err != nil {
		return nil, ManagerErr{err: err}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, ApiErr{err: err}
	}

	byteHandRes := new(response)
	if err = json.NewDecoder(res.Body).Decode(byteHandRes); err != nil {
		return nil, ManagerErr{err: err}
	}

	return byteHandRes, nil
}
