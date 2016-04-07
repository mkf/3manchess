package server

//import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "time"
import "log"

type Server struct {
	Database
}

type Database interface {
	Initialize(username string, password string, database string) error
	SaveGP(*GameplayData) (key int64, err error)
	LoadGP(key int64, gp *GameplayData) error
	SaveSD(sd *game.StateData) (key int64, err error)
	LoadSD(key int64, sd *game.StateData) error
	SaveMD(*MoveData) (key int64, err error)
	LoadMD(key int64, md *MoveData) error
	ListGP(uint) ([]GameplayFollow, error)
	AfterMDwe(beforegp int64) ([]MoveFollow, error)
	AfterMDwPlayers(beforegp int64, players [3]*int64) ([]MoveFollow, error)
	BeforeMD(aftergp int64) ([]MoveFollow, error)
	OwnersBots(owner int64) ([]BotFollow, error)
	GetAuth(playerid int64) (authkey []byte, err error)
	NewPlayer() (playerid int64, authkey []byte, err error)
	SignUp(login string, passwd string, name string) (userid int64, playerid int64, authkey []byte, err error)
	LogIn(login string, passwd string) (userid int64, authkey []byte, err error)
	Auth(userid int64, authkey []byte) (bool, error)
	BAuth(botid int64, authkey []byte) (bool, error)
	PAuth(playerid int64, authkey []byte) (bool, error)
	NewBot(whoami []byte, userid int64, uauth []byte, ownname string, settings []byte) (botid int64, playerid int64, botauth []byte, err error)
	WhoIsIt(playerid int64) (id int64, isitabot bool, err error)
	BotKey(botid int64, userid int64, uauth []byte) (playerid int64, botauth []byte, err error)
	UserInfo(userid int64) (login string, name string, playerid int64, err error)
	BotInfo(botid int64) (whoami []byte, owner int64, ownname string, player int64, settings []byte, err error)
}

func (sr Server) SaveState(sta *game.State) (key int64, err error) {
	return sr.SaveSD(sta.Data())
}

func (sr Server) LoadState(key int64, sta *game.State) (err error) {
	da := new(game.StateData)
	if err = sr.LoadSD(key, da); err == nil {
		sta.FromData(da)
	}
	return
}

type GameplayData struct {
	State int64     `json:"stateid"`
	White *int64    `json:"whiteplayer"`
	Gray  *int64    `json:"grayplayer"`
	Black *int64    `json:"blackplayer"`
	Date  time.Time `json:"when"`
}

type MoveData struct {
	FromTo        [4]int8 `json:"fromto"`
	BeforeGame    int64   `json:"beforegp"`
	AfterGame     int64   `json:"aftergp"`
	PawnPromotion int8    `json:"pawnpromotion"`
	Who           int64   `json:"playerid"`
}

func (md MoveData) fromTo() game.FromTo {
	return game.FromTo{
		game.Pos{md.FromTo[0], md.FromTo[1]},
		game.Pos{md.FromTo[2], md.FromTo[3]},
	}
}

func (md MoveData) FromToProm() game.FromToProm {
	return game.FromToProm{
		FromTo:        md.fromTo(),
		PawnPromotion: game.FigType(md.PawnPromotion),
	}
}

func (sr Server) AfterMD(beforegp int64, filterplayers [3]*int64) (out []MoveFollow, err error) {
	for i := range filterplayers {
		if filterplayers[i] != nil {
			return sr.AfterMDwPlayers(beforegp, filterplayers)
		}
	}
	return sr.AfterMDwe(beforegp)
}

func (sr Server) AddGame(st *game.State, players [3]*int64, when time.Time) (key int64, err error) {
	sk, err := sr.SaveState(st)
	if err != nil {
		return
	}
	gpd := GameplayData{sk, players[0], players[1], players[2], when}
	key, err = sr.SaveGP(&gpd)
	return
}

func (sr Server) MoveGame(before int64, ftp game.FromToProm, who int64) (mkey int64, aftkey int64, err error) {
	log.Println("MoveGame", sr, before, ftp, who)
	var befga GameplayData
	err = sr.LoadGP(before, &befga)
	log.Println(befga, err)
	if err != nil {
		return
	}
	befga.Date = time.Now()
	var sta game.State
	err = sr.LoadState(befga.State, &sta)
	log.Println(sta, err)
	if err != nil {
		return
	}
	switch sta.MovesNext {
	case game.White:
		befga.White = &who
	case game.Gray:
		befga.Gray = &who
	case game.Black:
		befga.Black = &who
	}
	mov := ftp.Move(&sta)
	log.Println("now MoveIt", mov)
	afts, err := mov.EvalAfter()
	if err != nil {
		return
	}
	log.Println("MoveIT returned", afts)
	aftskey, err := sr.SaveState(afts)
	if err != nil {
		return
	}
	befga.State = aftskey
	aftkey, err = sr.SaveGP(&befga)
	if err != nil {
		return
	}
	mdfin := MoveData{
		FromTo:        [4]int8{ftp.FromTo.From()[0], ftp.FromTo.From()[1], ftp.FromTo.To()[0], ftp.FromTo.To()[1]},
		BeforeGame:    before,
		AfterGame:     aftkey,
		PawnPromotion: int8(ftp.PawnPromotion),
		Who:           who}
	mkey, err = sr.SaveMD(&mdfin)
	if err != nil {
		return
	}
	var loadedmd MoveData
	err = sr.LoadMD(mkey, &loadedmd)
	if err == nil {
		aftskey = loadedmd.AfterGame
	}
	return
}

func nullminusone(q *int64) int64 {
	if q == nil {
		return -1
	}
	return *q
}

type GameplayFollow struct {
	Key          int64 `json:"id"`
	GameplayData `json:"game"`
}

type MoveFollow struct {
	Key      int64 `json:"id"`
	MoveData `json:"move"`
}

type AfterMoveFollow struct {
	MoveFollow  `json:"movefollow"`
	SamePlayers bool `json:"sameplayers"`
	//YourMoveNext bool `json:"yourmovenext"`
}

type StateDataFollow struct {
	Key int64
	*game.StateData
}

type StateFollow struct {
	Key         int64 `json:"key"`
	*game.State `json:"state"`
}

type InfoUser struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Player int64  `json:"playerid"`
}

type UserFollow struct {
	Key      int64 `json:"key"`
	InfoUser `json:"userinfo"`
}

type InfoBot struct {
	WhoAmI   []byte `json:"whoami"`
	Owner    int64  `json:"ownerid"`
	OwnName  string `json:"ownname"`
	Player   int64  `json:"playerid"`
	Settings []byte `json:"settings"`
}

type BotFollow struct {
	Key     int64 `json:"key"`
	InfoBot `json:"botinfo"`
}
