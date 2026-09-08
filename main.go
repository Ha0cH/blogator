package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Ha0cH/blogator/internal/config"
	"github.com/Ha0cH/blogator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.ReadConfigJson()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	var s state
	s.cfg = &cfg
	s.db = dbQueries

	var cmds commands
	cmds.cmdLookUp = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Program requires at least a command, and its arguments if needed.")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	args := os.Args[2:]
	cmd := command{name: cmdName, args: args}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalf("Error executing command: %v\n", err)
	}

}
