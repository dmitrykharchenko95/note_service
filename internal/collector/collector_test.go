package collector

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dmitrykharchenko95/note_service/store"
	"github.com/stretchr/testify/require"
)

func Test_removeNotesInFiles(t *testing.T) {
	path, err := os.MkdirTemp("./", "testdir")
	require.NoError(t, err, "actual err - %v", err)

	file1, err := os.CreateTemp(path, "test")
	require.NoError(t, err, "actual err - %v", err)

	_, err = file1.WriteString(
		`[{"ID":1,"CreateTime":"2021-12-14T12:06:56.575279361+03:00",
		"DeleteTime":"2021-12-14T12:07:01.575279361+03:00","Content":"test note 1"},
		{"ID":2,"CreateTime":"2023-12-14T12:06:56.575279361+03:00",
		"DeleteTime":"2023-12-14T12:07:01.575279361+03:00","Content":"test note 2"}]`)
	require.NoError(t, err, "actual err - %v", err)
	err = file1.Close()
	require.NoError(t, err, "actual err - %v", err)

	file2, err := os.CreateTemp(path, "empty")
	require.NoError(t, err, "actual err - %v", err)
	err = file2.Close()
	require.NoError(t, err, "actual err - %v", err)

	defer func() {
		err := os.RemoveAll(path)
		require.NoError(t, err, "actual err - %v", err)
	}()

	t.Run("base test", func(t *testing.T) {

		store.FilesDirectory = path
		err := removeNotesInFiles()
		require.NoError(t, err, "actual err - %v", err)

		files, err := ioutil.ReadDir(store.FilesDirectory)
		require.NoError(t, err, "actual err - %v", err)
		require.Equal(t, 1, len(files), "files not deleted:", files[0].Name())

		data, err := os.Open(file1.Name())
		require.NoError(t, err, "actual err - %v", err)
		expectData := `[{"ID":2,"CreateTime":"2023-12-14T12:06:56.575279361+03:00","DeleteTime":"2023-12-14T12:07:01.575279361+03:00","Content":"test note 2"}]`

		actualData, err := ioutil.ReadAll(data)
		require.NoError(t, err, "actual err - %v", err)
		require.Equal(t, expectData, string(actualData))
	})

}
