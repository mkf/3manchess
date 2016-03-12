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
import "strconv"

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
		{"APIAddGame", "POST", "/api/addgame", mu.APIAddGame},
		{"APIPlay", "GET", "/api/play/{gameId}", mu.APIPlay},
		{"APITurn", "POST", "/api/play/{gameId}", mu.APITurn},
		{"APIState", "GET", "/api/state/{stateId}", mu.APIState},
		{"APIMove", "GET", "/api/move/{moveId}", mu.APIMove},
		{"APILogin", "POST", "/api/login", mu.APILogin},
		{"APIWhoIsIt", "GET", "/api/player/{playerId}", mu.APIWhoIsIt},
		{"APIUserInfo", "GET", "/api/user/{userId}", mu.APIUserInfo},
		{"APIBotInfo", "GET", "/api/bot/{botId}", mu.APIBotInfo},
		{"APIBotKey", "POST", "/api/botkey", mu.APIBotKey},
	}
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}

type LoggingIn struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
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
		return
	}
	log.Println("signing up: ", su)
	var gi SignUpGive
	gi.User, gi.Player, gi.Auth, err = mu.Server.SignUp(su.Login, su.Passwd, su.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(gi); err != nil {
		panic(err)
	}
}

func (mu *Multi) APILogin(w http.ResponseWriter, r *http.Request) {
	var li LoggingIn
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &li); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var aut Authorization
	aut.ID, aut.AuthKey, err = mu.Server.LogIn(li.Login, li.Passwd)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(aut); err != nil {
		panic(err)
	}
}

type BotKeyGetting struct {
	BotID    int64         `json:"botid"`
	UserAuth Authorization `json:"userauth"`
}

func (mu *Multi) APIBotKey(w http.ResponseWriter, r *http.Request) {
	var bkg BotKeyGetting
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &bkg); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	var aut Authorization
	aut.ID, aut.AuthKey, err = mu.Server.BotKey(bkg.BotID, bkg.UserAuth.ID, bkg.UserAuth.AuthKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(aut); err != nil {
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
	if err := json.Unmarshal(body, &nbp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	var nbg NewBotGive
	nbg.Botid, nbg.PlayerID, nbg.AuthKey, err =
		mu.Server.NewBot(nbp.WhoAmI, nbp.UserAuth.ID, nbp.UserAuth.AuthKey, nbp.OwnName, nbp.Settings)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
	if err := json.Unmarshal(body, &gpp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	var gpg GameplayGive
	gpg.Key, err = server.AddGame(mu.Server, &gpp.State, [3]*int64{gpp.White, gpp.Gray, gpp.Black}, gpp.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(gpg); err != nil {
		panic(err)
	}
}

type TurnPost struct {
	game.FromToProm `json:"fromtoprom"`
	WhoPlayer       Authorization `json:"whoplayer"`
}

type MoveAndAfterKeys struct {
	MoveKey      int64 `json:"move"`
	AfterGameKey int64 `json:"after"`
}

func (mu *Multi) APITurn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var turnp TurnPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &turnp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	if ok, err := mu.Server.PAuth(turnp.WhoPlayer.ID, turnp.WhoPlayer.AuthKey); !(ok && err == nil) {
		if !ok {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(false); err != nil {
				panic(err)
			}
		} else if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		return
	}
	var maak MoveAndAfterKeys
	ourint, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	maak.MoveKey, maak.AfterGameKey, err = server.MoveGame(mu.Server, ourint, turnp.FromToProm, turnp.WhoPlayer.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(maak); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIPlay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp server.GameplayData
	key, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	err = mu.Server.LoadGP(key, &gp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp game.StateData
	key, err := strconv.ParseInt(vars["stateId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	err = mu.Server.LoadSD(key, &gp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp server.MoveData
	key, err := strconv.ParseInt(vars["moveId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	err = mu.Server.LoadMD(key, &gp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

type InfoWhoIsIt struct {
	ID       int64 `json:"id"`
	IsItABot bool  `json:"isitabot"`
}

func (mu *Multi) APIWhoIsIt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var iwit InfoWhoIsIt
	key, err := strconv.ParseInt(vars["playerId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	iwit.ID, iwit.IsItABot, err = mu.Server.WhoIsIt(key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(iwit); err != nil {
		panic(err)
	}
}

type InfoUser struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Player int64  `json:"playerid"`
}

func (mu *Multi) APIUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var iu InfoUser
	key, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	iu.Login, iu.Name, iu.Player, err = mu.Server.UserInfo(key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(iu); err != nil {
		panic(err)
	}
}

type InfoBot struct {
	WhoAmI   []byte `json:"whoami"`
	Owner    int64  `json:"ownerid"`
	OwnName  string `json:"ownname"`
	Player   int64  `json:"playerid"`
	Settings []byte `json:"settings"`
}

func (mu *Multi) APIBotInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sib InfoBot
	key, err := strconv.ParseInt(vars["botId"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	sib.WhoAmI, sib.Owner, sib.OwnName, sib.Player, sib.Settings, err = mu.Server.BotInfo(key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(421)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sib); err != nil {
		panic(err)
	}
}
