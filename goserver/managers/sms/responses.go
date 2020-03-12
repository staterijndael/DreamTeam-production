package sms

type response struct {
	Result string  `json:"result"`
	Count  *int64  `json:"count, omitempty"`
	ID     *string `json:"id, omitempty"`
}
