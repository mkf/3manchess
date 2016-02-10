package mojsql

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "log"

type MojSQL struct {
	conn *sql.DB
}

func (m *MojSQL) Initialize(username string, password string, database string) error {
	conn, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		return err
	}
	m.conn = conn
	res, err := conn.Exec("CREATE TABLE 3manchess if not exists;")
	log.Println(res)
	if err != nil {
		return err
	}
	return nil
}
