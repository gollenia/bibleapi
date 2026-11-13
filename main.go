package main

import (
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/ncruces/go-sqlite3/driver"
)

func main() {
	initStorage()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	app.Get("/books", func(c *fiber.Ctx) error {
		translationParam := getTranslationParam(c)
		booksMap := bookCache[translationParam]
		var booksList []Book
		for _, b := range booksMap {
			booksList = append(booksList, b)
		}
		return c.JSON(fiber.Map{"translation": translationParam, "books": booksList})
	})

	app.Get("/schema", func(c *fiber.Ctx) error {
		translationParam := getTranslationParam(c)
		rawSchema := schemaCache[translationParam]
		booksMap := bookCache[translationParam]

		type SortedBook struct {
			ID     string      `json:"id"`
			Name   string      `json:"name"`
			Order  int         `json:"order"`
			Verses map[int]int `json:"verses"`
		}

		var output []SortedBook

		for shortName, chaptersMap := range rawSchema {
			longName := shortName
			order := 999

			if b, exists := booksMap[shortName]; exists {
				longName = b.LongName
				order = b.Order
			}

			output = append(output, SortedBook{
				ID:     shortName,
				Name:   longName,
				Order:  order,
				Verses: chaptersMap,
			})
		}

		sort.Slice(output, func(i, j int) bool {
			return output[i].Order < output[j].Order
		})

		return c.JSON(fiber.Map{
			"translation": translationParam,
			"schema":      output,
		})
	})

	app.Get("/verses/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		translationParam := getTranslationParam(c)
		db := dbConnections[translationParam]

		req := parseBibleUri(path)
		if req.BookStr == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid URL format."})
		}

		book, bookExists := bookCache[translationParam][req.BookStr]
		if !bookExists {
			return c.Status(404).JSON(fiber.Map{"error": "Book not found: " + req.BookStr})
		}

		var rows *sql.Rows
		var err error

		if !req.IsRange {
			if req.StartVerse == 0 {
				rows, err = db.Query("SELECT book_id, chapter, verse, text, linebreak FROM verses WHERE book_id = ? AND chapter = ? ORDER BY verse", book.ID, req.StartChap)
			} else {
				rows, err = db.Query("SELECT book_id, chapter, verse, text, linebreak FROM verses WHERE book_id = ? AND chapter = ? AND verse = ? ORDER BY verse", book.ID, req.StartChap, req.StartVerse)
			}
		} else {
			if req.StartChap == req.EndChap {
				rows, err = db.Query("SELECT book_id, chapter, verse, text, linebreak FROM verses WHERE book_id = ? AND chapter = ? AND verse BETWEEN ? AND ? ORDER BY verse", book.ID, req.StartChap, req.StartVerse, req.EndVerse)
			} else {
				query := `SELECT book_id, chapter, verse, text, linebreak FROM verses 
                          WHERE book_id = ? AND (
                              (chapter = ? AND verse >= ?) OR 
                              (chapter = ? AND verse <= ?) OR 
                              (chapter > ? AND chapter < ?)
                          ) ORDER BY chapter, verse`
				rows, err = db.Query(query, book.ID, req.StartChap, req.StartVerse, req.EndChap, req.EndVerse, req.StartChap, req.EndChap)
			}
		}

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var results []VerseResponse
		for rows.Next() {
			var v VerseResponse
			if err := rows.Scan(&v.BookID, &v.Chapter, &v.Verse, &v.Text, &v.Linebreak); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			results = append(results, v)
		}

		return c.JSON(fiber.Map{"translation": translationParam, "book": book.LongName, "verses": results})
	})

	log.Fatal(app.Listen(":8080"))
}

func getTranslationParam(c *fiber.Ctx) string {
	param := strings.ToLower(c.Query("translation", ""))
	if _, exists := dbConnections[param]; exists {
		return param
	}
	for k := range dbConnections {
		return k
	}
	return ""
}
