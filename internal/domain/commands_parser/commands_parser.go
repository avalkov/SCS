package commandsparser

import (
	"errors"
	"regexp"
)

type CommandType int

const (
	AddItem CommandType = iota
	DeleteItem
	GetItem
	GetAllItems
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
}

type commandsParser struct {
	addItemPattern     *regexp.Regexp
	deleteItemPattern  *regexp.Regexp
	getItemPattern     *regexp.Regexp
	getAllItemsPattern *regexp.Regexp
}

// NewCommandsParser creates a new commandsParser
func NewCommandsParser() *commandsParser {
	return &commandsParser{
		addItemPattern:     regexp.MustCompile(`^addItem\(\s*'([^']*)'\s*,\s*'([^']*)'\s*\)$`),
		deleteItemPattern:  regexp.MustCompile(`^deleteItem\(\s*'([^']*)'\s*\)$`),
		getItemPattern:     regexp.MustCompile(`^getItem\(\s*'([^']*)'\s*\)$`),
		getAllItemsPattern: regexp.MustCompile(`^getAllItems\(\)$`),
	}
}

// ParseCommand parses the command string and returns the key used for consistent hashing,
// the parsed Command, and an error if any.
func (cp *commandsParser) parseCommand(command string) (string, Command, error) {
	if cp.addItemPattern.MatchString(command) {
		matches := cp.addItemPattern.FindStringSubmatch(command)
		return matches[1], Command{Type: AddItem, Key: matches[1], Value: matches[2]}, nil
	} else if cp.deleteItemPattern.MatchString(command) {
		matches := cp.deleteItemPattern.FindStringSubmatch(command)
		return matches[1], Command{Type: DeleteItem, Key: matches[1]}, nil
	} else if cp.getItemPattern.MatchString(command) {
		matches := cp.getItemPattern.FindStringSubmatch(command)
		return matches[1], Command{Type: GetItem, Key: matches[1]}, nil
	} else if cp.getAllItemsPattern.MatchString(command) {
		return "", Command{Type: GetAllItems}, nil
	}
	return "", Command{}, errors.New("invalid command format")
}

// ParseCommand parses the command string and returns the parsed Command and an error if any.
func (cp *commandsParser) ParseCommand(command string) (Command, error) {
	_, cmd, err := cp.parseCommand(command)
	return cmd, err
}

// GetCommandID returns the key used for consistent hashing from the command string.
func (cp *commandsParser) GetCommandID(command string) (string, error) {
	key, _, err := cp.parseCommand(command)
	return key, err
}
