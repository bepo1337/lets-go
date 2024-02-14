package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {

	queryString := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY ))`
	//	queryString := `INSERT INTO snippets (title, content, created, expires)
	//VALUES (?, ?, UTC_TIMESTAMP(), ?)`
	//	dateAddString := internal.DateAddString(expires, internal.DAY)
	result, err := model.DB.Exec(queryString, title, content, expires)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), err
}

func (model *SnippetModel) Get(id int) (*Snippet, error) {
	var snippet = &Snippet{}
	query := "SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
	err := model.DB.QueryRow(query, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snippet, nil
}

func (model *SnippetModel) LatestTen() ([]*Snippet, error) {
	return nil, nil
}
