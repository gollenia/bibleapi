package main

type Book struct {
	ID        int    `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
	Section   string `json:"section"`
	Order     int    `json:"order"`
	Chapters  int    `json:"chapters"`
	Testament string `json:"testament"`
	Number    int    `json:"number"`
}

type VerseResponse struct {
	BookID    int    `json:"book_id"`
	Chapter   int    `json:"chapter"`
	Verse     int    `json:"verse"`
	Text      string `json:"text"`
	Linebreak int    `json:"linebreak"`
}

type BibleRequest struct {
	BookStr     string
	StartChap   int
	StartVerse  int
	EndChap     int
	EndVerse    int
	IsRange     bool
}

type ChapterSchema struct {
	Chapter     int `json:"chapter"`
	TotalVerses int `json:"total_verses"`
}

type BookSchema struct {
	Book      string          `json:"book"`
	Order     int             `json:"order"` // Zum sauberen Sortieren
	Chapters  []ChapterSchema `json:"chapters"`
}