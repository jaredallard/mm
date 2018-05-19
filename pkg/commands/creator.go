package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jaredallard/mm/pkg/state"
	"github.com/jaredallard/mm/pkg/telegram"

	"github.com/jaredallard/mm/pkg/config"
	log "github.com/jaredallard/mm/pkg/logger"
	"github.com/jaredallard/mm/pkg/sheets"
	cron "github.com/robfig/cron"
)

// Command is a command structure
type Command struct {
	Column   string
	Date     string // this is the cron format
	SheetID  string
	Timezone string
	Comment  string
	Text     string
}

var commandTable []Command
var cfg *config.ConfigurationFile
var cronTable *cron.Cron

func initCron(cfg *config.ConfigurationFile) {
	l, _ := time.LoadLocation(cfg.Sheet.Timezone)
	cronTable = cron.NewWithLocation(l)

	log.Info("cron configured to run with timezone:", cfg.Sheet.Timezone)
}

// Initialize the command table
func Initialize(c *config.ConfigurationFile, sheet sheets.SheetContents) {
	log.Debug("started command init")

	cfg = c

	initCron(cfg)

	nc := numCmd(sheet)
	commandTable = make([]Command, nc)

	commands := 0
	for i := range sheet.Values {
		commandSlice := sheet.Values[i]

		// skip the first index, that's our mapper (for js)
		if i == 0 {
			continue
		}

		// filter out semi-bar data.
		if len(commandSlice) < 4 {
			continue
		}

		// text is optional, FIXME
		text := ""
		if len(commandSlice) == 6 {
			text = commandSlice[5]
		}

		commandTable[commands] = Command{
			Column:   commandSlice[0],
			Date:     commandSlice[1],
			SheetID:  commandSlice[2],
			Timezone: "",
			Comment:  commandSlice[4],
			Text:     text,
		}

		cmd := commandTable[commands]

		cronTable.AddFunc(cmd.Date, func() {
			nid := commands - 1
			RunCommand(nid)
		})

		log.Info("Registered command '" + cmd.Comment + "' to run at interval '" + cmd.Date + "'")

		commands++
	}

	log.Debug("commands initialized")

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to obtain working directory.")
	}
	log.Debug("Setting up state")
	state.Init(filepath.Join(workDir, "state.json"), commands)
}

// Get the number of valid commands in a Google Sheets splice
func numCmd(sheet sheets.SheetContents) int {
	valid := 0
	for i := range sheet.Values {
		// skip the first index, that's our mapper (for js)
		if i == 0 {
			continue
		}

		// filter out semi-bar data.
		if len(sheet.Values[i]) < 4 {
			continue
		}

		valid++
	}

	return valid
}

// RunCommand <- name implies is
func RunCommand(index int) error {
	if len(commandTable) < index {
		return errors.New("Index is out of range")
	}

	strIndex := strconv.Itoa(index)
	log.Debug("Executing command at index", strIndex)

	pointer := commandTable[index]

	// Handle pure-text posts.
	if pointer.Text != "" {
		err := telegram.SendToChannel(cfg.Telegram.ChannelID, pointer.Text)
		if err != nil {
			return err
		}

		return nil
	}

	prev := state.Get(strIndex)
	log.Debug("Iteration ID:", prev)

	contents, err := sheets.GetRange(pointer.SheetID, pointer.Column+prev+":"+pointer.Column+prev)
	if err != nil {
		return err
	}

	if len(contents) == 0 || contents[0] == "" {
		log.Debug("Empty contents for task, skipping and not bumping state index.")
		return nil
	}

	err = telegram.SendToChannel(cfg.Telegram.ChannelID, contents[0])
	if err != nil {
		return err
	}

	state.Bump(strIndex)

	return nil
}

// Start the cron watcher
func Start() {
	cronTable.AddFunc("@every 1h", func() {
		log.Debug("polling Google Sheet for new instructions")

		contents, err := sheets.GetSheet(cfg.Sheet.ID, "A:F")
		if err != nil {
			log.Fatal("Failed to get sheet contents.")
		}

		// FIXME: Only detects basic addition / removal, and replaces the entire struct
		nc := numCmd(contents)
		if nc != len(commandTable) {
			log.Info("Refreshing command list due to changes.")
			cronTable.Stop()
			Initialize(cfg, contents)
		}
	})
	cronTable.Start()
}
