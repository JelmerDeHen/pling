package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initSql() error {
	var err error

	if config.Dsn() == "" {
		return errors.New("dsn not configured, running without activity tracking")
	}

	db, err = Open(config.Dsn())
	if err != nil {
		log.Fatal(err)
	}

	CreateSchema()
	return nil
}

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateSchema() {
	sqlStmt := `
		create table if not exists activity (
			[id] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			[state] TEXT,
			[start] DATE,
			[stop] DATE
		);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

type ActivityRecord struct {
	Id    int
	State string
	Start time.Time
	Stop  time.Time
}

func LogActivity(record *ActivityRecord) {
	if db == nil {
		//log.Println("Running without database")
		return
	}

	if record.Stop.Sub(record.Start) < time.Second {
		log.Println("Activity duration was less than second, return")
		return
	}

	tx, err := db.Begin()
	stmt, err := tx.Prepare("insert into activity(state, start, stop) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.State, record.Start.UTC(), record.Stop.UTC())
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func GetActivityRecords() (records []*ActivityRecord) {
	if db == nil {
		log.Println("Running without database")
		return
	}

	rows, err := db.Query("select id, state, start, stop from activity")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var state string
		var start time.Time
		var stop time.Time

		err = rows.Scan(&id, &state, &start, &stop)
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, &ActivityRecord{
			Id:    id,
			State: state,
			Start: start,
			Stop:  stop,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func ListRecords() {
	records := GetActivityRecords()
	for _, record := range records {
		duration := record.Stop.Sub(record.Start)
		localizedStart := record.Start.Local().Format("20060102 15:04:05")
		localizedStop := record.Stop.Local().Format("20060102 15:04:05")
		log.Printf("state=%s; start=%s; stop=%s; duration=%s\n", record.State, localizedStart, localizedStop, duration)
	}
}
