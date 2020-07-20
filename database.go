package main

import (
	"database/sql"
)

// Repository represent the repositories
type Repository interface {
	Close()
	Up() error
	Create(siteOGMetaData ogMetaData) error
	Update(siteOGMetaData ogMetaData) error
	FindBySite(site string) (*ogMetaData, error)
}

// repository represent the repository model
type myDatabase struct {
	db *sql.DB
}

// NewRepository will create a variable that represent the Repository struct
func NewRepository(dialect, dsn string) (Repository, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &myDatabase{db}, nil
}

// Close attaches the provider and close the connection
func (r *myDatabase) Close() {
	r.db.Close()
}

func (r *myDatabase) Up() error {

	query :=
		"CREATE TABLE IF NOT EXISTS Inspiration (" +
			"id int(11) AUTO_INCREMENT," +
			"url VARCHAR(50)," +
			"meta_og_image VARCHAR(50)," +
			"meta_og_description VARCHAR(50)," +
			"meta_og_height INTEGER," +
			"meta_og_width INTEGER," +
			"PRIMARY KEY (`id`)" +
			")"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *myDatabase) Create(siteOGMetaData ogMetaData) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO Inspiration (url,meta_og_image, meta_og_description,meta_og_width, meta_og_height) VALUES (?,?,?,?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(siteOGMetaData.site, siteOGMetaData.ogImage, siteOGMetaData.ogDescription, siteOGMetaData.ogWidth, siteOGMetaData.ogHeight)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *myDatabase) FindBySite(site string) (*ogMetaData, error) {
	meta := new(ogMetaData)

	err := r.db.QueryRow("SELECT url, meta_og_image, meta_og_description, meta_og_height, meta_og_width FROM Inspiration WHERE url = ?", site).Scan(&meta.site, &meta.ogImage, &meta.ogDescription, &meta.ogHeight, &meta.ogWidth)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (r *myDatabase) Update(siteOGMetaData ogMetaData) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtStr := "UPDATE Inspiration SET meta_og_image=?, meta_og_description=?, meta_og_height=?, meta_og_width=? WHERE url=?"
	stmt, err := tx.Prepare(stmtStr)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(siteOGMetaData.ogImage, siteOGMetaData.ogDescription, siteOGMetaData.ogHeight, siteOGMetaData.ogWidth, siteOGMetaData.site)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
