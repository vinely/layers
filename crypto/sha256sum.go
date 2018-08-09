package crypto

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Sha256sum(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("Failed to ReadFile:%s\n", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func SaveSha256sum(file, cksumFile string) error {
	sum, err := Sha256sum(file)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%s  %s\n", sum, file)
	f, err := os.OpenFile(cksumFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file: %s", err)
	}
	defer f.Close()
	if _, err = f.WriteString(s); err != nil {
		return fmt.Errorf("Failed to write text: %s", err)
	}
	return nil
}

func GenSha256sum(dir string, cksumFile string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Read Dir[%s] Error: %s\n", dir, err)
	}
	p, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("Get Absolute Directory Error: %s\n", err)
	}
	dest := filepath.Join(p, filepath.Base(cksumFile))
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("Create File Error: %s\n", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			sum, err := Sha256sum(file.Name())
			if err != nil {
				return err
			}
			s := fmt.Sprintf("%s  %s", sum, file.Name())
			if _, err := fmt.Fprintln(f, s); err != nil {
				return fmt.Errorf("Failed write sha256sum for %s: %s\n", file.Name(), err)
			}
		}
	}
	return nil
}

func Check256sum(file string, sum string) bool {
	s, err := Sha256sum(file)
	if err != nil {
		log.Println(err)
		return false
	}
	return strings.Compare(s, sum) == 0
}

func Check256sumFromFile(sumFile string) bool {

	f, err := os.Open(sumFile)
	if err != nil {
		log.Printf("Failed to open file: %s\n", sumFile)
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		bf := bytes.Fields(scanner.Bytes())
		if len(bf) != 2 {
			continue
		}
		file := string(bf[1])
		chksum := string(bf[0])
		if !Check256sum(file, chksum) {
			return false
		} else {
			log.Println(scanner.Text(), "  OK!")
		}
	}

	if scanner.Err() != nil {
		log.Printf("error: %s\n", scanner.Err())
	}

	return true
}
