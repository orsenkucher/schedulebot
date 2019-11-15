package fbclient

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

//link :="https://us-central1-scheduleuabot.cloudfunctions.net/FetchUsersSubs"

//FetchUsersSubs is public
func FetchUsersSubs(ID string) [][]string {
	resp, err := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/FetchUsersSubs", "text/plain", bytes.NewBuffer([]byte(ID)))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return ParseSchedules(string(r))
}

// ParseSchedules is public
func ParseSchedules(str string) [][]string {
	nodes := [][]string{}
	i := 1
	path := []string{}
	for i < len(str)-1 {
		if str[i] == '[' {
			path = []string{}
		}
		if str[i] == ']' {
			nodes = append(nodes, path)
		}
		if str[i] == '"' {
			i++
			path = append(path, "")
			for str[i] != '"' {
				path[len(path)-1] += string([]byte{str[i]})
				i++
			}
		}
		i++
	}
	return nodes
}
