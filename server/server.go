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
	AfterMD(beforegp int64) ([]AfterMoveFollow, error)
	GetAuth(playerid int64) (authkey []byte, err error)
	NewPlayer() (playerid int64, authkey []byte, err error)
	SignUp(login string, passwd string, name string) (userid int64, playerid int64, authkey []byte, err error)
	LogIn(login string, passwd string) (userid int64, authkey []byte, err error)
	Auth(userid int64, authkey []byte) (bool, error)
	BAuth(botid int64, authkey []byte) (bool, error)
	Pauth(playerid int64, authkey []byte) (bool, error)
	NewBot(whoami []byte, userid int64, uauth []byte, ownname string, settings []byte) (botid int64, playerid int64, botauth []byte)
	WhoIsIt(playerid int64) (id int64, isitabot bool, err error)
	BotKey(botid int64, userid int64, uauth []byte) (playerid int64, botauth []byte, err error)
}

type GameplayData struct {
	State              int64
	White, Gray, Black *int64
	Date               time.Time
}

type MoveData struct {
	FromTo        [4]int8
	BeforeGame    int64
	AfterGame     int64
	PawnPromotion int8
	Who           int64
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

func SaveState(sr Server, st *game.State) (key int64, err error) {
	return sr.SaveSD(st.Data())
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
