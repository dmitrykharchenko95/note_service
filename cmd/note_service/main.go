package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/dmitrykharchenko95/note_service/internal/collector"
	"github.com/dmitrykharchenko95/note_service/internal/service"
	"github.com/dmitrykharchenko95/note_service/store"
)

func init()  {
	flag.DurationVar(&collector.RemovePeriod, "period", time.Minute * 10, "period of start note collector")
	flag.DurationVar(&store.NoteLifetime, "lifetime", time.Hour * 24, "notes lifetime")
	flag.StringVar(&store.FilesDirectory, "dir", "./notes_data", "path to directory for save data")
}

func main()  {
	flag.Parse()

	err := os.Mkdir(store.FilesDirectory, 0777)
	if err != nil && !errors.Is(err, os.ErrExist){
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collector.NoteCollector(ctx)

	OUTLOOP:
	for {
		select {
		case <- ctx.Done():
			break OUTLOOP
		default:
			err := service.NoteService(os.Stdin, os.Stdout, cancel)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}


