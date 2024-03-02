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

type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	LatestTen() ([]*Snippet, error)
}

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {

	queryString := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY ))`
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
	query := "SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"
	rows, err := model.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []*Snippet
	for rows.Next() {
		var snippet = &Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
