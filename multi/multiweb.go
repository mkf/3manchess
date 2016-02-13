package multi

//import "github.com/ArchieT/3manchess/game"
//import "github.com/ArchieT/3manchess/player"
import "net/http"
import "github.com/ArchieT/3manchess/server"

//import "golang.org/x/net/context"
//import "html/template"
import "time"

//import "io"
import "log"
import "fmt"
import "github.com/gorilla/mux"

type GameInfo struct {
	Id   int64     `json:"id"`
	Date time.Time `json:"date"`
}

type DataProvider interface {
	GameList() []GameInfo
}

type Multi struct {
	DataProvider
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
	{"APIAuth", "POST", "/api/auth", APIAuth},
	{"APICreate", "POST", "/api/play", APICreate},
	{"APIPlay", "GET", "/api/play/{gameId}", APIPlay},
	{"APIMove", "POST", "/api/play/{gameId}", APIMove},
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
