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

func (m *MojSQL) SaveSD(*game.StateData) (key int64, err error) {

}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (string, error) {
}
