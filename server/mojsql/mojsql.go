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
	res, err = conn.Exec(`
	CREATE TABLE if not exists 3manst (
		id bigint auto_increment primary key, 
		board blob not null, 
		moatzero boolean not null, 
		moatone boolean not null,
		moattwo boolean not null,
		movesnext tinyint not null, 
		castling tinyint not null, 
		enpasprevrow tinyint not null, 
		enpasprevfile tinyint not null,
		enpascurrow tinyint not null, 
		enpascurfile tinyint not null,
		halfmoveclock tinyint not null, 
		fullmovenumber smallint not null, 
		alivewhite bool not null,
		alivegray bool not null,
		aliveblack bool not null
	) ENGINE = InnoDB`)
	log.Println(res)
	if err != nil {
		return err
	}
	res, err = conn.Exec(`
	create table if not exists 3manplayer (
		id bigint auto_increment primary key,
		whoami varchar(100) not null,
		name varchar(100) not null,
		precise double,
		coefficient double,
		pawnpromotion tinyint
	) engine = InnoDB`)
	log.Println(res)
	if err != nil {
		return err
	}
	res, err := conn.Exec(`
	CREATE TABLE if not exists 3mangp (
		id bigint auto_increment primary key, 
		state bigint not null, 
		white bigint not null, 
		gray bigint not null, 
		black bigint not null, 
		date datetime not null,
		constraint
			foreign key (state) references 3manst (id)
			on update restrict,
		constraint
			foreign key (white) references 3manplayer (id)
			on update restrict,
		constraint
			foreign key (gray) references 3manplayer (id)
			on update restrict,
		constraint
			foreign key (black) references 3manplayer (id)
			on update restrict
	) ENGINE = InnoDB`)
	log.Println(res)
	return err
}

func (m *MojSQL) SaveGP(gpd *server.GameplayData) (string, error) {
}
