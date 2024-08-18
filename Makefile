zerodisk: Makefile
zerodisk: go.*	# go.mod, go.sum, etc
zerodisk: *.go 	# source code proper
	go mod tidy
	env GOOS=darwin GOARCH=arm64 go build -o zerodisk.darwin-arm64 .
	env go build -o zerodisk .

clean:
	git clean --force -X	# Removed ignored files only
