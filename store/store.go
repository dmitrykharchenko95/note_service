package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Note is struct of user's note
type Note struct {
	ID         int
	CreateTime time.Time
	DeleteTime time.Time
	Content    string
}

var (
	FilesDirectory string
	NoteLifetime   time.Duration
)

type Notes []Note

// NewNote creates a new note
func NewNote() *Note {
	return &Note{CreateTime: time.Now()}
}

// Save saves the user's notes in JSON format from notes to a file userLogin in the directory ./notes_data. If the file does not
//exist, a new file will be created.
func Save(userLogin string, notes Notes) error {
	filePath := fmt.Sprintf("%s/%s", FilesDirectory, userLogin)

	err := os.Remove(filePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	userData, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	_, err = file.Write(userData)
	if err != nil {
		return err
	}

	return nil
}

// Load loads data from file userLogin in directory ./notes_data and return pointer for Notes.  If the file does not
// exist, Load return pointer for empty Notes and a new file will be created.
func Load(userLogin string) (*Notes, error) {
	filePath := fmt.Sprintf("%s/%s", FilesDirectory, userLogin)

	file, err := os.Open(filePath)

	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(filePath)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		return &Notes{}, err
	} else if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return new(Notes), nil
	}

	loadNotes := make(Notes, 0, 16)
	err = json.Unmarshal(data, &loadNotes)
	if err != nil {
		return nil, err
	}

	return &loadNotes, nil
}
