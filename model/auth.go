package model

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Authenticate(username, password string, c echo.Context) (bool, error) {
	var realData UsrAuth
	if strings.ContainsAny(username, "\"`'") || strings.ContainsAny(password, "\"`'") {
		return false, fmt.Errorf("Auth invalidChar")
	}
	rows := DB.QueryRow(`SELECT u.usuario_nombre, u.usuario_psswd, u.usuario_activo, p.privilegio_nombre FROM usuario u
	JOIN privilegio p on p.privilegio_id = u.usuario_privilegio
WHERE u.usuario_nombre = ?`, username)
	if err := rows.Scan(&realData.Usuario, &realData.Passwd, &realData.Activo, &realData.Privilegio); err != nil {
		return false, fmt.Errorf("auth/scan: %v", err)
	}
	if realData.Activo == 0 {
		return false, fmt.Errorf("auth/disableduser")
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("auth/rowserr: %v", err)
	}
	if subtle.ConstantTimeCompare([]byte(username), []byte(realData.Usuario)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(realData.Passwd)) == 1 {
		c.SetCookie(&http.Cookie{
			Name:     "usrtype",
			Value:    realData.Privilegio, //admin, user, manager validos.
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})
		c.SetCookie(&http.Cookie{
			Name:     "login",
			Value:    "yes",
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})
		c.SetCookie(&http.Cookie{
			Name:     "usrname",
			Value:    realData.Usuario,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})
		return true, nil
	}
	return false, nil
}
