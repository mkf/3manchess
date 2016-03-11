package multi

//import "github.com/ArchieT/3manchess/game"
//import "github.com/ArchieT/3manchess/player"
import "net/http"
import "github.com/ArchieT/3manchess/server"

//import "golang.org/x/net/context"
//import "html/template"
import "time"

import "fmt"
import "io"
import "io/ioutil"
import "log"
import "fmt"
import "github.com/gorilla/mux"

type Multi struct {
	server.Server
}

func (mu *Multi) Run() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8082", router))
}

type Route struct {
	Name, Method, Pattern string
	http.HandlerFunc
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			inner.ServerHTTP(w, r)
			log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
		})
}

type Routes []Route

func (mu *Multi) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	var routes = Routes{
		{"Index", "GET", "/", mu.APIIndex},
		{"APIIndex", "GET", "/api", mu.APIIndex},
		{"APISignUp", "POST", "/api/signup", mu.APISignUp},
		{"APIAuth", "POST", "/api/newbot", mu.APINewBot},
		{"APICreate", "POST", "/api/addgame", mu.APIAddGame},
		{"APIPlay", "GET", "/api/play/{gameId}", mu.APIPlay},
		{"APITurn", "POST", "/api/play/{gameId}", mu.APITurn},
		{"APIState", "GET", "/api/state/{stateId}", mu.APIState},
		{"APIMove", "GET", "/api/move/{moveId}", mu.APIMove},
	}
	for _, router := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}

func (mu *Multi) APIIndex(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Index is here") }

type SignUpPost struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
	Name   string `json:"name"`
}

type SignUpGive struct {
	User   int64  `json:"userid"`
	Player int64  `json:"playerid"`
	Auth   []byte `json:"authkey"`
}

func (mu *Multi) APISignUp(w http.ResponseWriter, r *http.Request) {
	var su SignUpPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &su); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) //unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var gi SignUpGive
	uu, pp, aa, ee := mu.Server.SignUp(su.Login, su.Passwd, su.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(ee); err != nil {
			panic(err)
		}
	}
	gi.User = uu
	gi.Player = pp
	gi.Auth = aa
	//ciag dalszy, wg poradnika restful thenewstack
}
