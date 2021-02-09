package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func openEditor() (string, error) {
	f, err := ioutil.TempFile("", "gitflow-toolkit")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	// write utf8 bom
	bom := []byte{0xef, 0xbb, 0xbf}
	_, err = f.Write(bom)
	if err != nil {
		return "", err
	}

	// get os editor
	var editor string
	if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	} else if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else {
		switch runtime.GOOS {
		case "windows":
			// vscode
			_, err = exec.LookPath("code")
			if err != nil {
				_, err = exec.LookPath("notepad")
				if err != nil {
					return "", fmt.Errorf("cannot find any editor (code/notepad)")
				}
				editor = "notepad"
			}
			editor = "code"
		case "linux", "darwin":
			_, err = exec.LookPath("vim")
			if err != nil {
				_, err = exec.LookPath("vi")
				if err != nil {
					_, err = exec.LookPath("nano")
					if err != nil {
						return "", fmt.Errorf("cannot find any editor (vi/vim/nano)")
					}
					editor = "nano"
				}
				editor = "vi"
			}
			editor = "vim"
		default:
			return "", fmt.Errorf("unsupported platform")
		}
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	raw, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}

	return string(bytes.TrimPrefix(raw, bom)), nil
}

func gitCommand(out io.Writer, cmds []string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("git.exe", cmds...)
	case "linux", "darwin":
		cmd = exec.Command("git", cmds...)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = os.Stdin
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}

	return cmd.Run()
}

func testSend(out io.Writer, cmds []string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("git.exe", cmds...)
	case "linux", "darwin":
		cmd = exec.Command("git", cmds...)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = os.Stdin
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}

	return cmd.Run()
}
