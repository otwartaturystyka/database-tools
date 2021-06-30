GC := go

all: cmd/main.go
	$(GC) build -o touristdb cmd/main.go

notify: ./cmd/notify/notify.go
	$(GC) build ./cmd/notify/notify.go

walk: ./cmd/walk/walk.go
	$(GC) build ./cmd/walk/walk.go

clean:
	rm -f ./touristdb ./notify ./optimize ./walk
