package squashfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type CmdResult struct {
	Result int
	Output []byte
	Error  error
}

func FileExisted(file string) bool {
	var err error
	if _, err = os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func IsDir(file string) bool {
	var err error
	var f os.FileInfo
	if f, err = os.Stat(file); os.IsExist(err) {
		return f.IsDir()
	}
	return err == nil
}

func GetFileList(path string) string {
	var buffer bytes.Buffer
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return ""
	}
	for _, f := range files {
		buffer.WriteString(f.Name())
		buffer.WriteString(" ")
	}
	return buffer.String()
}

func MoveFile(src, dst string) error {
	inf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outf, err := os.Create(dst)
	if err != nil {
		inf.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outf.Close()
	_, err = io.Copy(outf, inf)
	inf.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func AppendFile(file, text string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file:%s", err)
	}

	if _, err = f.WriteString(text); err != nil {
		return fmt.Errorf("Failed to write text: %s", err)
	}
	return f.Close()
}
