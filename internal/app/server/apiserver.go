package server

import (
	"database/sql"
	"github.com/GlebSr/nuhaiShort/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	userStore := sqlstore.NewUserStore(db)
	urlStore := sqlstore.NewUrlStore(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(userStore, urlStore, sessionStore) //умоиооалдтомфулирцлуватирыуатцли затычка
	return http.ListenAndServe(config.BindAddr, cors.Default().Handler(srv))
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
