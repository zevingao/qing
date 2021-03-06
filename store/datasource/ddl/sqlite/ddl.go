package sqlite

import (
	"database/sql"
)

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createMigration(db); err != nil {
		return err
	}
	if err := createTable(db); err != nil {
		return err
	}
	// migration
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {
			continue
		}
		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}
	}
	return nil
}

// 这里的设计还是比较屌的，先创建一个数据库迁移表，然后再用这个表
// 记录数据库迁移的记录
func createTable(db *sql.DB) error {
	for _, sql := range index {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

var index = []string{
	createTableUsers,
	createTableUserBinds,
	createTableDiscussions,
	createTablePosts,
	createTableLikes,
	createTableTags,
	createTableTagRel,
	createTableNotifications,
	createTableChats,
	createTableMeta,
	createTableMedia,
	createTableReport,
}

//eg:
//{
//  name: "posts.add_field_favor_count",
//  stmt: "ALTER TABLE posts ADD COLUMN favor_count INTEGER  default 0 AFTER author_id",
// },
var migrations = []struct{
	name string
	stmt string
} {

}

func createMigration(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`


const createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,
updated_at INTEGER,

sign_count INTEGER,
exp_count  INTEGER,
last_login INTEGER,

blocked_at  INTEGER,
silenced_at INTEGER,

login 	 VARCHAR(250),
nickname VARCHAR(255),
email 	 VARCHAR(255),
phone 	 VARCHAR(200),
avatar 	 VARCHAR(255),
summary  VARCHAR(500),

hash 	      VARCHAR(255),
password_hash VARCHAR(255),

UNIQUE(login)
);`

const createTableUserBinds = `
CREATE TABLE IF NOT EXISTS user_binds (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

user_id 	INTEGER,
bind_id 	TEXT,
union_id    TEXT,
kind 		TEXT,

UNIQUE(kind, bind_id)
);`


const createTableDiscussions = `
CREATE TABLE IF NOT EXISTS discussions (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,
updated_at INTEGER,

title 		VARCHAR(255),
content 	TEXT,
author_id 	INTEGER,

first_post    INTEGER,
last_post     INTEGER,
comment_count INTEGER
);`

const createTablePosts = `
CREATE TABLE IF NOT EXISTS posts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

discussion_id INTEGER,
author_id 	  INTEGER,
reply_id      INTEGER,
parent_id 	  INTEGER,
like_count    INTEGER,

content 	  TEXT
);`

const createTableLikes = `
CREATE TABLE IF NOT EXISTS likes (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

status 		INTEGER,
author_id 	INTEGER,
post_id 	INTEGER,

UNIQUE(author_id, post_id)
);`

const createTableTags = `
CREATE TABLE IF NOT EXISTS tags (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

_order      INTEGER,
parent_id   INTEGER,
color       INTEGER,
text 		VARCHAR(200),

summary 	TEXT,

UNIQUE(text)
);`

const createTableTagRel = `
CREATE TABLE IF NOT EXISTS tag_discussions (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

discussion_id INTEGER,
tag_id INTEGER,

UNIQUE(discussion_id, tag_id)
);`

const createTableNotifications = `
CREATE TABLE IF NOT EXISTS notifications (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

entity_id  INTEGER,
entity_ty  INTEGER,
from_id    INTEGER,
to_id      INTEGER,
status     INTEGER
);`

const createTableChats = `
CREATE TABLE IF NOT EXISTS chats (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

content    TEXT,
_type      INTEGER,

chat_id    BIGINT,
status     INTEGER,
from_id    INTEGER,
to_id      INTEGER
);`


const createTableMeta = `
CREATE TABLE IF NOT EXISTS metadata (
id INTEGER PRIMARY KEY AUTOINCREMENT,

kv_key 	  TEXT,
kv_value  TEXT,

UNIQUE(kv_key)
);`

const createTableSignIn = `
CREATE TABLE IF NOT EXISTS sign_logs (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at INTEGER,

seq_count   INTEGER,
sum_count   INTEGER,
user_id 	INTEGER,
date_day    TEXT,

UNIQUE(user_id, date_day)
);`

const createTableMedia = `
CREATE TABLE IF NOT EXISTS medias (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at  INTEGER,

post_id 	INTEGER,
author_id   INTEGER,
_type 		INTEGER,

path	TEXT,
meta	TEXT
);`

const createTableReport = `
CREATE TABLE IF NOT EXISTS reports (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at  INTEGER,
updated_at  INTEGER,

entity_id   INTEGER,
entity_ty   INTEGER,
counter     INTEGER,
status      INTEGER,
user_id     INTEGER,
report_ty   INTEGER,

content     string,
other       string,
images      string
);`