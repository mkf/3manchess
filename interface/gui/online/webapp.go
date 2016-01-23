package online

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/interface/gui"
import "fmt"
import "net/http"
import "google.golang.org/appengine"
import "google.golang.org/appengine/user"
import "google.golang.org/appengine/datastore"

func Handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
}
