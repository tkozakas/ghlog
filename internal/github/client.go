package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
)

var (
	ErrGHNotInstalled     = errors.New("gh CLI is not installed")
	ErrGHNotAuthenticated = errors.New("gh CLI is not authenticated")
)

func CheckGHInstalled() error {
	if err := exec.Command("gh", "--version").Run(); err != nil {
		return ErrGHNotInstalled
	}
	return nil
}

func CheckGHAuthenticated() error {
	if err := exec.Command("gh", "auth", "status").Run(); err != nil {
		return ErrGHNotAuthenticated
	}
	return nil
}

func runGH(args ...string) ([]byte, error) {
	cmd := exec.Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String())
	}
	return stdout.Bytes(), nil
}

func runGHWithJSON(result any, args ...string) error {
	output, err := runGH(args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(output, result)
}
