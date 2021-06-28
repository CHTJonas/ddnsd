package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"

	zonefile "github.com/bwesterb/go-zonefile"
)

func updateResourceRecord(username, contents string) error {
	// Open zonefile
	f, err := os.OpenFile(zonefilePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// Load zonefile
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// Parse zonefile
	zf, err := zonefile.Load(data)
	if err != nil {
		return err
	}

	for _, e := range zf.Entries() {
		// Find RR
		if !bytes.Equal(e.Domain(), []byte(username)) {
			continue
		}

		// Check RR
		if !bytes.Equal(e.Type(), []byte("TXT")) {
			return errors.New("resource record type in zonefile is not TXT")
		}

		// Update RR
		e.SetValue(0, []byte(contents))

		// Write zonefile
		err = f.Truncate(0)
		if err != nil {
			return err
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}
		_, err = f.Write(zf.Save())
		return err
	}

	return errors.New("could not find resource record in zonefile")
}

func callHook(path string) error {
	cmd := exec.Command(path)
	_, err := cmd.Output()
	return err
}
