package model

import (
	"database/sql"
	"env"
)

type Author struct {
	Author_id   int    `author_id`
	Author_name string `author_name`
}

func GetAuthorByName(name string) (*Author, error) {
	var stmt *sql.Stmt

	if stmt, err := env.Db.Prepare(
		`SELECT author_id
         FROM novel_author_info
         WHERE author_name=?`); err != nil {
		env.Log.Warn("[SQL][SELECT][author:%s][error:%s]", name, err.Error())
		return nil, err
	}
	defer stmt.Close()

	author := &Author{
		Author_name: name,
	}

	if err := stmt.QueryRow(name).
		Scan(&author.Author_id); err != nil {
		env.Log.Warn("[SQL][SELECT][author:%s][error:%s]", name, err.Error())
		return nil, err
	}

	return author, nil
}
