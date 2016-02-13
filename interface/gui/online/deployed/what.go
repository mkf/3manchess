package deployed

import "net/http"
import "github.com/ArchieT/3manchess/interface/gui/online"

func init() {
	http.HandleFunc("/", online.MainPage)
	http.HandleFunc("/play", online.PlayPage)
	http.HandleFunc("/move", online.MovePage)
}
