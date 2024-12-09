package model

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Usuario struct {
	Usuario_id         int64
	Usuario_nombre     string
	Usuario_psswd      string
	Usuario_privilegio int64
	Usuario_activo     bool
}

type UsrAuth struct {
	Usuario    string
	Passwd     string
	Activo     int
	Privilegio string
}

func AllUsers(db *sql.DB) ([]Usuario, error) {
	query, err := DB.Query(`SELECT * FROM usuarios`)
	if err != nil {
		return nil, fmt.Errorf("AllUsers SelectUsuarios %v", err)
	}
	var table []Usuario
	defer query.Close()
	for query.Next() {
		var usr Usuario
		if err := query.Scan(&usr.Usuario_id, &usr.Usuario_nombre, &usr.Usuario_psswd, &usr.Usuario_privilegio, &usr.Usuario_activo); err != nil {
			return nil, fmt.Errorf("AllUsers ScanError: %v", err)
		}
		table = append(table, usr)
	}
	return table, nil
}

func AddNewUser(c echo.Context, usr string, pwd string, priv string) error {
	verifyQuery := DB.QueryRow(`SELECT usuario_nombre FROM usuario WHERE usuario_nombre = ?`, usr)
	var dupeUser string
	if verifyError := verifyQuery.Scan(&dupeUser); verifyError == nil {
		c.Response().Header().Add("HX-Trigger", "duplicateError")
		return fmt.Errorf("AddNewUser verify: %s already exists", dupeUser)
	}

	query := DB.QueryRow(`SELECT privilegio_id FROM privilegio WHERE privilegio_nombre = ?`, priv)
	var privInt int
	if queryError := query.Scan(&privInt); queryError != nil {
		c.Response().Header().Add("HX-Trigger", "invalidPrivilegeError")
		return fmt.Errorf("AddNewUser query: %v", queryError)
	}

	_, err := DB.Exec(`INSERT INTO usuario (usuario_nombre, usuario_psswd, usuario_activo, usuario_privilegio)
	VALUES (?,?,?,?)`, usr, pwd, 1, privInt)
	if err != nil {
		c.Response().Header().Add("HX-Trigger", "insertError")
		return fmt.Errorf("AddNewUser: %v", err)
	}

	return nil
}

func getIDFromUsername(username string) (int, error) {
	response := DB.QueryRow("SELECT usuario_id FROM usuario WHERE usuario_nombre = ?", username)
	var id int
	if err := response.Scan(&id); err != nil {
		return -1, fmt.Errorf("getIDFromUsername: %v", err)
	}
	return id, nil
}
