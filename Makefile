all:
        GOOS=darwin GOARCH=386 go build -o darwin-i386-qrsctl .
	GOOS=darwin GOARCH=amd64 go build -o darwin-amd64-qrsctl .
	GOOS=windows GOARCH=386 go build -o windows-i386-qrsctl .
	GOOS=windows GOARCH=amd64 go build -o windows-amd64-qrsctl .
