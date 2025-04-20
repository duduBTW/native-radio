package dbManager

import (
	"database/sql"

	"github.com/dudubtw/osu-radio-native/lib"
)

func InsertSong(db *sql.DB, song *lib.Song) error {
	_, err := db.Exec(`
	INSERT INTO SONGS (
		audio, bg, osu_file, path, date_added,
		title, artist, creator, duration,
		beatmap_set_id, mode, title_unicode, artist_unicode,
		primary_color, secondary_color
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		song.Audio,
		song.Bg,
		song.OsuFile,
		song.Path,
		song.DateAdded,
		song.Title,
		song.Artist,
		song.Creator,
		song.Duration,
		song.BeatmapSetID,
		song.Mode,
		song.TitleUnicode,
		song.ArtistUnicode,
		song.PrimaryColor,
		song.SecondaryColor,
	)

	return err
}

func SelectAllSongs(db *sql.DB, searchString string, targetID int) ([]lib.Song, int, error) {
	query := `
		SELECT
			id, audio, bg, osu_file, path, date_added,
			title, artist, creator, duration,
			beatmap_set_id, mode, title_unicode, artist_unicode,
			primary_color, secondary_color
		FROM SONGS
	`

	var args []interface{}

	if searchString != "" {
		query += `
			WHERE title LIKE $1 OR artist LIKE $1
		`
		args = append(args, "%"+searchString+"%")
	}

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()

	var songs []lib.Song

	containsID := -1
	index := 0
	for rows.Next() {
		var song lib.Song
		err := rows.Scan(
			&song.ID,
			&song.Audio,
			&song.Bg,
			&song.OsuFile,
			&song.Path,
			&song.DateAdded,
			&song.Title,
			&song.Artist,
			&song.Creator,
			&song.Duration,
			&song.BeatmapSetID,
			&song.Mode,
			&song.TitleUnicode,
			&song.ArtistUnicode,
			&song.PrimaryColor,
			&song.SecondaryColor,
		)
		if err != nil {
			return nil, containsID, err
		}
		songs = append(songs, song)

		if song.ID == targetID {
			containsID = index
		}
		index++
	}

	return songs, containsID, rows.Err()
}

func NewSongTableFromDb(db *sql.DB) (*lib.SongTable, error) {
	songs, _, err := SelectAllSongs(db, "", 0)
	if err != nil {
		return nil, err
	}

	return &lib.SongTable{
		Songs: songs,
	}, nil
}
