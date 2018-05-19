package state

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/jaredallard/mm/pkg/logger"
)

type stateStruct map[string]string

var statePath string
var stateRep stateStruct

// Init the state store
func Init(path string, projectedLength int) error {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		// ignore
	} else if err != nil {
		return err
	}

	stateRep = make(stateStruct, projectedLength)
	statePath = path

	err = json.Unmarshal(b, &stateRep)

	return err
}

// Set the state store
func Set(id string, value string) {
	stateRep[id] = value

	Serialize()
}

// Bump the state store
func Bump(id string) {
	prev := stateRep[id]
	prevNum, _ := strconv.Atoi(prev)
	stateRep[id] = strconv.Itoa(prevNum + 1)

	Serialize()
}

// Serialize the state store (dump it)
func Serialize() {
	log.Debug("dumping state to", statePath)

	data, err := json.Marshal(stateRep)
	if err != nil {
		log.Fatal("Failed to marshall state to JSON:", err.Error())
	}

	err = ioutil.WriteFile(statePath, data, 0644)
	if err != nil {
		log.Fatal("Failed to save state file:", err.Error())
	}
}

// Get the state store
func Get(id string) string {
	if d, ok := stateRep[id]; !ok || d == "" {
		stateRep[id] = "2"
	}

	return stateRep[id]
}
