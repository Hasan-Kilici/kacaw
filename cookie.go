package main

import (
	"net/http"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	Expires  time.Time
	Secure   bool
	HttpOnly bool
}

type Session struct {
	cookie *Cookie
}

func NewSession(cookie *Cookie) *Session {
	return &Session{
		cookie: cookie,
	}
}

func (s *Session) Get(key string) string {
	return s.cookie.Value
}

func (s *Session) Set(key, value string) {
	s.cookie.Value = value
}

type CookieManager interface {
	SetCookie(w http.ResponseWriter, cookie *Cookie)
	GetCookie(r *http.Request, name string) (*http.Cookie, error)
}

type DefaultCookieManager struct{}

func (cm *DefaultCookieManager) SetCookie(w http.ResponseWriter, cookie *Cookie) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
	})
}

func (cm *DefaultCookieManager) GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}
