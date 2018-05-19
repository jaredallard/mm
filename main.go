package main

import (
	"os"
	"path/filepath"

	"github.com/jaredallard/mm/pkg/commands"
	config "github.com/jaredallard/mm/pkg/config"
	log "github.com/jaredallard/mm/pkg/logger"
	"github.com/jaredallard/mm/pkg/sheets"
	telegram "github.com/jaredallard/mm/pkg/telegram"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to obtain working directory.")
	}

	cfg := config.Load(filepath.Join(workDir, "config.yaml"))

	log.Debug("Setting up Google")
	sheets.Init(&cfg)

	log.Debug("Setting up Telegram")
	telegram.Setup(&cfg)

	log.Debug("Parsing command sheet")

	contents, err := sheets.GetSheet(cfg.Sheet.ID, "A:F")
	if err != nil {
		log.Fatal("Failed to get sheet contents.")
	}

	commands.Initialize(&cfg, contents)

	err = commands.RunCommand(1)
	if err != nil {
		log.Fatal("Failed to run command:", err.Error())
	}

	log.Info("It worked!")
}
