package mojsql

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/server"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

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

func makebit(b bool) byte {
	if b {
		return '1'
	} else {
		return '0'
	}
}

func tobit(b []bool) []byte {
	a := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		a = append(a, makebit(b[i]))
	}
	return a
}

func makebool(b byte) bool {
	switch b {
	case '1':
		return true
	case '0':
		return false
	default:
		panic(b)
	}
}

func tobool(b []byte) []bool {
	a := make([]bool, 0, len(b))
	for i := 0; i < len(b); i++ {
		a = append(a, makebool(b[i]))
	}
	return a
}

func (m *MojSQL) SaveSD(*game.StateData) (key int64, err error) {

}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (string, error) {
}
