default:
	go build -o fuu -o fuu cmd/server/main.go
	mkdir -p build
	mv fuu* build

clean:
	rm -r build

app: 
	cd cmd/server/frontend && pnpm build

multiarch:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o fuu_linux-arm cmd/server/main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o fuu_linux-arm64 cmd/server/main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o fuu_linux-amd64 cmd/server/main.go
	mkdir -p build
	mv fuu* build

linuxamd64:
	GOOS=linux GOARCH=amd64 go build -o fuu cmd/server/main.go

dev:
	go run cmd/server/main.go -c "/Users/marco/dev/homebrew/fuu/Fuufile"