package deployed

import "net/http"
import "github.com/ArchieT/3manchess/interface/online"

func init() {
	http.HandleFunc("/", online.Handler)
}
