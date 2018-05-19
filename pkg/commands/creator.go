package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

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

// Initialize the command table
func Initialize(c *config.ConfigurationFile, sheet sheets.SheetContents) {
	log.Debug("started command init")

	cfg = c
	cronTable = cron.New()

	// FIXME: Will be larger than needed.
	commands := 0
	for i := range sheet.Values {
		if i == 0 {
			continue
		}
		if len(sheet.Values[i]) < 4 {
			continue
		}
		if sheet.Values[i][1] == "" {
			continue
		}
		commands++
	}

	// use the prealloc
	commandTable = make([]Command, commands)

	commands = 0
	for i := range sheet.Values {
		commandSlice := sheet.Values[i]

		// pre-checks
		if i == 0 {
			continue
		}
		if len(commandSlice) < 4 {
			continue
		}
		if sheet.Values[i][1] == "" {
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
			Timezone: commandSlice[3],
			Comment:  commandSlice[4],
			Text:     text,
		}

		cmd := commandTable[commands]

		cronTable.AddFunc(cmd.Date, func() {
			id := strconv.Itoa(commands)
			log.Debug("running command, via cron:", id)

			RunCommand(commands)
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

	err = telegram.SendToChannel(cfg.Telegram.ChannelID, contents[0])
	if err != nil {
		return err
	}

	state.Bump(strIndex)

	return nil
}

// Start the cron watcher
func Start() {
	cronTable.Start()
}
