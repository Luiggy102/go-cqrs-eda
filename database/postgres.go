package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/Luiggy102/go-cqrs-eda/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	Db *sql.DB
}

func NewPostgresRepo(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{Db: db}, nil
}
func (repo *PostgresRepository) Close() {
	repo.Db.Close()
}
func (repo *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.Db.ExecContext(ctx, "insert into feeds (id, title, description) values ($1, $2, $3)",
		feed.ID, feed.Title, feed.Description)
	if err != nil {
		return err
	}
	return nil
}
func (repo *PostgresRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repo.Db.QueryContext(ctx, "select id, title, description, created_at from feeds")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	feeds := []*models.Feed{}
	for rows.Next() {
		f := &models.Feed{}
		err := rows.Scan(&f.ID, &f.Title, &f.Description, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}
