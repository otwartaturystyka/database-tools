all: generate compress notify upload walk

GENERATE := ./cmd/generate
COMPRESS := ./cmd/compress
NOTIFY := ./cmd/notify
UPLOAD := ./cmd/upload

generate: $(GENERATE)/generate.go $(GENERATE)/parsers.go
	go build $(GENERATE)/generate.go $(GENERATE)/parsers.go

compress: $(COMPRESS)/compress.go
	go build $(COMPRESS)/compress.go

notify: $(NOTIFY)/notify.go
	go build $(NOTIFY)/notify.go

upload: $(UPLOAD)/upload.go $(UPLOAD)/parsers.go $(UPLOAD)/types.go
	go build $(UPLOAD)/upload.go $(UPLOAD)/parsers.go $(UPLOAD)/types.go

walk: ./cmd/walk/walk.go
	go build ./cmd/walk/walk.go

clean:
	rm -f ./compress ./generate ./upload ./walk
