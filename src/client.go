package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ensureClientID(c echo.Context) (string, error) {
	cookie, err := c.Cookie("client_id")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	clientID, err := newClientID()
	if err != nil {
		return "", err
	}
	c.SetCookie(&http.Cookie{
		Name:     "client_id",
		Value:    clientID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60,
	})
	return clientID, nil
}

func newClientID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
