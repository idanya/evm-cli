package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func openInEditor(text []byte) error {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "evm-cli-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal("Cannot find editor", err)
	}

	cmd := exec.Command(executable, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
