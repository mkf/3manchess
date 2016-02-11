package mojsql

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
	res, err := conn.Exec("CREATE TABLE if not exists 3mangp (id bigint auto_increment primary key, state bigint, white bigint, gray bigint, black bigint, date datetime)")
	log.Println(res)
	return err
}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (string, error) {
}
