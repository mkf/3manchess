drop table if exists 3manst;
create table 3manst (
	id bigint auto_increment primary key, 
	board binary(144) not null, 
	moats bit(3) not null,
	movesnext tinyint not null, 
	castling bit(6) not null, 
	enpassant binary(4) not null,
	halfmoveclock tinyint not null, 
	fullmovenumber smallint not null, 
	alive bit(3) not null,
	unique key everything(
		board, moats, movesnext, castling,
		enpassant,
		halfmoveclock, fullmovenumber,
		alive
	)
) ENGINE = InnoDB;

drop table if exists 3manplayer;
create table 3manplayer (
	id bigint auto_increment primary key,
	whoami varbinary(20) not null, -- player type identifier
	name varchar(100) not null, -- name, not parsed
	precise double, -- bot depth
	coefficient double, -- bot coefficient
	pawnpromotion tinyint -- bot promotion if not auto
) engine = InnoDB;

drop table if exists 3mangp;
create table 3mangp (
	id bigint auto_increment primary key, 
	state bigint not null, 
	white bigint not null, 
	gray bigint not null, 
	black bigint not null, 
	created datetime not null
--	constraint
--		foreign key (state) references 3manst (id)
--		on update restrict,
--	constraint
--		foreign key (white) references 3manplayer (id)
--		on update restrict,
--	constraint
--		foreign key (gray) references 3manplayer (id)
--		on update restrict,
--	constraint
--		foreign key (black) references 3manplayer (id)
--		on update restrict
) ENGINE = InnoDB;

drop table if exists 3manmv;
create table 3manmv (
	id bigint auto_increment primary key,
	fromto binary(4) not null,
	beforestate bigint not null,
	promotion tinyint not null
--	constraint
--		foreign key (before) references 3manst (id)
--		on update restrict
) engine = InnoDB;

-- vi:ft=mysql