package service

import (
	"testing"
	"time"

	"github.com/dmitrykharchenko95/note_service/store"
	"github.com/stretchr/testify/require"
)

var testNotes = store.Notes{
	store.Note{
		ID:         1,
		CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
		DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
		Content:    "test note #1",
	},
	store.Note{
		ID:         2,
		CreateTime: time.Date(2022, time.December, 24, 12, 00, 00, 00, time.Local),
		DeleteTime: time.Date(2023, time.December, 24, 12, 00, 00, 00, time.Local),
		Content:    "test note #2",
	}}

func TestCreateNote(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		actualNotes := testNotes
		CreateNote(&actualNotes, time.Minute, "test note #3")
		expectedNotes := store.Notes{
			store.Note{
				ID:         1,
				CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
				DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
				Content:    "test note #1",
			},
			store.Note{
				ID:         2,
				CreateTime: time.Date(2022, time.December, 24, 12, 00, 00, 00, time.Local),
				DeleteTime: time.Date(2023, time.December, 24, 12, 00, 00, 00, time.Local),
				Content:    "test note #2",
			}, store.Note{
				ID:         3,
				CreateTime: time.Now(),
				DeleteTime: time.Now().Add(time.Minute),
				Content:    "test note #3",
			}}
		require.Equal(t, expectedNotes[:1], actualNotes[:1], "Notes no equal")
		require.Equal(t, expectedNotes[2].ID, actualNotes[2].ID, "Notes no equal")
		require.Equal(t, expectedNotes[2].Content, actualNotes[2].Content, "Notes no equal")
		require.Equal(t, expectedNotes[2].CreateTime.Round(time.Second),
			actualNotes[2].CreateTime.Round(time.Second), "Notes no equal")
		require.Equal(t, expectedNotes[2].DeleteTime.Round(time.Second),
			actualNotes[2].DeleteTime.Round(time.Second), "Notes no equal")
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		actualNotes := make(store.Notes, 2)
		copy(actualNotes, testNotes)

		DeleteNote(&actualNotes, 1)

		expectedNotes := store.Notes{
			store.Note{
				ID:         2,
				CreateTime: time.Date(2022, time.December, 24, 12, 00, 00, 00, time.Local),
				DeleteTime: time.Date(2023, time.December, 24, 12, 00, 00, 00, time.Local),
				Content:    "test note #2",
			}}
		require.Equal(t, expectedNotes, actualNotes, "Notes no equal")
	})
}

func TestGetAll(t *testing.T) {

	t.Run("base", func(t *testing.T) {
		expectedData := "id: 2\tSat, 24 Dec 12:00:00:\ntest note #2\nThe note will be deleted at Sun, 24 Dec 12:00:00\n" +
			"\nid: 1\tFri, 23 Dec 12:00:00:\ntest note #1\nThe note will be deleted at Sat, 23 Dec 12:00:00\n\n"

		actualData, err := GetAll(&testNotes)

		require.NoError(t, err)
		require.Equal(t, expectedData, actualData, "Notes no equal")
	})
}

func TestGetLast(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		expectedData := "id: 2\tSat, 24 Dec 12:00:00:\ntest note #2\nThe note will be deleted at Sun, 24 Dec 12:00:00\n\n"

		actualData := GetLast(&testNotes)

		require.Equal(t, expectedData, actualData, "Notes no equal")
	})
}

func TestGetOldest(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		expectedData := "id: 1\tFri, 23 Dec 12:00:00:\ntest note #1\nThe note will be deleted at Sat, 23 Dec 12:00:00\n\n"

		actualData := GetOldest(&testNotes)

		require.Equal(t, expectedData, actualData, "Notes no equal")
	})
}
