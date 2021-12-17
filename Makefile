build:
	go build -v -o ./bin/note_service ./cmd/note_service/

run: build
	./bin/note_service -dir ./bin/note_data

test:
	go test ./...

