package server

import (
	"database/sql"
	"log"

	"github.com/huandu/go-sqlbuilder"
)

type SavedUrl struct {
	Url         string `json:"url"`
	Description string `json:"description"`
	Created_at  string `json:"created_at"`
	Read_at     string `json:"read_at,omitempty"`
	Priority    int    `json:"priority"`
}

var defaultFlavor = sqlbuilder.SQLite

func getUrlsQuery(pending bool) string {
	sb := defaultFlavor.NewSelectBuilder()

	sb.Select(
		"url",
		"description",
		"created_at",
		"read_at",
		"priority",
	).From("urls")

	if pending {
		sb.Where(sb.IsNull("read_at"))
	} else {
		sb.Where(sb.IsNotNull("read_at"))
	}

	sql, _ := sb.Build()
	return sql
}

func GetUrls(db *sql.DB, pending bool) ([]SavedUrl, error) {
	query := getUrlsQuery(pending)
	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]SavedUrl, 0)

	for rows.Next() {
		u := SavedUrl{}
		read_at := sql.NullString{}

		err = rows.Scan(&u.Url, &u.Description, &u.Created_at, &read_at, &u.Priority)

		if err != nil {
			log.Printf("error scanning row: %v", err)
			continue
		}

		if read_at.Valid {
			u.Read_at = read_at.String
		}

		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func markUrlAsDoneQuery(url string) (string, []interface{}) {
	sb := defaultFlavor.NewUpdateBuilder()

	sb.Update("urls")
	sb.Set("read_at = datetime('now')")
	sb.Where(sb.Equal("url", url))

	return sb.Build()
}

func MarkUrlAsDone(db *sql.DB, url string) error {
	query, args := markUrlAsDoneQuery(url)
	log.Println(query, args)
	_, err := db.Exec(query, args...)
	return err
}

func updateUrlQuery(url string, data map[string]interface{}) (string, []interface{}) {
	sb := defaultFlavor.NewUpdateBuilder()

	sb.Update("urls")
	for key, value := range data {
		sb.Set(
			sb.Assign(key, value),
		)
	}
	sb.Where(sb.Equal("url", url))

	return sb.Build()
}

func UpdateUrl(db *sql.DB, url string, data map[string]interface{}) error {
	query, args := updateUrlQuery(url, data)
	log.Println(query, args)
	_, err := db.Exec(query, args...)
	return err
}

func insertUrlQuery(data []InsertUrlRequestData) (string, []interface{}) {
	sb := defaultFlavor.NewInsertBuilder()

	sb.InsertIgnoreInto("urls")
	sb.Cols("url", "description", "created_at", "read_at", "priority")

	for _, d := range data {
		sb.Values(d.Url, d.Description, sqlbuilder.Raw("datetime('now')"), nil, 0)
	}
	return sb.Build()
}

func InsertUrls(db *sql.DB, data []InsertUrlRequestData) error {
	query, args := insertUrlQuery(data)
	log.Printf("Inserting %d urls", len(data))
	_, err := db.Exec(query, args...)
	return err
}
