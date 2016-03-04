package mojsql

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/server"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "strconv"

import "log"

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

func (m *MojSQL) SaveSD(sd *game.StateData) (key int64, err error) {
	board := string(sd.Board[:])
	moats := string(tobit(sd.Moats[:]))
	castling := string(tobit(sd.Castling[:]))
	enpassant := string([4]byte(sd.EnPassant)[:])
	alive := string(tobit(sd.Alive[:]))
	whether, err := m.conn.Query(`select id from 3manst where board="` + board +
		`" and moats="` + moats + `" and movesnext=` + strconv.Itoa(sd.MovesNext) +
		` and castling="` + castling + `" and enpassant="` + enpassant + `" and halfmoveclock=` + strconv.Itoa(sd.HalfmoveClock) +
		` and fullmovenumber=` + strconv.Itoa(sd.FullmoveNumber) + ` and alive="` + alive)
	if err != nil {
		return -1, err
	}
	if whether.Next() {
		nasz := int64(-1)
		err := whether.Scan(&nasz)
		return nasz, err
	}
	res, err := m.conn.Exec(`insert into 3manst (board,moats,movesnext,castling,enpassant,halfmoveclock,fullmovenumber,alive) values ("` +
		board + `","` + moats + `","` + strconv.Itoa(sd.MovesNext) + `","` + castling + `","` + enpassant + `","` + strconv.Itoa(sd.HalfmoveClock) + "," + strconv.Itoa(sd.FullmoveNumber) + `,"` + alive + `")`)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (string, error) {
}
