package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func getCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %v", err)
	}

	hash := strings.TrimSpace(string(out))
	return hash, nil
}

func getCommitTag() (string, error) {
	cmd := exec.Command("git", "tag", "--points-at")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git tag: %v", err)
	}

	tag := strings.TrimSpace(string(out))
	return tag, nil
}

func main() {
	commitHash, err := getCommitHash()
	if err != nil {
		log.Fatalf("get commit hash: %v", err)
	}
	fmt.Printf("commit hash: %s\n", commitHash)

	commitTag, err := getCommitTag()
	if err != nil {
		log.Fatalf("get commit tag: %v", err)
	}
	fmt.Printf("commit tag: %s\n", commitTag)
}
