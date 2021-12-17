package store

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	type args struct {
		userLogin string
	}
	tests := []struct {
		name    string
		args    args
		want    *Notes
		wantErr bool
	}{{
		name: "base",
		args: args{
			userLogin: "test",
		},
		want: &Notes{
			Note{
				ID:         1,
				CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
				DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
				Content:    "test note #1",
			},
			Note{
				ID:         2,
				CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
				DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
				Content:    "test note #2",
			},
		},
		wantErr: false,
	}, {
		name: "empty",
		args: args{
			userLogin: "testempty",
		},
		want:    &Notes{},
		wantErr: false,
	}}
	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				path, err := os.MkdirTemp("./", "testdir")
				require.NoError(t, err)

				defer func() {
					err := os.RemoveAll(path)
					require.NoError(t, err, "actual err - %v", err)
				}()

				FilesDirectory = path

				data, err := json.Marshal(Notes{
					Note{
						ID:         1,
						CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
						DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
						Content:    "test note #1",
					},
					Note{
						ID:         2,
						CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
						DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
						Content:    "test note #2",
					}})
				require.NoError(t, err, "actual err - %v", err)

				err = os.WriteFile(path+"/"+tt.args.userLogin, data, 0666)
				require.NoError(t, err, "actual err - %v", err)

				got, err := Load(tt.args.userLogin)
				if (err != nil) != tt.wantErr {
					t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Load() got = %v, want %v", got, tt.want)
				}
			})
		} else {
			t.Run(tt.name, func(t *testing.T) {
				path, err := os.MkdirTemp("./", "testdir")
				require.NoError(t, err)

				defer func() {
					err := os.RemoveAll(path)
					require.NoError(t, err, "actual err - %v", err)
				}()

				FilesDirectory = path

				got, err := Load(tt.args.userLogin)
				if (err != nil) != tt.wantErr {
					t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Load() got = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestSave(t *testing.T) {
	type args struct {
		userLogin string
		notes     Notes
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: "new user",
		args: args{
			userLogin: "user1",
			notes: Notes{
				Note{
					ID:         1,
					CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
					DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
					Content:    "test note #1",
				},
				Note{
					ID:         2,
					CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.Local),
					DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.Local),
					Content:    "test note #2",
				},
			},
		},
		wantErr: false,
	}, {
		name: "old user",
		args: args{
			userLogin: "user1",
			notes: Notes{
				Note{
					ID:         1,
					CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.FixedZone("UTC-8", 0)),
					DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.FixedZone("UTC-8", 0)),
					Content:    "test note #1",
				},
				Note{
					ID:         2,
					CreateTime: time.Date(2022, time.December, 23, 12, 00, 00, 00, time.FixedZone("UTC-8", 0)),
					DeleteTime: time.Date(2023, time.December, 23, 12, 00, 00, 00, time.FixedZone("UTC-8", 0)),
					Content:    "test note #2",
				},
			},
		},
		wantErr: false,
	},
	}
	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				path, err := os.MkdirTemp("./", "testdir")
				require.NoError(t, err)

				defer func() {
					err := os.RemoveAll(path)
					require.NoError(t, err, "actual err - %v", err)
				}()

				FilesDirectory = path
				if err := Save(tt.args.userLogin, tt.args.notes); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}

				actualData, err := os.ReadFile(path + "/" + tt.args.userLogin)
				require.NoError(t, err, "actual err - %v", err)

				expectedData := `[{"ID":1,"CreateTime":"2022-12-23T12:00:00+03:00","DeleteTime":"2023-12-23T12:00:00+03:00","Content":"test note #1"},{"ID":2,"CreateTime":"2022-12-23T12:00:00+03:00","DeleteTime":"2023-12-23T12:00:00+03:00","Content":"test note #2"}]`
				require.Equal(t, expectedData, string(actualData))

			})
		} else {
			t.Run(tt.name, func(t *testing.T) {
				path, err := os.MkdirTemp("./", "testdir")
				require.NoError(t, err)

				defer func() {
					err := os.RemoveAll(path)
					require.NoError(t, err, "actual err - %v", err)
				}()

				file, err := os.Create(path + "/" + tt.args.userLogin)
				require.NoError(t, err, "actual err - %v", err)

				note, err := json.Marshal(Notes{
					Note{
						ID:         1,
						CreateTime: time.Date(2022, time.December, 23, 13, 00, 00, 00, time.FixedZone("UTC-8", 0)),
						DeleteTime: time.Date(2023, time.December, 23, 13, 00, 00, 00, time.FixedZone("UTC-8", 0)),
						Content:    "test note #2",
					}})
				require.NoError(t, err)

				_, err = file.Write(note)
				require.NoError(t, err)

				FilesDirectory = path
				if err := Save(tt.args.userLogin, tt.args.notes); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}

				actualData, err := os.ReadFile(path + "/" + tt.args.userLogin)
				require.NoError(t, err, "actual err - %v", err)

				expectedData := `[{"ID":1,"CreateTime":"2022-12-23T12:00:00Z","DeleteTime":"2023-12-23T12:00:00Z","Content":"test note #1"},{"ID":2,"CreateTime":"2022-12-23T12:00:00Z","DeleteTime":"2023-12-23T12:00:00Z","Content":"test note #2"}]`
				require.Equal(t, expectedData, string(actualData))

			})
		}
	}
}
