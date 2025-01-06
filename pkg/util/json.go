package util

import (
	"bytes"
	"encoding/json"
)

func BeautifyJSON(compressedJSON string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(compressedJSON), "", "  ")
	if err != nil {
		return compressedJSON
	}
	return out.String()
}
