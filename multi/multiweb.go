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
import "errors"
import "reflect"

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

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			w.Header().Set("Access-Control-Allow-Origin", "*")
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
		{"APIAfter", "GET", "/api/play/{gameId}/after", mu.APIAfter},
		{"APIBefore", "GET", "/api/play/{gameId}/before", mu.APIBefore},
		{"APIState", "GET", "/api/state/{stateId}", mu.APIState},
		{"APIVFTPGen", "GET", "/api/state/{stateId}/vftpgen", mu.APIVFTPGen},
		{"APIMove", "GET", "/api/move/{moveId}", mu.APIMove},
		{"APIDiff", "GET", "/api/move/{moveId}/diff", mu.APIDiff},
		{"APILogin", "POST", "/api/login", mu.APILogin},
		{"APIWhoIsIt", "GET", "/api/player/{playerId}", mu.APIWhoIsIt},
		{"APIUserInfo", "GET", "/api/user/{userId}", mu.APIUserInfo},
		{"APIOwnersBots", "GET", "/api/user/{userId}/bots", mu.APIOwnersBots},
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

type OurError struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Where   string `json:"where"`
}

func (oe *OurError) Error() string {
	a, e := json.Marshal(oe)
	if e != nil {
		panic(e)
	}
	return string(a)
}

type APIListErr struct {
	Errors []OurError `json:"errors"`
}

func (ale APIListErr) Empty() bool {
	return len(ale.Errors) == 0
}

func (ale APIListErr) Error() string {
	if ale.Empty() {
		return ""
	}
	return ale.Errors[0].Error()
}

func (ale APIListErr) ToErr() error {
	if ale.Empty() {
		return nil
	}
	return ale
}

func giveerror(w http.ResponseWriter, r *http.Request, e error, h int, where string) {
	var ale APIListErr
	ale.giveerror(w, r, e, h, where)
}

func (ale *APIListErr) giveerror(w http.ResponseWriter, r *http.Request, e error, h int, where string) {
	ale.put(e, where)
	ale.give(w, r, h)
}

func (ale APIListErr) give(w http.ResponseWriter, r *http.Request, h int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(h)
	log.Println(h, ale)
	if err := json.NewEncoder(w).Encode(ale); err != nil {
		panic(err)
	}
}

func makeoe(e error, where string) OurError {
	return OurError{
		Type:    reflect.TypeOf(e).String(),
		Content: e.Error(),
		Where:   where,
	}
}

func (ale *APIListErr) Oeappend(oel ...OurError) {
	ale.Errors = append(ale.Errors, oel...)
}

func (ale *APIListErr) add(oe OurError) {
	log.Println(oe)
	ale.Oeappend(oe)
}

func (ale *APIListErr) put(e error, where string) {
	ale.add(makeoe(e, where))
}

type hcod int

func (h *hcod) m(i int) {
	if *h == 0 {
		*h = hcod(i)
	}
}

func (mu *Multi) APISignUp(w http.ResponseWriter, r *http.Request) {
	var su SignUpPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &su); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	log.Println("signing up: ", su)
	var gi SignUpGive
	gi.User, gi.Player, gi.Auth, err = mu.Server.SignUp(su.Login, su.Passwd, su.Name)
	if err != nil {
		giveerror(w, r, err, 422, "server_signup")
		return
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
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &li); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	var aut Authorization
	aut.ID, aut.AuthKey, err = mu.Server.LogIn(li.Login, li.Passwd)
	if err != nil {
		giveerror(w, r, err, http.StatusForbidden, "server_login")
		return
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
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &bkg); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	var aut Authorization
	aut.ID, aut.AuthKey, err = mu.Server.BotKey(bkg.BotID, bkg.UserAuth.ID, bkg.UserAuth.AuthKey)
	if err != nil {
		giveerror(w, r, err, http.StatusForbidden, "server_botkey")
		return
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
	Botid    int64  `json:"botid"`
	PlayerID int64  `json:"playerid"`
	AuthKey  []byte `json:"authkey"`
}

func (mu *Multi) APINewBot(w http.ResponseWriter, r *http.Request) {
	var nbp NewBotPost
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &nbp); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	var nbg NewBotGive
	nbg.Botid, nbg.PlayerID, nbg.AuthKey, err =
		mu.Server.NewBot(nbp.WhoAmI, nbp.UserAuth.ID, nbp.UserAuth.AuthKey, nbp.OwnName, nbp.Settings)
	if err != nil {
		giveerror(w, r, err, 422, "server_newbot")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(nbg); err != nil {
		panic(err)
	}
}

type GameplayPost struct {
	State game.State `json:"state"`
	White *int64     `json:"whiteplayer"`
	Gray  *int64     `json:"grayplayer"`
	Black *int64     `json:"blackplayer"`
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
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &gpp); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	var gpg GameplayGive
	gpg.Key, err = mu.Server.AddGame(&gpp.State, [3]*int64{gpp.White, gpp.Gray, gpp.Black}, time.Now())
	if err != nil {
		giveerror(w, r, err, 422, "server_addgame")
		return
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
	if err = r.Body.Close(); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &turnp); err != nil {
		giveerror(w, r, err, 422, "unmarshal")
		return
	}
	log.Println("TURNunmarshal", turnp)
	ok, err := mu.Server.PAuth(turnp.WhoPlayer.ID, turnp.WhoPlayer.AuthKey)
	if !(ok && err == nil) {
		if !ok {
			giveerror(w, r, errors.New("Auth failed"), http.StatusForbidden, "server_pauth_notok")
		} else if err != nil {
			giveerror(w, r, err, 422, "server_pauth")
		}
		return
	}
	ourint, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	var maak MoveAndAfterKeys
	maak.MoveKey, maak.AfterGameKey, err = mu.Server.MoveGame(ourint, turnp.FromToProm, turnp.WhoPlayer.ID)
	log.Println(maak, err)
	if err != nil {
		giveerror(w, r, err, 422, "server_movegame")
		return
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
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	err = mu.Server.LoadGP(key, &gp)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadgp")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp game.State
	key, err := strconv.ParseInt(vars["stateId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	err = mu.Server.LoadState(key, &gp)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadsd")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

type VFTPGenGive struct {
	game.State  `json:"state"`
	FromToProms []game.FromToProm `json:"fromtoproms"`
}

func (mu *Multi) APIVFTPGen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp VFTPGenGive
	key, err := strconv.ParseInt(vars["stateId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
	}
	err = mu.Server.LoadState(key, &gp.State)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadsd")
		return
	}
	gp.FromToProms = make([]game.FromToProm, 0, 50)
	for ftp := range game.VFTPGen(&gp.State) {
		gp.FromToProms = append(gp.FromToProms, ftp)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

type DiffGive struct {
	server.MoveData `json:"move"`
	BeforeGame      server.GameplayData `json:"beforegame"`
	AfterGame       server.GameplayData `json:"aftergame"`
	BeforeState     game.State          `json:"beforestate"`
	AfterState      game.State          `json:"afterstate"`
	DiffBoard       []game.BoardDiff    `json:"diffboard"`
}

func (mu *Multi) APIDiff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var dig DiffGive
	key, err := strconv.ParseInt(vars["moveId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
	}
	err = mu.Server.LoadMD(key, &dig.MoveData)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadmd")
	}
	err = mu.Server.LoadGP(dig.MoveData.BeforeGame, &dig.BeforeGame)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadgp_diff_before")
	}
	err = mu.Server.LoadGP(dig.MoveData.AfterGame, &dig.AfterGame)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadgp_diff_after")
	}
	err = mu.Server.LoadState(dig.BeforeGame.State, &dig.BeforeState)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadsd_diff_before")
	}
	err = mu.Server.LoadState(dig.AfterGame.State, &dig.AfterState)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadsd_diff_after")
	}
	dig.DiffBoard = dig.BeforeState.Board.Diff(dig.AfterState.Board)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dig); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var gp server.MoveData
	key, err := strconv.ParseInt(vars["moveId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	err = mu.Server.LoadMD(key, &gp)
	if err != nil {
		giveerror(w, r, err, 421, "server_loadmd")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gp); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIAfter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sthem := [3]string{r.FormValue("white"), r.FormValue("gray"), r.FormValue("black")}
	var them [3]*int64
	for no := range sthem {
		if len(sthem[no]) > 0 {
			v, err := strconv.ParseInt(sthem[no], 10, 64)
			if err != nil {
				continue
			}
			them[no] = &v
		}
	}
	key, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	what, err := mu.Server.AfterMD(key, them)
	if err != nil {
		giveerror(w, r, err, 421, "server_aftermd")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(what); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIBefore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	what, err := mu.Server.BeforeMD(key)
	if err != nil {
		giveerror(w, r, err, 421, "server_beforemd")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(what); err != nil {
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
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	iwit.ID, iwit.IsItABot, err = mu.Server.WhoIsIt(key)
	if err != nil {
		giveerror(w, r, err, 421, "server_whoisit")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(iwit); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var iu server.InfoUser
	key, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	iu.Login, iu.Name, iu.Player, err = mu.Server.UserInfo(key)
	if err != nil {
		giveerror(w, r, err, 421, "server_userinfo")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(iu); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIBotInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sib server.InfoBot
	key, err := strconv.ParseInt(vars["botId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	sib.WhoAmI, sib.Owner, sib.OwnName, sib.Player, sib.Settings, err = mu.Server.BotInfo(key)
	if err != nil {
		giveerror(w, r, err, 421, "server_botinfo")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sib); err != nil {
		panic(err)
	}
}

func (mu *Multi) APIOwnersBots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		giveerror(w, r, err, http.StatusBadRequest, "parseint")
		return
	}
	what, err := mu.Server.OwnersBots(key)
	if err != nil {
		giveerror(w, r, err, 421, "server_ownersbots")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(what); err != nil {
		panic(err)
	}
}
