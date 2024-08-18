package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	progressbar "github.com/schollz/progressbar/v3"
)

var log func(string) (int, error)
var logf func(string, ...any) (int, error)
var logln func(...any) (int, error)

func init() {
	log = os.Stderr.WriteString
	logf = func(f string, a ...any) (int, error) {
		return log(fmt.Sprintf(f, a...))
	}
	logln = func(a ...any) (int, error) {
		return log(fmt.Sprintln(a...))
	}
}

func maybe_panic(err error) {
	if err != nil {
		panic(err)
	}
}

func get_block_device_size(fn string) int64 {
	type Disk struct {
		Size int64 `json:"Size"`
	}

	var out bytes.Buffer

	// Construct command
	shell_cmd := fmt.Sprintf("diskutil info -plist %s | plutil -convert json - -o -", fn)
	cmd_s := []string{"sh", "-c", shell_cmd}

	cmd := exec.Command(cmd_s[0], cmd_s[1:]...)
	cmd.Stdout = &out

	// Run command and check for problems
	logf("Running sub-process: %#v\n", cmd_s)
	if err := cmd.Run(); err != nil {
		logln("Error running sub-process: ", err)
	}

	// Parse output
	var d Disk
	err := json.Unmarshal(out.Bytes(), &d)
	maybe_panic(err)

	// Celebrate
	return d.Size
}

func main() {
	fn := "-"

	switch len(os.Args) {
	case 1:
		log("Missing required argument: output destination file.")
		os.Exit(1)
	case 2:
		fn = os.Args[1]
	default:
		logf("Wrong number of arguments! Got %d expecting 1: %v\n", len(os.Args), os.Args)
		os.Exit(1)
	}

	size := get_block_device_size(fn)

	var f *os.File
	if fn == "-" {
		f = os.Stdout
	} else {
		logf("Opening %s for writing ...", fn)
		var err error
		f, err = os.OpenFile(fn, os.O_WRONLY|os.O_SYNC, 0644)
		maybe_panic(err)
	}

	// Create a block of 0xFF bytes
	bs := 256 * 1024
	block := make([]byte, bs)
	for i := range block {
		block[i] = byte(0xFF)
	}
	block_r := bytes.NewReader(block)

	bar := progressbar.DefaultBytes(
		size,
		"zeroing",
	)
	for {
		if _, err := io.Copy(io.MultiWriter(f, bar), block_r); err != nil {
			logf("Error during io.Copy: %s\n", err)
			break
		}
		if _, err := block_r.Seek(0, io.SeekStart); err != nil {
			logf("Error during block_r.Seek: %s\n", err)
			break
		}
	}
	log(fmt.Sprintln("All done!"))
}
