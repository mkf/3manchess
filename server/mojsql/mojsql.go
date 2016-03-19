package mojsql

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/server"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "log"

type MojSQL struct {
	conn *sql.DB
}

func (m *MojSQL) Interface() server.Server { return m }

func (m *MojSQL) Initialize(username string, password string, database string) error {
	conn, err := sql.Open("mysql", username+":"+password+"@unix(/var/run/mysql/mysql.sock)/"+database)
	//go func() { defer conn.Close() }()
	m.conn = conn
	if err != nil {
		return err
	}
	err = m.conn.Ping()
	return err
}

func (m *MojSQL) TransactionStart() error {
	return nil
}

func (m *MojSQL) TransactionEnd() error {
	return nil
}

func (m *MojSQL) SaveSD(sd *game.StateData) (key int64, err error) {
	board := string(sd.Board[:])
	moats := string(tobit(sd.Moats[:]))
	castling := string(tobit(sd.Castling[:]))
	eenp := fourbyte(sd.EnPassant)
	enpassant := string(eenp[:])
	alive := string(tobit(sd.Alive[:]))
	whetherstmt, err := m.conn.Prepare("select id from 3manst where board=? and moats=? and movesnext=? and castling=? and enpassant=? and halfmoveclock=? and fullmovenumber=? and alive=?")
	log.Println(whetherstmt, err)
	if err != nil {
		return -1, err
	}
	whether, err := whetherstmt.Query(board, moats, sd.MovesNext, castling, enpassant, sd.HalfmoveClock, sd.FullmoveNumber, alive)
	log.Println(whether, err)
	if err != nil {
		return -1, err
	}
	if whether.Next() {
		nasz := int64(-1)
		err := whether.Scan(&nasz)
		log.Println(nasz, err)
		return nasz, err
	}
	resstmt, err := m.conn.Prepare("insert into 3manst (board,moats,movesnext,castling,enpassant,halfmoveclock,fullmovenumber,alive) values (?,?,?,?,?,?,?,?)")
	log.Println(resstmt, err)
	if err != nil {
		return -1, err
	}
	res, err := resstmt.Exec(board, moats, sd.MovesNext, castling, enpassant, sd.HalfmoveClock, sd.FullmoveNumber, alive)
	log.Println(res, err)
	if err != nil {
		return -1, err
	}
	var lid int64
	lid, err = res.LastInsertId()
	log.Println(lid, err)
	return lid, err
}

func (m *MojSQL) LoadSD(key int64, sd *game.StateData) error {
	givestmt, err := m.conn.Prepare("select board,moats,movesnext,castling,enpassant,halfmoveclock,fullmovenumber,alive from 3manst where id=?")
	if err != nil {
		return err
	}
	give := givestmt.QueryRow(key)
	var board, moats, castling, enpassant, alive []byte
	err = give.Scan(&board, &moats, &sd.MovesNext, &castling, &enpassant, &sd.HalfmoveClock, &sd.FullmoveNumber, &alive)
	if err != nil {
		return err
	}
	var bmoats, bcastling, balive []bool
	bmoats, bcastling, balive = tobool(moats), tobool(castling), tobool(alive)
	sd.Moats, sd.Castling, sd.EnPassant, sd.Alive = bas3(bmoats), bas6(bcastling), fourint8(yas4(enpassant)), bas3(balive)
	return err
}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (int64, error) {
	stmt, err := m.conn.Prepare("insert into 3mangp (state,created,white,gray,black) values (?,?,?,?,?)")
	if err != nil {
		return -1, err
	}
	players := make([]interface{}, 0, 3)
	players = append(players, gpd.State, gpd.Date)
	if gpd.White != nil {
		players = append(players, *(gpd.White))
	}
	if gpd.Gray != nil {
		players = append(players, *(gpd.Gray))
	}
	if gpd.Black != nil {
		players = append(players, *(gpd.Black))
	}
	res, err := stmt.Exec(players...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (m *MojSQL) LoadGP(key int64, gpd *server.GameplayData) error {
	stmt, err := m.conn.Prepare("select state,white,gray,black,created from 3mangp where id=?")
	if err != nil {
		return err
	}
	var w, g, b sql.NullInt64
	err = stmt.QueryRow(key).Scan(&gpd.State, &w, &g, &b, &gpd.Date)
	nullint64(&gpd.White, w)
	nullint64(&gpd.Gray, g)
	nullint64(&gpd.Black, b)
	return err
}

func (m *MojSQL) ListGP(many uint) (h []server.GameplayFollow, err error) {
	qstr := "select id,state,white,gray,black,created from 3mangp order by created desc"
	l := many != 0
	ar := make([]interface{}, 0, 1)
	ile := uint(500)
	if l {
		qstr += " limit ?"
		ar = append(ar, many)
		ile = many
	}
	stmt, err := m.conn.Prepare(qstr)
	if err != nil {
		return
	}
	rows, err := stmt.Query(ar...)
	if err != nil {
		return
	}
	h = make([]server.GameplayFollow, 0, ile)
	for rows.Next() {
		var id int64
		var gd server.GameplayData
		var w, g, b sql.NullInt64
		err = rows.Scan(&id, &gd.State, &w, &g, &b, &gd.Date)
		nullint64(&gd.White, w)
		nullint64(&gd.Gray, g)
		nullint64(&gd.Black, b)
		h = append(h, server.GameplayFollow{id, gd})
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

func (m *MojSQL) SaveMD(md *server.MoveData) (key int64, err error) {
	stmt, err := m.conn.Prepare("insert into 3manmv (fromto,beforegame,aftergame,promotion,who) values (?,?,?,?,?)")
	key = -1
	if err != nil {
		return
	}
	fb := fourbyte(md.FromTo)
	res, err := stmt.Exec(fb[:], md.BeforeGame, md.AfterGame, md.PawnPromotion, md.Who)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func (m *MojSQL) LoadMD(key int64, md *server.MoveData) (err error) {
	stmt, err := m.conn.Prepare("select fromto,beforegame,aftergame,promotion,who from 3manmv where id=?")
	if err != nil {
		return
	}
	var ft []byte
	err = stmt.QueryRow(key).Scan(&ft, &md.BeforeGame, &md.AfterGame, &md.PawnPromotion, &md.Who)
	md.FromTo = fourint8(yas4(ft))
	return
}

func (m *MojSQL) AfterMD(beforegp int64) (out []server.MoveFollow, err error) {
	stmt, err := m.conn.Prepare("select id,fromto,aftergame,promotion,who from 3manmv where beforegame=?")
	if err != nil {
		return
	}
	rows, err := stmt.Query(beforegp)
	if err != nil {
		return
	}
	out = make([]server.MoveFollow, 0, 1)
	for rows.Next() {
		var neww server.MoveFollow
		var ft []byte
		neww.MoveData.BeforeGame = beforegp
		err = rows.Scan(&neww.Key, &ft, &neww.MoveData.AfterGame, &neww.MoveData.PawnPromotion, &neww.MoveData.Who)
		neww.MoveData.FromTo = fourint8(yas4(ft))
		out = append(out, neww)
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

//AfterMDwPlayers takes the before GameplayID and W,G,B players and returns the after gameplays with the same players
func (m *MojSQL) AfterMDwPlayers(beforegp int64, players [3]int64) (out []server.MoveFollow, err error) {
	stmt, err := m.conn.Prepare("select id,fromto,aftergame,promotion,who from 3manmv join 3mangp g on g.id=aftergame where beforegame=? and n.white=? and n.gray=? and n.black=?")
	if err != nil {
		return
	}
	rows, err := stmt.Query(beforegp, players[0], players[1], players[2])
	if err != nil {
		return
	}
	out = make([]server.MoveFollow, 0, 1)
	for rows.Next() {
		var neww server.MoveFollow
		var ft []byte
		neww.MoveData.BeforeGame = beforegp
		err = rows.Scan(&neww.Key, &ft, &neww.MoveData.AfterGame, &neww.MoveData.PawnPromotion, &neww.MoveData.Who)
		neww.MoveData.FromTo = fourint8(yas4(ft))
		out = append(out, neww)
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

//GetAuth : PLAYER(ID) → PLAYER(AUTH)
func (m *MojSQL) GetAuth(playerid int64) (authkey []byte, err error) {
	stmt, err := m.conn.Prepare("select auth from 3manplayer where id=?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(playerid)
	err = row.Scan(&authkey)
	return authkey, err
}

//NewPlayer creates a new player and gives PLAYER(ID + AUTH)
func (m *MojSQL) NewPlayer() (playerid int64, authkey []byte, err error) {
	res, err := m.conn.Exec("insert into 3manplayer (auth) values (md5(rand()))")
	if err != nil {
		return -1, nil, err
	}
	playerid, err = res.LastInsertId()
	if err != nil {
		return playerid, nil, err
	}
	authkey, err = m.GetAuth(playerid)
	return playerid, authkey, err
}

//SignUp : NEW_USER(LOGIN + PASSWD + NAME) → USER(ID + PLAYERID + AUTH)
func (m *MojSQL) SignUp(login string, passwd string, name string) (userid int64, playerid int64, authkey []byte, err error) {
	playerid, authkey, err = m.NewPlayer()
	if err != nil {
		return -1, playerid, authkey, err
	}
	stmt, err := m.conn.Prepare("insert into chessuser (login,passwd,name,player) values (?,sha2(?,256),?,?)")
	if err != nil {
		return -1, playerid, authkey, err
	}
	res, err := stmt.Exec(login, passwd, name, playerid)
	if err != nil {
		return -1, playerid, authkey, err
	}
	userid, err = res.LastInsertId()
	return userid, playerid, authkey, err
}

//LogIn : USER(LOGIN + PASSWD) → USER(ID + AUTH)
func (m *MojSQL) LogIn(login string, passwd string) (userid int64, authkey []byte, err error) {
	stmt, err := m.conn.Prepare("select id,3manplayer.auth from chessuser inner join 3manplayer where login=? and passwd=sha2(?,256) and player=3manplayer.id")
	if err != nil {
		return
	}
	row := stmt.QueryRow(login, passwd)
	err = row.Scan(&userid, &authkey)
	return
}

//Auth authenticates by UserID
func (m *MojSQL) Auth(userid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from chessuser join 3manplayer where id=? and 3manplayer.auth=? and player=3manplayer.id)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(userid, authkey).Scan(&a)
	return a, err
}

//BAuth authenticates by BotID
func (m *MojSQL) BAuth(botid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from chessbot join 3manplayer where id=? and 3manplayer.auth=? and player=3manplayer.id)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(botid, authkey).Scan(&a)
	return a, err
}

//PAuth authenticates by PlayerID
func (m *MojSQL) PAuth(playerid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from 3manplayer where id=? and auth=?)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(playerid, authkey).Scan(&a)
	return a, err
}

//NewBot : NEW_BOT(WHOAMI+OWNNAME+SETTINGS) + USER(ID+AUTH) → BOT(ID + PLAYERID + AUTH)
func (m *MojSQL) NewBot(whoami []byte, userid int64, uauth []byte, ownname string, settings []byte) (botid int64, playerid int64, botauth []byte, err error) {
	botid, playerid = -1, -1
	ok, err := m.Auth(userid, uauth)
	if !(ok && err == nil) {
		return
	}
	playerid, botauth, err = m.NewPlayer()
	if err != nil {
		return
	}
	stmt, err := m.conn.Prepare("insert into chessbot (whoami,owner,ownname,player,settings) values (?,?,?,?,?)")
	if err != nil {
		return
	}
	res, err := stmt.Exec(whoami, userid, ownname, playerid, settings)
	if err != nil {
		return
	}
	botid, err = res.LastInsertId()
	return
}

//BotOwnerLoginAndName : BOTID → OWNER(LOGIN+NAME)   //TO BE DEPRECATED
func (m *MojSQL) BotOwnerLoginAndName(botid int64) (login string, name string, err error) {
	stmt, err := m.conn.Prepare("select chessuser.login,chessuser.name from chessbot inner join chessuser where owner=chessuser.id and id=?")
	if err != nil {
		return
	}
	err = stmt.QueryRow(botid).Scan(&login, &name)
	return
}

func (m *MojSQL) UserInfo(userid int64) (login string, name string, playerid int64, err error) {
	stmt, err := m.conn.Prepare("select login,name,player from chessuser where id=?")
	if err != nil {
		return
	}
	err = stmt.QueryRow(userid).Scan(&login, &name, &playerid)
	return
}

func (m *MojSQL) BotInfo(botid int64) (whoami []byte, owner int64, ownname string, player int64, settings []byte, err error) {
	stmt, err := m.conn.Prepare("select whoami,owner,ownname,player,settings from chessbot where id=?")
	if err != nil {
		return
	}
	err = stmt.QueryRow(botid).Scan(&whoami, &owner, &ownname, &player, &settings)
	return
}

//WhoIsIt : PLAYERID → (BOT/USER)(ID) + ?[BOT/USER]¿
//WhoIsIt takes a playerid, and returns userid or bot id, then true if it is a bot or false if it's a user
func (m *MojSQL) WhoIsIt(playerid int64) (id int64, isitabot bool, err error) {
	stmt, err := m.conn.Prepare(
		`select id, '0' as isitabot from chessuser where player=? 
		union all 
		select id, '1' as isitabot from chessbot where player=?`)
	if err != nil {
		return -1, false, err
	}
	row := stmt.QueryRow(playerid, playerid)
	err = row.Scan(&id, &isitabot)
	return
}

//BotKey : BOTID+USER(ID+AUTH) → BOT(PLAYERID+AUTH)
func (m *MojSQL) BotKey(botid int64, userid int64, uauth []byte) (playerid int64, botauth []byte, err error) {
	if ok, err := m.Auth(userid, uauth); !(ok && err == nil) {
		return -1, nil, err
	}
	stmt, err := m.conn.Prepare("select p.id,p.auth from chessbot b join 3manplayer p on b.player=p.id where b.id=?")
	if err != nil {
		return -1, nil, err
	}
	row := stmt.QueryRow(botid)
	err = row.Scan(&playerid, &botauth)
	return playerid, botauth, err
}
