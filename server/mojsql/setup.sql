drop table if exists 3manst;
create table 3manst (
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
) ENGINE = InnoDB;

create trigger dbo.BlockDuplicates3manst
	on dbo.3manst
	instead of insert
as
begin
	set nocount on;
	if not exists (select 1 from inserted as i
		inner join dbo.3manst as t
		on i.board = t.board
		and i.moatzero = t.moatzero
		and i.moatone = t.moatone
		and i.moattwo = t.moattwo
		and i.movesnext = t.movesnext
		and i.castling = t.castling
		and i.enpasprevrow = t.enpasprevrow
		and i.enpasprevfile = t.enpasprevfile
		and i.enpascurrow = t.enpascurrow
		and i.enpascurfile = t.enpascurfile
		and i.halfmoveclock = t.halfmoveclock
		and i.fullmovenumber = t.fullmovenumber
		and i.alivewhite = t.alivewhite
		and i.alivegray = t.alivegray
		and i.aliveblack = t.aliveblack
	)
	begin
		insert dbo.3manst(board,
			moatzero, moatone, moattwo, 
			movesnext, castling,
			enpasprevrow, enpasprevfile, enpascurrow, enpascurfile,
			halfmoveclock, fullmovenumber,
			alivewhite, alivegray, aliveblack
		)
			select board,
				moatzero, moatone, moattwo, 
				movesnext, castling,
				enpasprevrow, enpasprevfile, enpascurrow, enpascurfile,
				halfmoveclock, fullmovenumber,
				alivewhite, alivegray, aliveblack
				from inserted;
	end
	else
	begin
		select board,
			moatzero, moatone, moattwo, 
			movesnext, castling,
			enpasprevrow, enpasprevfile, enpascurrow, enpascurfile,
			halfmoveclock, fullmovenumber,
			alivewhite, alivegray, aliveblack
			from inserted;
	end
end
go


drop table if exists 3manplayer;
create table 3manplayer (
	id bigint auto_increment primary key,
	whoami varchar(100) not null,
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
