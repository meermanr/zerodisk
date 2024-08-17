zerodisk: Makefile
zerodisk: $(wildcard *.go,go.*)
	env GOOS=darwin GOARCH=arm64 go build -o zerodisk .

clean:
	git clean --force -X	# Removed ignored files only
