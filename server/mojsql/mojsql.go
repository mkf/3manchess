package mojsql

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/server"
import "time"
import "encoding/hex"
import "database/sql"
import _ "github.com/go-sql-driver/mysql" //importing sql driver is idiomatically done using a blank import

import "log"

//MojSQL is an instance of sql binding
type MojSQL struct {
	conn *sql.DB
}

//Interface returns server.Server(m)
func (m *MojSQL) Interface() server.Server { return server.Server{m} }

//Initialize initializes the sql binding
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

//SaveSD inserts StateData into db
func (m *MojSQL) SaveSD(sd *game.StateData) (key int64, err error) {
	key = -1
	moats := bitint(sd.Moats[:])       // string(tobit(sd.Moats[:]))
	castling := bitint(sd.Castling[:]) // string(tobit(sd.Castling[:]))
	eenp := fourbyte(sd.EnPassant)
	alive := bitint(sd.Alive[:]) // string(tobit(sd.Alive[:]))
	ultbo := hex.EncodeToString(sd.Board[:])
	ultenp := hex.EncodeToString(eenp[:])
	log.Println("vals:",
		ultbo,
		moats,
		sd.MovesNext,
		castling,
		ultenp,
		sd.HalfmoveClock,
		sd.FullmoveNumber,
		alive)
	resstmt, err := m.conn.Prepare(
		`insert into 3manst (
			board,
			moats,
			movesnext,
			castling,
			enpassant,
			halfmoveclock,
			fullmovenumber,
			alive
		) values (
			unhex(?),  -- bo  hex
			?,         -- mt
			?,         -- mn
			?,         -- ct
			unhex(?),  -- ep  hex
			?,         -- hm
			?,         -- fm
			?          -- al
		) on duplicate key update id=last_insert_id(id)`)
	defer resstmt.Close()
	log.Println(resstmt, err)
	if err != nil {
		return
	}
	res, err := resstmt.Exec(
		ultbo,
		moats,
		sd.MovesNext,
		castling,
		ultenp,
		sd.HalfmoveClock,
		sd.FullmoveNumber,
		alive)
	log.Println(res, err)
	if err != nil {
		return
	}
	key, err = res.LastInsertId()
	log.Println(key, err)
	return
}

//LoadSD gets StateData from db
func (m *MojSQL) LoadSD(key int64, sd *game.StateData) error {
	givestmt, err := m.conn.Prepare(
		`select
			board,
			moats+0,
			movesnext,
			castling+0,
			enpassant,
			halfmoveclock,
			fullmovenumber,
			alive+0
		from 3manst where id=?`)
	defer givestmt.Close()
	if err != nil {
		return err
	}
	give := givestmt.QueryRow(key)
	var moats, castling, alive uint8
	var board, enpassant []byte
	err = give.Scan(&board, &moats, &sd.MovesNext, &castling, &enpassant, &sd.HalfmoveClock, &sd.FullmoveNumber, &alive)
	if err != nil {
		return err
	}
	var bmoats, bcastling, balive []bool
	bmoats, bcastling, balive = intbit(moats, 3), intbit(castling, 6), intbit(alive, 3)
	sd.Moats, sd.Castling, sd.EnPassant, sd.Alive = bas3(bmoats), bas6(bcastling), fourint8(yas4(enpassant)), bas3(balive)
	sd.Board = game.Byte144(board)
	return err
}

//SaveGP inserts GameplayData into db
func (m *MojSQL) SaveGP(gpd *server.GameplayData) (int64, error) {
	stmt, err := m.conn.Prepare(
		`insert into 3mangp (state,created,white,gray,black) values (?,?,?,?,?)`)
	//		on duplicate key update id=last_insert_id(id)`)
	defer stmt.Close()
	if err != nil {
		return -1, err
	}
	players := make([]interface{}, 0, 3)
	players = append(players, gpd.State, gpd.Date)
	if gpd.White != nil {
		players = append(players, *(gpd.White))
	} else {
		players = append(players, nil)
	}
	if gpd.Gray != nil {
		players = append(players, *(gpd.Gray))
	} else {
		players = append(players, nil)
	}
	if gpd.Black != nil {
		players = append(players, *(gpd.Black))
	} else {
		players = append(players, nil)
	}
	res, err := stmt.Exec(players...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

//LoadGP gets GameplayData from db
func (m *MojSQL) LoadGP(key int64, gpd *server.GameplayData) error {
	stmt, err := m.conn.Prepare("select state,white,gray,black,created from 3mangp where id=?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	var w, g, b sql.NullInt64
	var dat string
	err = stmt.QueryRow(key).Scan(&gpd.State, &w, &g, &b, &dat)
	nullint64(&gpd.White, w)
	nullint64(&gpd.Gray, g)
	nullint64(&gpd.Black, b)
	var er error
	gpd.Date, er = time.Parse("2006-01-02 15:04:05", dat)
	if err == nil {
		return er
	}
	return err
}

//ListGP selects {number} newest Gameplays in db
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
	defer stmt.Close()
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
		h = append(h, server.GameplayFollow{Key: id, GameplayData: gd})
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

//SaveMD inserts MoveData into db
func (m *MojSQL) SaveMD(md *server.MoveData) (key int64, err error) {
	key = -1
	trans, err := m.conn.Begin()
	if err != nil {
		return
	}
	stmt, err := trans.Prepare(
		`insert into 3manmv (fromto,beforegame,promotion,who) 
		values (?,?,?,?) on duplicate key update id=last_insert_id(id)`)
	defer stmt.Close()
	if err != nil {
		log.Println(stmt)
		return
	}
	fb := fourbyte(md.FromTo)
	defer trans.Rollback()
	res, err := stmt.Exec(fb[:], md.BeforeGame, md.PawnPromotion, md.Who)
	if err != nil {
		log.Println(res)
		return
	}
	stmt, err = trans.Prepare("update 3manmv set aftergame=? where id=last_insert_id() and aftergame is null limit 1")
	if err != nil {
		log.Println(stmt)
		return
	}
	resp, err := stmt.Exec(md.AfterGame)
	if err != nil {
		log.Println(resp)
		return
	}
	lidd, er := res.LastInsertId()
	if er != nil {
		return lidd, er
	}
	return lidd, trans.Commit()
}

//LoadMD gets MoveData from db
func (m *MojSQL) LoadMD(key int64, md *server.MoveData) (err error) {
	stmt, err := m.conn.Prepare("select fromto,beforegame,aftergame,promotion,who from 3manmv where id=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	var ft []byte
	err = stmt.QueryRow(key).Scan(&ft, &md.BeforeGame, &md.AfterGame, &md.PawnPromotion, &md.Who)
	md.FromTo = fourint8(yas4(ft))
	return
}

//AfterMD lists moves after the selected gameplay
func (m *MojSQL) AfterMDwe(beforegp int64) (out []server.MoveFollow, err error) {
	stmt, err := m.conn.Prepare("select id,fromto,aftergame,promotion,who from 3manmv where beforegame=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	rows, err := stmt.Query(beforegp)
	return procafter(beforegp, rows, err)
}

func procafter(beforegp int64, rows *sql.Rows, er error) (out []server.MoveFollow, err error) {
	err = er
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

func (m *MojSQL) BeforeMD(aftergp int64) (out []server.MoveFollow, err error) {
	stmt, err := m.conn.Prepare("select id,fromto,beforegame,promotion,who from 3manmv where aftergame=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	rows, err := stmt.Query(aftergp)
	if err != nil {
		return
	}
	out = make([]server.MoveFollow, 0, 1)
	for rows.Next() {
		var neww server.MoveFollow
		var ft []byte
		neww.MoveData.AfterGame = aftergp
		err = rows.Scan(&neww.Key, &ft, &neww.MoveData.BeforeGame, &neww.MoveData.PawnPromotion, &neww.MoveData.Who)
		neww.MoveData.FromTo = fourint8(yas4(ft))
		out = append(out, neww)
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

func (m *MojSQL) OwnersBots(owner int64) (out []server.BotFollow, err error) {
	stmt, err := m.conn.Prepare("select id,whoami,ownname,player,settings from chessbot where owner=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	rows, err := stmt.Query(owner)
	if err != nil {
		return
	}
	out = make([]server.BotFollow, 0, 3)
	for rows.Next() {
		var neww server.BotFollow
		neww.Owner = owner
		err = rows.Scan(&neww.Key, &neww.WhoAmI, &neww.OwnName, &neww.Player, &neww.Settings)
		out = append(out, neww)
		if err != nil {
			return
		}
	}
	err = rows.Close()
	return
}

var minusoneint64 int64 = -1

//AfterMDwPlayers takes the before GameplayID and W,G,B players and returns the after gameplays with the same players
func (m *MojSQL) AfterMDwPlayers(beforegp int64, players [3]*int64) (out []server.MoveFollow, err error) {
	stmt, err := m.conn.Prepare(
		`select m.id,m.fromto,m.aftergame,m.promotion,m.who from 3manmv m 
		join 3mangp g on g.id=m.aftergame 
		where m.beforegame=? and 
			(g.white=? and "prawda"=?) and 
			(g.gray=? and "prawda"=?) and
			(g.black=? and "prawda"=?)
	`)
	defer stmt.Close()
	if err != nil {
		return
	}
	tstr := [3]string{"prawda", "prawda", "prawda"}
	for numer := range players {
		if players[numer] == nil {
			tstr[numer] = "nienie"
			players[numer] = &minusoneint64
		}
	}
	rows, err := stmt.Query(
		beforegp,
		*players[0], tstr[0],
		*players[1], tstr[1],
		*players[2], tstr[2],
	)
	return procafter(beforegp, rows, err)
}

//GetAuth : PLAYER(ID) → PLAYER(AUTH)
func (m *MojSQL) GetAuth(playerid int64) (authkey []byte, err error) {
	stmt, err := m.conn.Prepare("select auth from 3manplayer where id=?")
	defer stmt.Close()
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
	defer stmt.Close()
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
	stmt, err := m.conn.Prepare("select u.id,p.auth from chessuser u inner join 3manplayer p where u.login=? and u.passwd=sha2(?,256) and u.player=p.id")
	defer stmt.Close()
	if err != nil {
		return
	}
	row := stmt.QueryRow(login, passwd)
	err = row.Scan(&userid, &authkey)
	return
}

//Auth authenticates by UserID
func (m *MojSQL) Auth(userid int64, authkey []byte) (bool, error) {
	return m.procauth(
		`select exists (
			select u.id from chessuser u 
			join 
			3manplayer p 
			where u.id=? and p.auth=? and u.player=p.id
		)`,
		userid, authkey)
}

//BAuth authenticates by BotID
func (m *MojSQL) BAuth(botid int64, authkey []byte) (bool, error) {
	return m.procauth(
		`select exists (
			select b.id from chessbot b 
			join 
			3manplayer p 
			where b.id=? and p.auth=? and b.player=p.id
		)`,
		botid, authkey)
}

//PAuth authenticates by PlayerID
func (m *MojSQL) PAuth(playerid int64, authkey []byte) (bool, error) {
	return m.procauth("select exists (select id from 3manplayer where id=? and auth=?)",
		playerid, authkey)
}

func (m *MojSQL) procauth(query string, id int64, authkey []byte) (bool, error) {
	stmt, err := m.conn.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return false, err
	}
	var a bool
	err = stmt.QueryRow(id, authkey).Scan(&a)
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
	defer stmt.Close()
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
	stmt, err := m.conn.Prepare("select u.login,u.name from chessbot b inner join chessuser u where b.owner=u.id and b.id=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	err = stmt.QueryRow(botid).Scan(&login, &name)
	return
}

//UserInfo gets info about a user from db
func (m *MojSQL) UserInfo(userid int64) (login string, name string, playerid int64, err error) {
	stmt, err := m.conn.Prepare("select login,name,player from chessuser where id=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	err = stmt.QueryRow(userid).Scan(&login, &name, &playerid)
	return
}

//BotInfo gets info about a bot from db
func (m *MojSQL) BotInfo(botid int64) (whoami []byte, owner int64, ownname string, player int64, settings []byte, err error) {
	stmt, err := m.conn.Prepare("select whoami,owner,ownname,player,settings from chessbot where id=?")
	defer stmt.Close()
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
	defer stmt.Close()
	if err != nil {
		return -1, false, err
	}
	row := stmt.QueryRow(playerid, playerid)
	err = row.Scan(&id, &isitabot)
	return
}

//BotKey : BOTID+USER(ID+AUTH) → BOT(PLAYERID+AUTH)
func (m *MojSQL) BotKey(botid int64, userid int64, uauth []byte) (playerid int64, botauth []byte, err error) {
	ok, err := m.Auth(userid, uauth)
	if !(ok && err == nil) {
		return -1, nil, err
	}
	stmt, err := m.conn.Prepare("select p.id,p.auth from chessbot b join 3manplayer p on b.player=p.id where b.id=?")
	defer stmt.Close()
	if err != nil {
		return -1, nil, err
	}
	row := stmt.QueryRow(botid)
	err = row.Scan(&playerid, &botauth)
	return playerid, botauth, err
}
