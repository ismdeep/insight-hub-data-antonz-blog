package core

import (
	"errors"
	"io"
	"strings"
)

type Store struct {
	loaded map[string]bool
	output io.Writer
}

func NewStore(writer io.Writer) *Store {
	return &Store{
		loaded: make(map[string]bool),
		output: writer,
	}
}

func (receiver *Store) Load(reader io.Reader) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		items := strings.Split(line, "|")

		receiver.loaded[items[0]] = true
	}

	return nil
}

func (receiver *Store) Exists(id string) bool {
	return receiver.loaded[id]
}

func (receiver *Store) URLExists(link string) bool {
	return receiver.loaded[LinkToRecordID(link)]
}

func LinkIsTidy(link string) bool {
	link = strings.ReplaceAll(link, "https://", "")
	link = strings.ReplaceAll(link, "http://", "")
	if strings.Contains(link, "//") {
		return false
	}
	return true
}

func (receiver *Store) Write(r Record) error {

	if receiver.output == nil {
		return ErrOutputIsNil
	}

	if !LinkIsTidy(r.Link) {
		return ErrLinkIsNotTidy
	}

	id := RecordID(r)

	if receiver.loaded[id] {
		return ErrAlreadyExists
	}

	if _, err := receiver.output.Write([]byte(RecordMarshal(r) + "\n")); err != nil {
		return err
	}

	receiver.loaded[id] = true

	return nil
}

func (receiver *Store) Save(r Record) error {
	if err := receiver.Write(r); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return nil
		}
		return err
	}
	return nil
}
