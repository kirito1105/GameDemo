package msg

import (
	"encoding/json"
	"net/http"
)

const (
	SUCCESS        = 0
	REGISTERED     = 1
	NOUSER         = 2
	PWDERR         = 3
	OUTTIMESESSION = 4
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Send(w *http.ResponseWriter, msg *Res) error {
	if err := json.NewEncoder(*w).Encode(*msg); err != nil {
		return err
	}
	return nil
}
