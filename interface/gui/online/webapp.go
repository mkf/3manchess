package online

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"

//import "github.com/ArchieT/3manchess/interface/gui"
import "net/http"
import "golang.org/x/net/context"
import "google.golang.org/appengine"
import "html/template"
import "google.golang.org/appengine/user"
import "google.golang.org/appengine/datastore"

import "time"

var mainTemplate = template.Must(template.New("main").ParseFiles("static/main.html"))

func allGameplaysKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Gamesbase", "default_gamesbase", 0, nil)
}

type presentMain struct {
	LoginStr, LoginURL string
	Gameplays          []GameplayData
	GameKeys           []string
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	var pre presentMain
	var err error
	if u == nil {
		pre.LoginURL, _ = user.LoginURL(c, "/")
		pre.LoginStr = "Sign in or register to play"
	} else {
		pre.LoginURL, _ = user.LogoutURL(c, "/")
		pre.LoginStr = "Click here to sign out (logged as" + u.String() + ")"
	}
	pre.Gameplays = make([]GameplayData, 0, 20)
	q := datastore.NewQuery("Gameplay").Ancestor(allGameplaysKey(c)).Order("-Date").Limit(20)
	qk, err := q.GetAll(c, &pre.Gameplays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pre.GameKeys = make([]string, 0, len(qk))
	for iq := range qk {
		pre.GameKeys = append(pre.GameKeys, qk[iq].Encode())
	}
	if err := mainTemplate.Execute(w, pre); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PlayPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Deadline() //just a placeholder senseless and probably harmful, just delete and forget this line
}

func MovePage(w http.ResponseWriter, r *http.Request) {
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

func SavePlayer(pl player.Player, c context.Context) (*datastore.Key, error) {
	s := pl.Data()
	return datastore.Put(c, datastore.NewIncompleteKey(c, "Player", nil), &s)
}

type GameplayData struct {
	State, White, Gray, Black *datastore.Key
	Date                      time.Time
}

func SaveGameplay(gp player.Gameplay, c context.Context) (*datastore.Key, error) {
	st, err := SaveState(gp.State, c)
	if err != nil {
		return nil, err
	}
	w, err := SavePlayer(gp.Players[game.White], c)
	if err != nil {
		return nil, err
	}
	g, err := SavePlayer(gp.Players[game.Gray], c)
	if err != nil {
		return nil, err
	}
	b, err := SavePlayer(gp.Players[game.Black], c)
	if err != nil {
		return nil, err
	}
	d := GameplayData{State: st, White: w, Gray: g, Black: b, Date: time.Now()}
	return datastore.Put(c, datastore.NewIncompleteKey(c, "Gameplay", allGameplaysKey(c)), &d)
}
