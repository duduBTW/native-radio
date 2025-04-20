package dbManager

import "database/sql"

var USER_ID = 1

func hasUser(db *sql.DB) bool {
	var id int
	err := db.QueryRow("SELECT id FROM USER_SETTINGS WHERE id = ?", USER_ID).Scan(&id)
	if err != nil {
		return false
	}

	return id == USER_ID
}

func createDefaultUser(db *sql.DB) error {
	_, err := db.Exec(`
	INSERT INTO USER_SETTINGS (
		volume
	)
	VALUES (?)
	`,
		50,
	)

	return err
}

func GetUserVolume(db *sql.DB) int {
	var volume int
	err := db.QueryRow("SELECT volume FROM USER_SETTINGS WHERE id = ?", USER_ID).Scan(&volume)
	if err != nil {
		return 0
	}

	return volume
}

func UpdateVolume(db *sql.DB, newVolume int) error {
	_, err := db.Exec("UPDATE USER_SETTINGS SET volume = ? WHERE id = ?", newVolume, USER_ID)
	return err
}

func GetUserSelectedIndex(db *sql.DB) int {
	var selectedIndex int
	err := db.QueryRow("SELECT selected_song_index FROM USER_SETTINGS WHERE id = ?", USER_ID).Scan(&selectedIndex)
	if err != nil {
		return 0
	}

	return selectedIndex
}

func UpdateSelectedIndex(db *sql.DB, newSelectedIndex int) error {
	_, err := db.Exec("UPDATE USER_SETTINGS SET selected_song_index = ? WHERE id = ?", newSelectedIndex, USER_ID)
	return err
}

func SetupUser(db *sql.DB) error {
	if hasUser(db) {
		return nil
	}

	return createDefaultUser(db)
}
