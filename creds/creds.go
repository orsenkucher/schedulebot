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
