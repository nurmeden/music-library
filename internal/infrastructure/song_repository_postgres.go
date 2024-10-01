package postgres

import (
	"database/sql"
	"fmt"
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/repository"
)

type PostgresSongRepository struct {
	DB *sql.DB
}

func NewPostgresSongRepository(db *sql.DB) repository.SongRepository {
	return &PostgresSongRepository{DB: db}
}

func (r *PostgresSongRepository) FetchAll(filters map[string]interface{}, limit, offset int) ([]entity.Song, error) {
	var songs []entity.Song
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE 1=1"

	args := []interface{}{}
	if groupName, ok := filters["group_name"]; ok {
		query += " AND group_name = ?"
		args = append(args, groupName)
	}
	if songName, ok := filters["song_name"]; ok {
		query += " AND song_name = ?"
		args = append(args, songName)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song entity.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *PostgresSongRepository) FetchByID(id int) (*entity.Song, error) {
	var song entity.Song
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = ?"
	err := r.DB.QueryRow(query, id).Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song with id %d not found", id)
		}
		return nil, err
	}

	return &song, nil
}

func (r *PostgresSongRepository) Store(song *entity.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) 
              VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSongRepository) Update(song *entity.Song) error {
	query := `UPDATE songs SET group_name = ?, song_name = ?, release_date = ?, text = ?, link = ? 
              WHERE id = ?`
	_, err := r.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSongRepository) Delete(id int) error {
	query := "DELETE FROM songs WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
