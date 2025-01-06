package util

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type HTTPResponse[D any] struct {
	Code    int64  `json:"code"`
	Data    *D     `json:"data"`
	Message string `json:"message"`
}

func DecodeResponse[D any](data []byte) (*D, error) {
	var resp HTTPResponse[D]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.Errorf("response error, code:%d, msg:%s", resp.Code, resp.Message)
	}
	return resp.Data, nil
}
