package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/logger"
	"github.com/nurmeden/music-library/internal/repository"
	"strings"
)

type PostgresSongRepository struct {
	DB     *sql.DB
	logger logger.Logger
}

func NewPostgresSongRepository(db *sql.DB, logger logger.Logger) repository.SongRepository {
	return &PostgresSongRepository{DB: db, logger: logger}
}

func (r *PostgresSongRepository) FetchAll(filters map[string]interface{}, limit, offset int) ([]entity.Song, error) {
	r.logger.Infof("Fetching all songs with filters: %v, limit: %d, offset: %d", filters, limit, offset)

	var songs []entity.Song
	var query strings.Builder
	query.WriteString(`SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE is_deleted = FALSE and 1=1`)

	args := []interface{}{}
	if groupName, ok := filters["group_name"]; ok {
		query.WriteString(` AND group_name = ?`)
		args = append(args, groupName)
	}
	if songName, ok := filters["song_name"]; ok {
		query.WriteString(` AND song_name = ?`)
		args = append(args, songName)
	}

	query.WriteString(` LIMIT ? OFFSET ?`)
	args = append(args, limit, offset)

	r.logger.Debugf("Executing query: %s with args: %v", query.String(), args)

	rows, err := r.DB.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song entity.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			r.logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error iterating rows: %v", err)
		return nil, err
	}

	r.logger.Infof("Successfully fetched %d songs", len(songs))

	return songs, nil
}

func (r *PostgresSongRepository) FetchByID(id int) (*entity.Song, error) {
	r.logger.Infof("Fetching song by ID: %d", id)

	var song entity.Song
	var query strings.Builder
	query.WriteString(`SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = ?`)
	err := r.DB.QueryRow(query.String(), id).Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warnf("Song with ID %d not found", id)
			return nil, fmt.Errorf("song with id %d not found", id)
		}
		r.logger.Errorf("Error fetching song by ID: %v", err)
		return nil, err
	}

	r.logger.Infof("Successfully fetched song: %s", song.SongName)
	return &song, nil
}

func (r *PostgresSongRepository) Store(song *entity.Song) error {
	r.logger.Infof("Storing new song: %s by group: %s", song.SongName, song.GroupName)

	var query strings.Builder
	query.WriteString(`INSERT INTO songs (group_name, song_name, release_date, text, link) 
              VALUES (?, ?, ?, ?, ?)`)
	_, err := r.DB.Exec(query.String(), song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		r.logger.Errorf("Error storing song %s: %v", song.SongName, err)
		return err
	}

	r.logger.Infof("Successfully stored song: %s", song.SongName)
	return nil
}

func (r *PostgresSongRepository) Update(song *entity.Song) error {
	r.logger.Infof("Updating song with ID: %d", song.ID)

	var query strings.Builder
	query.WriteString(`UPDATE songs SET group_name = ?, song_name = ?, release_date = ?, text = ?, link = ? 
              WHERE id = ?`)
	_, err := r.DB.Exec(query.String(), song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		r.logger.Errorf("Error updating song with ID %d: %v", song.ID, err)
		return err
	}

	r.logger.Infof("Successfully updated song with ID: %d", song.ID)
	return nil
}

func (r *PostgresSongRepository) Delete(id int) error {
	r.logger.Infof("Marking song with ID: %d as deleted", id)

	var query strings.Builder
	query.WriteString(`UPDATE songs SET is_deleted = TRUE, updated_at = NOW() WHERE id = ?`)
	_, err := r.DB.Exec(query.String(), id)
	if err != nil {
		r.logger.Errorf("Error marking song with ID %d as deleted: %v", id, err)
		return err
	}

	r.logger.Infof("Successfully marked song with ID: %d as deleted", id)
	return nil
}
