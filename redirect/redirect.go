package redirect

import (
	"log"
	"net/http"
)

const serverURL = "https://schedulebot-x2gm2h2g4a-uc.a.run.app"
const serverPort = "8080"

// Redirect redirects content from cloud fn hook to cloud run server
func Redirect(w http.ResponseWriter, r *http.Request) {
	_, err := http.Post(serverURL, "application/json", r.Body)
	if err != nil {
		log.Fatalln(err)
	}
}
