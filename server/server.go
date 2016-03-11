package server

//import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "time"

type Server interface {
	Initialize(username string, password string, database string) error
	SaveGP(*GameplayData) (key int64, err error)
	LoadGP(key int64, gp *GameplayData) error
	SaveSD(sd *game.StateData) (key int64, err error)
	LoadSD(key int64, sd *game.StateData) error
	SaveMD(*MoveData) (key int64, err error)
	LoadMD(key int64, md *MoveData) error
	ListGP(uint) ([]GameplayFollow, error)
	AfterMD(beforegp int64) ([]MoveFollow, error)
	AfterMDwPlayers(beforegp int64, players [3]int64) ([]MoveFollow, error)
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

func (md MoveData) Move(sr Server) game.Move {
	s := new(game.State)
	LoadState(sr, md.BeforeGame, s)
	return game.Move{
		From:          game.Pos{md.FromTo[0], md.FromTo[1]},
		To:            game.Pos{md.FromTo[2], md.FromTo[3]},
		Before:        s,
		PawnPromotion: game.FigType(md.PawnPromotion),
	}
}

func AddGame(sr Server, st *game.StateData, players [3]*int64, when time.Time) (key int64, err error) {
	sk, err := sr.SaveSD(st)
	if err != nil {
		return
	}
	gpd := GameplayData{sk, players[0], players[1], players[2], when}
	key, err = sr.SaveGP(&gpd)
	return
}

func MoveGame(sr Server, before int64, ftp game.FromToProm, who int64) (mkey int64, aftkey int64, err error) {
	var befga GameplayData
	err = sr.LoadGP(before, &befga)
	if err != nil {
		return
	}
	befga.Date = time.Now()
	var sta game.State
	err = LoadState(sr, befga.State, &sta)
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
	afts := MoveIt(&mov, [3]int64{nullminusone(befga.White), nullminusone(befga.Gray), nullminusone(befga.Black)})
	aftskey, err := sr.SaveSD(afts.Data())
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
	return
}

func nullminusone(q *int64) int64 {
	if q == nil {
		return -1
	}
	return *q
}

func LoadState(sr Server, key int64, s *game.State) error {
	var sd game.StateData
	err := sr.LoadSD(key, &sd)
	s.FromData(&sd)
	return err
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

type StateFollow struct {
	Key int64
	*game.StateData
}
