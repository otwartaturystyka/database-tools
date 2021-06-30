GC := go
GENERATE := ./cmd/generate
COMPRESS := ./cmd/compress
NOTIFY := ./cmd/notify
UPLOAD := ./cmd/upload

all: cmd/main.go
	$(GC) build -o touristdb cmd/main.go

notify: $(NOTIFY)/notify.go
	go build $(NOTIFY)/notify.go

walk: ./cmd/walk/walk.go
	go build ./cmd/walk/walk.go

clean:
	rm -f ./compress ./generate ./notify ./upload ./walk
