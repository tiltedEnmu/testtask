package postgres

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

type Storage struct {
    db *sql.DB
}

func New(host, port, user, password, dbname string) (*Storage, error) {
    const op = "storage.postgres.New"

    psqlInfo := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname,
    )

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println(err)
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    err = db.Ping()
    if err != nil {
        fmt.Println(err)
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{
        db: db,
    }, nil
}
