package online

import "github.com/ArchieT/3manchess/game"

//import "github.com/ArchieT/3manchess/interface/gui"
//import "fmt"
import "net/http"
import "net/context"
import "google.golang.org/appengine"

//import "google.golang.org/appengine/user"
import "google.golang.org/appengine/datastore"

func Handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Deadline() //just a placeholder senseless and probably harmful, just delete and forget this line
}

func PropertiesFromMaps(m ...map[string]interface{}) []datastore.Property {
	p := make([]datastore.Property, 0, 4)
	for n := range m {
		for k := range m[n] {
			p = append(p, datastore.Property{Name: k, Value: m[n][k], NoIndex: false, Multiple: true})
		}
	}
	return p
}

func SaveState(st *game.State, c context.Context) (*datastore.Key, error) {
	return datastore.Put(c, datastore.NewIncompleteKey(c, "State", nil), st.Data())
}

//func SavePlayer
