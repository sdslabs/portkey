package utils

import (
	"errors"
)

func Validate(key string, sendPath string, receive bool, receivePath string, certPath string, privateKeyPath string, doBenchmarking bool) error {
	if !receive && sendPath == "" {
		return errors.New("invalid flags combination: specify atleast one of -r or -s <sendpath>")
	}

	certSpecified := (certPath != "")
	privateKeySpecified := (privateKeyPath != "")
	if certSpecified != privateKeySpecified {
		return errors.New("invalid flags combination: specify both -c and -p or neither")
	}

	return nil
}
