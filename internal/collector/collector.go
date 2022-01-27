package collector

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dmitrykharchenko95/note_service/store"
)

var RemovePeriod time.Duration

// NoteCollector start removeNotesInFiles every removePeriod while the program is running
func NoteCollector(ctx context.Context) {
	ticker := time.NewTicker(RemovePeriod)
	for {
		select {
		case <-ticker.C:
			err := removeNotesInFiles()
			if err != nil {
				log.Printf("NoteCollector err: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// removeNotesInFiles iterates over the files in the store.FilesDirectory directory and deletes notes if value of field
// DeleteTime less than time.Now. If there are no notes in the file, the file will be deleted.
func removeNotesInFiles() error {
	files, err := ioutil.ReadDir(store.FilesDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		notes, err := store.Load(file.Name())
		if err != nil {
			return err
		}

		if len(*notes) == 0 {
			if err = os.Remove(fmt.Sprintf("%s/%s", store.FilesDirectory, file.Name())); err!= nil{
				return err
			}
			continue
		}

		var newNotes store.Notes
		for _, note := range *notes {
			if note.DeleteTime.Before(time.Now()) {
				continue
			}
			newNotes = append(newNotes, note)
		}

		if err = store.Save(file.Name(), newNotes); err != nil {
			log.Printf("NoteCollector save err: %v", err)
			return err
		}
	}
	return nil
}
