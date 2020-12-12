generate:
	go build ./cmd/generate/generate.go ./cmd/generate/parsers.go

compress:
	go build ./cmd/compress/compress.go

upload:
	go build ./cmd/upload/upload.go ./cmd/upload/parsers.go

walk:
	go build ./cmd/walk/walk.go

all: generate compress upload
