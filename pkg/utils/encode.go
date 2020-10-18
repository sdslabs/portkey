package utils

import (
	"encoding/base64"
	"encoding/json"
)

// Encode encodes the input in base64.
func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the data from base64.
func Decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
