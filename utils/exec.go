package utils

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func ExecWithDir(cmdLine string, dir string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdLine)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func ExecSyncOutput(cmdLine string, dir string) error {
	cmd := exec.Command("bash", "-c", cmdLine)
	cmd.Dir = dir
	cmdStdoutPipe, _ := cmd.StdoutPipe()
	cmdStderrPipe, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}
	go syncLog(cmdStdoutPipe)
	go syncLog(cmdStderrPipe)
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func syncLog(reader io.ReadCloser) {
	// f, _ := os.OpenFile("file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// defer f.Close()
	buf := make([]byte, 1024, 1024)
	for {
		strNum, err := reader.Read(buf)
		if strNum > 0 {
			outputByte := buf[:strNum]
			// f.WriteString(string(outputByte))
			fmt.Printf(string(outputByte))
		}
		if err != nil {
			// read tail
			if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
				err = nil
			}
		}
	}
}
