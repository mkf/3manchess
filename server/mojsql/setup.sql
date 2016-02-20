drop table if exists 3manst;
create table 3manst (
	id bigint auto_increment primary key, 
	board binary(144) not null, 
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
	aliveblack bool not null,
	unique everything(
		board, moatzero, moatone, moattwo, movesnext, castling,
		enpasprevrow, enpasprevfile, enpascurrow, enpascurfile,
		halfmoveclock, fullmovenumber,
		alivewhite, alivegray, aliveblack
	)
) ENGINE = InnoDB;

drop table if exists 3manplayer;
create table 3manplayer (
	id bigint auto_increment primary key,
	whoami varbinary(20) not null,
	name varchar(100) not null,
	precise double,
	coefficient double,
	pawnpromotion tinyint
) engine = InnoDB;

drop table if exists 3mangp;
create table 3mangp (
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
) ENGINE = InnoDB;

drop table if exists 3manmv;
create table 3manmv (
	id bigint auto_increment primary key,
	fromrank tinyint not null,
	fromfile tinyint not null,
	torank tinyint not null,
	tofile tinyint not null,
	before bigint not null,
	promotion tinyint not null,
	constraint
		foreign key (before) references 3manst (id)
		on update restrict
) engine = InnoDB;
