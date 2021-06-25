package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	zonefile "github.com/bwesterb/go-zonefile"
)

func updateResourceRecord(username, contents string) error {
	// Load zonefile
	data, err := ioutil.ReadFile(zonefilePath)
	if err != nil {
		return err
	}

	// Parse zonefile
	zf, err := zonefile.Load(data)
	if err != nil {
		return err
	}

	// Update RR
	for _, e := range zf.Entries() {
		if !bytes.Equal(e.Domain(), []byte(username)) {
			continue
		}
		if !bytes.Equal(e.Type(), []byte("TXT")) {
			return errors.New("resource record type in zonefile is not TXT")
		}
		e.SetValue(0, []byte(contents))
		fh, err := os.OpenFile(zonefilePath, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		_, err = fh.Write(zf.Save())
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("could not find resource record in zonefile")
}

func callHook(path string) error {
	cmd := exec.Command(path)
	_, err := cmd.Output()
	return err
}
