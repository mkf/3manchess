package multi

import "github.com/ArchieT/3manchess/game"

//import "github.com/ArchieT/3manchess/player"
import "net/http"
import "github.com/ArchieT/3manchess/server"

//import "golang.org/x/net/context"
import "time"
import "encoding/json"
import "io"
import "io/ioutil"
import "log"
import "fmt"
import "github.com/gorilla/mux"

type Multi struct {
	server.Server
}

func (mu *Multi) Run() {
	router := mu.NewRouter()
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
			inner.ServeHTTP(w, r)
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
	for _, route := range routes {
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &su); err != nil {
		w.WriteHeader(422) //unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var gi SignUpGive
	uu, pp, aa, ee := mu.Server.SignUp(su.Login, su.Passwd, su.Name)
	if err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(ee); err != nil {
			panic(err)
		}
	}
	gi.User = uu
	gi.Player = pp
	gi.Auth = aa
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(gi); err != nil {
		panic(err)
	}
}

type Authorization struct {
	ID      int64  `json:"id"`
	AuthKey []byte `json:"authkey"`
}

type NewBotPost struct {
	WhoAmI   []byte        `json:"whoami"`
	UserAuth Authorization `json:"owner"`
	OwnName  string        `json:"ownname"`
	Settings []byte        `json:"settings"`
}

type NewBotGive struct {
	Botid    int64
	PlayerID int64
	AuthKey  []byte
}

func (mu *Multi) APINewBot(w http.ResponseWriter, r *http.Request) {
	var nbp NewBotPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &nbp); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var nbg NewBotGive
	nbg.Botid, nbg.PlayerID, nbg.AuthKey, err =
		mu.Server.NewBot(nbp.WhoAmI, nbp.UserAuth.ID, nbp.UserAuth.AuthKey, nbp.OwnName, nbp.Settings)
	if err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(nbg); err != nil {
		panic(err)
	}
}

type GameplayPost struct {
	State game.StateData `json:"state"`
	White *int64         `json:"whiteplayer"`
	Gray  *int64         `json:"grayplayer"`
	Black *int64         `json:"blackplayer"`
	Date  time.Time      `json:"when"`
}

type GameplayGive struct {
	Key int64 `json:"gameid"`
}

func (mu *Multi) APIAddGame(w http.ResponseWriter, r *http.Request) {
	var gpp GameplayPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &gpp); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var gpg GameplayGive
	gpg.Key, err = server.AddGame(mu.Server, &gpp.State, [3]*int64{gpp.White, gpp.Gray, gpp.Black}, gpp.Date)
	if err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(gpg); err != nil {
		panic(err)
	}
}

type TurnPost struct {
	Before          int64 `json:"beforegame"`
	game.FromToProm `json:"fromtoprom"`
	WhoPlayer       Authorization `json:"whoplayer"`
}

type MoveAndAfterKeys struct {
	MoveKey      int64 `json:"move"`
	AfterGameKey int64 `json:"after"`
}

func (mu *Multi) APITurn(w http.ResponseWriter, r *http.Request) {
	var turnp TurnPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	w.Header().Set("Content=Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &turnp); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	if ok, err := mu.Server.PAuth(turnp.WhoPlayer.ID, turnp.WhoPlayer.AuthKey); !(ok && err == nil) {
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(false); err != nil {
				panic(err)
			}
		} else if err != nil {
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
	}
	var maak MoveAndAfterKeys
	maak.MoveKey, maak.AfterGameKey, err = server.MoveGame(mu.Server, turnp.Before, turnp.FromToProm, turnp.Who.ID)
	if err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(maak); err != nil {
		panic(err)
	}
}
