package internal

import "os/exec"

func notify(title, message string) error {
	cmd := exec.Command(
		"terminal-notifier",
		"-title", title,
		"-message", message,
		"-activate", "com.mitchellh.ghostty",
	)
	return cmd.Run()
}
