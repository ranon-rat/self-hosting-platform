package database

import (
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Setup() {
	db = sqlx.MustConnect("sqlite3", "file:db/database.db?_foreign_keys=on")
	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(15)
	content := readAndCleanComments()
	executeSQLStatements(content)
}

func readAndCleanComments() string {
	file, err := os.Open("db/init.sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	contentBuffer, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	content := string(contentBuffer)
	// regex \/\*[a-z''A-Z0-9\-\>\s\.\n\t\,\;\_]+\*\/
	re := regexp.MustCompile(`\/\*[\w\d''\-\>\s\n\\t.\,\;\:\_]+\*\/`)
	matches := re.ReplaceAllString(content, "")
	return strings.ReplaceAll(matches, "/**/", "")

}

func executeSQLStatements(content string) {
	sqlStatements := strings.SplitSeq(content, ";")
	for stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			if _, err := db.Exec(stmt); err != nil {
				log.Fatal("Error executing SQL statement: ", err, stmt)
				os.Exit(1)
			}
		}
	}
}
func GetDB() *sqlx.DB {
	return db
}
