package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"local-rag/internal/config"
	"time"
	_ "modernc.org/sqlite"
)

type MetadataDB struct {
	db *sql.DB
}

type Source struct {
	ID       int64
	FilePath string
	CheckSum string
	Metadata map[string]any
	AddedAt  time.Time
}

func NewMetadataDB() (*MetadataDB, error) {
	db, err := sql.Open("sqlite", config.GetConfig().SqliteDb)
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	schema := `
	CREATE TABLE IF NOT EXISTS sources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT UNIQUE,
		checksum TEXT,
		metadata TEXT,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &MetadataDB{db: db}, nil
}

func (m *MetadataDB) SaveSource(ctx context.Context, source *Source) (int64, error) {
	metaJSON, _ := json.Marshal(source.Metadata)

	query := `
	INSERT INTO sources (file_path, checksum, metadata)
	VALUES (?, ?, ?)
	ON CONFLICT(file_path) DO UPDATE SET
		checksum = excluded.checksum,
		metadata = excluded.metadata
	RETURNING id;
	`

	var id int64
	err := m.db.QueryRowContext(ctx, query, source.FilePath, source.CheckSum, string(metaJSON)).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MetadataDB) Close() error {
	return m.db.Close()
}