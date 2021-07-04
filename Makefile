GC := go

all: cmd/main.go
	$(GC) build -o touristdb cmd/main.go

walk: ./cmd/walk/walk.go
	$(GC) build ./cmd/walk/walk.go

clean:
	rm -f ./touristdb ./walk
