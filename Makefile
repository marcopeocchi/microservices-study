default:
	go build -o fuu *.go
	mkdir -p build
	mv fuu* build

clean:
	rm -r build

app: 
	cd frontend && pnpm build

multiarch:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o fuu_linux-arm main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o fuu_linux-arm64 main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o fuu_linux-amd64 main.go
	mkdir -p build
	mv fuu* build

linuxamd64:
	GOOS=linux GOARCH=amd64 go build -o fuu *.go