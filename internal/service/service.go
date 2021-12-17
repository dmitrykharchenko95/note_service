package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dmitrykharchenko95/note_service/internal/collector"
	"github.com/dmitrykharchenko95/note_service/store"
)

var login string

// CreateNote creates and adds a new note to the notes with a lifetime t and content
func CreateNote(notes *store.Notes, t time.Duration, content string) {
	note := store.NewNote()
	if len(*notes) == 0 {
		note.ID = 1
	} else {
		note.ID = (*notes)[len(*notes)-1].ID + 1
	}

	note.Content = content
	note.DeleteTime = note.CreateTime.Add(t)
	*notes = append(*notes, *note)
	fmt.Println("Note added!")
}

// DeleteNote deletes note with id equal to noteID
func DeleteNote(notes *store.Notes, noteID int) {
	for i, n := range *notes {
		if n.ID != noteID {
			continue
		}
		*notes = append((*notes)[:i], (*notes)[i+1:]...)
		fmt.Println("Note deleted!")
		return
	}
}

// GetAll returns all notes from notes as a string
func GetAll(notes *store.Notes) (string, error) {

	var sb strings.Builder
	var s string

	for i := len(*notes) - 1; i >= 0; i-- {
		s = fmt.Sprintf("id: %v\t%v:\n%v\nThe note will be deleted at %v\n\n", (*notes)[i].ID,
			(*notes)[i].CreateTime.Format("Mon, 2 Jan 15:04:05"), (*notes)[i].Content,
			(*notes)[i].DeleteTime.Round(collector.RemovePeriod).Format("Mon, 2 Jan 15:04:05"))
		_, err := sb.WriteString(s)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

// GetOldest returns the oldest note from notes as a string
func GetOldest(notes *store.Notes) string {
	return fmt.Sprintf("id: %v\t%v:\n%v\nThe note will be deleted at %v\n\n", (*notes)[0].ID,
		(*notes)[0].CreateTime.Format("Mon, 2 Jan 15:04:05"), (*notes)[0].Content,
		(*notes)[0].DeleteTime.Round(collector.RemovePeriod).Format("Mon, 2 Jan 15:04:05"))
}

// GetLast returns the newest note from notes as a string
func GetLast(notes *store.Notes) string {
	l := len(*notes) - 1
	return fmt.Sprintf("id: %v\t%v:\n%v\nThe note will be deleted at %v\n\n", (*notes)[l].ID,
		(*notes)[l].CreateTime.Format("Mon, 2 Jan 15:04:05"), (*notes)[l].Content,
		(*notes)[l].DeleteTime.Round(collector.RemovePeriod).Format("Mon, 2 Jan 15:04:05"))
}

// NoteService accepts commands from in and processes the data depending on the command. For authorization,
// enter your username. After NoteService accepts the following commands:
// add [lifetime] [content] - add a note;
// del [id] - delete note with id;
// get - show all notes of user;
// last - show the last note of user;
// old - show the oldest note of user;
// out - logout;
// q - stop the program
// Responses of the program is recorded in out.
func NoteService(in io.ReadCloser, out io.Writer, cancel context.CancelFunc) error {
	scanner := bufio.NewScanner(in)

	fmt.Println(`Enter Login or "Q" for quit`)

	for scanner.Scan() {
		login = scanner.Text()
		break
	}

	if strings.ToUpper(login) == "Q" {
		cancel()
		return nil
	}

	fmt.Println("Log in as", login)

	for scanner.Scan() {
		args := strings.SplitN(scanner.Text(), " ", 3)

		switch strings.ToUpper(args[0]) {
		case "ADD":
			if len(args) == 1 {
				fmt.Println("Enter 'add [lifetime] [content]'")
				continue
			}

			lt, err := time.ParseDuration(args[1])
			if err != nil {
				lt = store.NoteLifetime
				args = args[1:]
			} else {
				args = args[2:]
			}

			loadNotes, err := store.Load(login)
			if err != nil {
				return err
			}

			CreateNote(loadNotes, lt, args[0])
			err = store.Save(login, *loadNotes)
			if err != nil {
				log.Println(err)
			}

		case "DEL":
			if len(args) == 1 {
				fmt.Println("Enter 'del [note id]'")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				log.Println(err)
				continue
			}
			loadNotes, err := store.Load(login)
			if err != nil {
				return err
			}

			DeleteNote(loadNotes, id)

			err = store.Save(login, *loadNotes)
			if err != nil {
				log.Println(err)
			}

		case "GET":
			loadNotes, err := store.Load(login)
			if err != nil {
				return err
			}

			notes, err := GetAll(loadNotes)
			if err != nil {
				log.Fatal(err)
			}

			_, err = out.Write([]byte(notes))
			if err != nil {
				log.Fatal(err)
			}

		case "OLD":
			loadNotes, err := store.Load(login)
			if err != nil {
				return err
			}
			_, err = out.Write([]byte(GetOldest(loadNotes)))
			if err != nil {
				log.Fatal(err)
			}

		case "LAST":
			loadNotes, err := store.Load(login)
			if err != nil {
				return err
			}
			_, err = out.Write([]byte(GetLast(loadNotes)))
			if err != nil {
				log.Fatal(err)
			}

		case "OUT":
			login = ""
			return nil

		case "Q":
			cancel()
			return nil
		default:
			_, err := out.Write([]byte("Unknown command"))
			if err != nil {
				log.Fatal(err)
			}
		}
		continue
	}
	return nil
}
