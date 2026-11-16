package repo

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

const schema = `
create table if not exists users (
	id serial primary key,
	login text unique not null,
	password text not null,
	role text not null default 'user'
);
create table if not exists categories (
	id serial primary key,
	name text not null
);
create table if not exists products (
	id serial primary key,
	name text not null,
	description text not null default '',
	price int not null,
	image text not null default '',
	category_id int references categories(id) on delete set null
);
create table if not exists cart (
	user_id int references users(id) on delete cascade,
	product_id int references products(id) on delete cascade,
	qty int not null default 1,
	primary key (user_id, product_id)
);`

type Repo struct {
	Db *sql.DB
}

func New() *Repo {
	db, err := sql.Open("pgx", "postgres://postgres:пароль@localhost:5432/shop?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	if _, err = db.Exec(schema); err != nil {
		panic(err)
	}

	var n int
	err = db.QueryRow("select count(*) from users where login = 'admin'").Scan(&n)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("insert into users (login, password, role) values ('admin', $1, 'admin')", string(hash))
		if err != nil {
			panic(err)
		}
	}
	return &Repo{Db: db}
}
