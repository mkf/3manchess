package deployed

import "net/http"
import "github.com/ArchieT/3manchess/interface/gui/online"

func init() {
	og := new(online.Online)
	og.Initialize()
	http.HandleFunc("/", og.MainPage)
	http.HandleFunc("/play", og.PlayPage)
	http.HandleFunc("/move", og.MovePage)
}
