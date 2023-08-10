package store

type UserStore interface {
	User() UserRepository
}

type UrlStore interface {
	Url() UrlRepository
}
