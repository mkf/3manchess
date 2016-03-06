package mojsql

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/server"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "strconv"

import "log"

import "errors"

type MojSQL struct {
	conn *sql.DB
}

func (m *MojSQL) Initialize(username string, password string, database string) error {
	conn, err := sql.Open("mysql", username+":"+password+"@/"+database)
	m.conn = conn
	if err != nil {
		return err
	}
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
	enpassant := string([4]byte(sd.EnPassant)[:])
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
	var bmoats, bcastling, benpassant, balive []bool
	bmoats, bcastling, balive = tobool(moats), tobool(castling), tobool(balive)
	sd.Moats, sd.Castling, sd.EnPassant, sd.Alive = [3]bool(bmoats), [6]bool(bcastling), [4]int8(enpassant), [3]bool(balive)
}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (int64, error) {
	stmt, err := m.conn.Prepare("insert into 3mangp (state,created,white,gray,black) values (?,?,?,?,?)")
	if err != nil {
		return -1, err
	}
	players := make([]int64, 0, 3)
	if gpd.White != nil {
		players = append(players, *(gpd.White))
	}
	if gpd.Gray != nil {
		players = append(players, *(gpd.Gray))
	}
	if gpd.Black != nil {
		players = append(players, *(gpd.Black))
	}
	res, err := stmt.Exec(gpd.State, gpd.Date, players...)
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

func (m *MojSQL) GetAuth(playerid int64) (authkey []byte, err error) {
	stmt, err := m.conn.Prepare("select auth from 3manplayer where id=?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(playerid)
	var authkey []byte
	err = row.Scan(&authkey)
	return authkey, err
}

func (m *MojSQL) NewPlayer() (playerid int64, authkey []byte, err error) {
	res, err := m.conn.Exec("insert into 3manplayer (auth) values (md5(rand()))")
	if err != nil {
		return -1, nil, err
	}
	playerid, err := res.LastInsertId()
	if err != nil {
		return playerid, nil, err
	}
	authkey, err := m.GetAuth(playerid)
	return playerid, authkey, err
}

func (m *MojSQL) SignUp(login string, passwd string, name string) (userid int64, playerid int64, authkey []byte, err error) {
	playerid, authkey, err := m.NewPlayer()
	if err != nil {
		return -1, playerid, authkey, err
	}
	stmt, err = m.conn.Prepare("insert into chessuser (login,passwd,name,player) values (?,sha2(?,256),?,?)")
	if err != nil {
		return -1, playerid, authkey, err
	}
	res, err = stmt.Exec(login, passwd, name, playerid)
	if err != nil {
		return -1, playerid, authkey, err
	}
	userid, err := res.LastInsertId()
	return userid, playerid, authkey, err
}

func (m *MojSQL) LogIn(login string, passwd string) (userid int64, authkey []byte, err error) {
	stmt, err := m.conn.Prepare("select id,3manplayer.auth from chessuser inner join 3manplayer where login=? and passwd=sha2(?,256) and player=3manplayer.id")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(login, passwd)
	var userid int64
	var authkey []byte
	err = row.Scan(&userid, &authkey)
	return authkey, err
}

func (m *MojSQL) Auth(userid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from chessuser join 3manplayer where id=? and 3manplayer.auth=? and player=3manplayer.id)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(userid, authkey).Scan(&a)
	return a, err
}

func (m *MojSQL) BAuth(botid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from chessbot join 3manplayer where id=? and 3manplayer.auth=? and player=3manplayer.id)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(botid, authkey).Scan(&a)
	return a, err
}

func (m *MojSQL) PAuth(playerid int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare("select exists (select id from 3manplayer where id=? and auth=?)")
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(playerid, authkey).Scan(&a)
	return a, err
}

func (m *MojSQL) NewBot(whoami []byte, userid int64, uauth []byte, ownname string, settings []byte) (botid int64, playerid int64, botauth []byte, err error) {
	if ok, err := m.Auth(userid, uauth); !(ok && err == nil) {
		return -1, -1, nil, err
	}
	playerid, authkey, err := m.NewPlayer()
	if err != nil {
		return -1, playerid, botauth, err
	}
	stmt, err := m.conn.Prepare("insert into chessbot (whoami,owner,ownname,player,settings) values (?,?,?,?,?)")
	if err != nil {
		return -1, playerid, botauth, err
	}
	res, err := stmt.Exec(whoami, userid, ownname, playerid, settings)
	if err != nil {
		return -1, playerid, botauth, err
	}
	botid, err := res.LastInsertId()
	return botid, playerid, botauth, err
}

func (m *MojSQL) BotOwnerLogin(botid int64) (login string, err error) {
	stmt, err := m.conn.Prepare("select chessuser.login from chessbot inner join chessuser where owner=chessuser.id and id=?")
	if err != nil {
		return "", err
	}
	row, err := m.conn.QueryRow(botid)
	if err != nil {
		return "", err
	}
	var login string
	err = row.Scan(&login)
	return login, err
}

//WhoIsIt takes a playerid, and returns userid or bot id, then true if it is a bot or false if it's a user
func (m *MojSQL) WhoIsIt(playerid int64) (id int64, isitabot bool, err error) {
	stmt, err := m.conn.Prepare("set @playerid = ?; select id, '0' as isitabot from chessuser where player=@playerid union all select id, '1' as isitabot from chessbot where player=@playerid")
	if err != nil {
		return -1, false, err
	}
	row := stmt.QueryRow(playerid)
	var id int64
	var isitabot bool
	err = row.Scan(&id, &isitabot)
	return id, isitabot, err
}

func (m *MojSQL) BotKey(botid int64, userid int64, uauth []byte) (playerid int64, botauth []byte, err error) {
	if ok, err := m.Auth(userid, uauth); !(ok && err == nil) {
		return -1, nil, err
	}
	stmt, err := m.conn.Prepare("select p.id,p.auth from chessbot b join 3manplayer p on b.player=p.id where b.id=?")
	if err != nil {
		return -1, nil, err
	}
	row := stmt.QueryRow(botid)
	var playerid int64
	var botauth []byte
	err = row.Scan(&playerid, &botauth)
	return playerid, botauth, err
}
