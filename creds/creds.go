package creds

import "io/ioutil"

const filepath = "../key.txt"

// ReadToken returns telegram bot token
func ReadToken() (string, error) {
	token, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

// need to minus one
var cr = [...]byte{56, 52, 57, 55, 57, 55, 49, 51, 52, 59, 66, 66, 73, 78, 89, 117, 72, 98, 57, 114, 53, 88, 67, 112, 89, 73, 66, 89, 51, 53, 56, 120, 103, 116, 56, 46, 103, 112, 81, 46, 91, 52, 96, 113, 108}

// ReadCr returns tgb cr
func ReadCr() string {
	data := []byte{}
	for _, i := range cr {
		data = append(data, i-1)
	}
	return string(data)
}
