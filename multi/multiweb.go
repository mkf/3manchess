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

func Run() {
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

var routes = Routes{
	{"Index", "GET", "/", Index},
	{"APIIndex", "GET", "/api", APIIndex},
	{"APISignUp", "POST", "/api/signup", APISignUp},
	{"APIAuth", "POST", "/api/newbot", APINewBot},
	{"APICreate", "POST", "/api/addgame", APIAddGame},
	{"APIPlay", "GET", "/api/play/{gameId}", APIPlay},
	{"APITurn", "POST", "/api/play/{gameId}", APITurn},
	{"APIState", "GET", "/api/state/{stateId}", APIState},
	{"APIMove", "GET", "/api/move/{moveId}", APIMove},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, router := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}

type SignUp struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
	Name   string `json:"name"`
}

func APIIndex(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Index is here") }

func APISignUp(w http.ResponseWriter, r *http.Request) {
	var su SignUp
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
}
