drop table if exists 3manmv;
drop table if exists 3mangp;
drop table if exists chessbot;
drop table if exists chessuser;
drop table if exists 3manplayer;

drop table if exists 3manst;
create table 3manst (
	id bigint auto_increment primary key, 
	board binary(144) not null, 
	moats tinyint not null, -- _ _ _ _ _ W G B
	movesnext tinyint not null, 
	castling tinyint not null, -- _ _ K Q x u k q , where x,q are gray king&queen
	enpassant binary(4) not null,
	halfmoveclock tinyint not null, 
	fullmovenumber smallint not null, 
	alive tinyint not null,  -- _ _ _ _ _ W G B
	constraint everything unique (
		board, moats, movesnext, castling,
		enpassant,
		halfmoveclock, fullmovenumber,
		alive
	)
) ENGINE = InnoDB;

create table 3manplayer (
	id bigint auto_increment primary key,
	auth varbinary(100) not null
	-- name varchar(100) not null,
) engine = InnoDB default charset=utf8;

create table chessuser (
	id bigint auto_increment primary key,
	login varchar(20) unique key,
	passwd varchar(100) not null,
	name varchar(100),
	player bigint not null unique key,
	constraint
		foreign key (player) references 3manplayer (id)
		on update restrict
) engine = InnoDB default charset=utf8;

create table chessbot (
	id bigint auto_increment primary key,
	whoami varbinary(20) not null, -- ai type identifier
	owner bigint not null,
	ownname varchar(50),
	player bigint not null unique key,
	settings varbinary(500),
	constraint everything unique ( whoami, owner, settings ),
	constraint
		foreign key (owner) references chessuser (id)
		on update restrict,
	constraint
		foreign key (player) references 3manplayer (id)
		on update restrict
) engine = InnoDB default charset=utf8;

create table 3mangp (
	id bigint auto_increment primary key, 
	state bigint,
	white bigint, 
	gray bigint, 
	black bigint, 
	created timestamp default current_timestamp,
--	constraint everything unique ( state,white,gray,black ),
	constraint
		foreign key (white) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (gray) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (black) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (state) references 3manst (id)
		on update restrict
) ENGINE = InnoDB;

create table 3manmv (
	id bigint auto_increment primary key,
	fromto binary(4) not null,
	beforegame bigint not null,
	aftergame bigint,
	promotion tinyint not null,
	who bigint not null,
	constraint
		foreign key (beforegame) references 3mangp (id)
		on update restrict,
	constraint
		foreign key (who) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (aftergame) references 3mangp (id)
		on update cascade,
	constraint
		foreign key (aftergame) references 3mangp (id)
		on delete restrict,
	constraint onemove unique (fromto, beforegame, promotion, who)
) engine = InnoDB;


-- vi:ft=mysql
