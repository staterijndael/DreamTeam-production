package views

import (
	"encoding/json"
)

type JSONRPCNotification struct {
	Namespace string
	Method    string
	Params    interface{}
}

func (n *JSONRPCNotification) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if n.Namespace == "" {
		m["method"] = "notification." + n.Method
	} else {
		m["method"] = n.Namespace + "." + n.Method
	}

	m["jsonrpc"] = "2.0"
	m["params"] = n.Params
	return json.Marshal(m)
}
