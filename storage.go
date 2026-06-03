package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var dbConnections map[string]*sql.DB
var bookCache map[string]map[string]Book
var schemaCache map[string]map[string]map[int]int

func initStorage() {
	dbConnections = make(map[string]*sql.DB)
	bookCache = make(map[string]map[string]Book)
	schemaCache = make(map[string]map[string]map[int]int)
	
	dbFolder := "./db"

	files, err := os.ReadDir(dbFolder)
	if err != nil {
		log.Fatalf("Error reading DB folder '%s': %v", dbFolder, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" {
				translationName := strings.ToLower(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
				fullPath := filepath.Join(dbFolder, file.Name())

				db, err := sql.Open("sqlite3", fullPath)
				if err != nil {
					log.Printf("WARNING: Could not open %s: %v", fullPath, err)
					continue
				}
				dbConnections[translationName] = db
				log.Printf("Translation loaded: %s", translationName)

				bookCache[translationName] = make(map[string]Book)
				bookRows, err := db.Query("SELECT id, short_name, long_name, section, [order], chapters, testament, number FROM books")
				if err != nil {
					log.Printf("WARNING: Could not load books for %s: %v", translationName, err)
					continue
				}
				for bookRows.Next() {
					var b Book
					if err := bookRows.Scan(&b.ID, &b.ShortName, &b.LongName, &b.Section, &b.Order, &b.Chapters, &b.Testament, &b.Number); err == nil {
						bookCache[translationName][strings.ToLower(b.ShortName)] = b
					}
				}
				bookRows.Close()

				schemaCache[translationName] = make(map[string]map[int]int)
				schemaRows, err := db.Query("SELECT book_id, chapter, COUNT(verse) FROM verses GROUP BY book_id, chapter")
				if err != nil {
					log.Printf("WARNING: Could not calculate schema for %s: %v", translationName, err)
					continue
				}

				rowCount := 0
				for schemaRows.Next() {
					var bookID, chapter, totalVerses int
					if err := schemaRows.Scan(&bookID, &chapter, &totalVerses); err != nil {
						continue
					}
					rowCount++

					shortName := strconv.Itoa(bookID)
					if transBooks, ok := bookCache[translationName]; ok {
						for _, b := range transBooks {
							if b.ID == bookID {
								shortName = b.ShortName
								break
							}
						}
					}

					if schemaCache[translationName][shortName] == nil {
						schemaCache[translationName][shortName] = make(map[int]int)
					}
					schemaCache[translationName][shortName][chapter] = totalVerses
				}
				schemaRows.Close()
				log.Printf("-> Schema for %s cached in RAM! (%d chapter entries processed)", translationName, rowCount)
			}
		}
	}

	if len(dbConnections) == 0 {
		log.Fatal("No valid SQLite databases found!")
	}
}