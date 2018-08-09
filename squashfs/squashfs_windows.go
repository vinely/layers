package squashfs

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	Paths         = []string{"."}
	MksquashfsCmd = "mksquashfs.exe"
	UnsquashfsCmd = "unsquashfs.exe"
)

func ExtraSquashfsPackage(pack string, target string) error {
	if FileExisted(target) {
		return errors.New(fmt.Sprintf("%s already existed!", target))
	}
	cmd := fmt.Sprintf("%s -d %s %s \n", UnsquashfsCmd, target, pack)
	//log.Println(cmd)
	out := Cmd(cmd)
	if out.Result != 0 {
		return out.Error
	}
	return nil
}

func MakeSquashfsPackage(src string, pack string) error {
	if !IsDir(src) {
		return errors.New(fmt.Sprintf("%s is not a valid directory", src))
	}
	list := GetFileList(src)
	cmd := fmt.Sprintf("cd %s && %s %s %s -comp xz -b 512k -keep-as-directory\n", src, MksquashfsCmd, list, pack)
	//log.Println(cmd)
	out := Cmd(cmd)
	if out.Result != 0 {
		return out.Error
	}
	return nil
}

func CheckUnSquashfs() bool {
	var ok bool
	ok, UnsquashfsCmd = CheckCommandExisted(Paths, "unsquashfs.exe")
	// if ok {
	// 	fmt.Printf("%s existed!\n", MksquashfsCmd)
	// }
	return ok
}

func CheckMkSquashfs() bool {
	var ok bool
	ok, MksquashfsCmd = CheckCommandExisted(Paths, "mksquashfs.exe")
	if ok {
		//fmt.Printf("%s existed!\n", MksquashfsCmd)
		out := Cmd(MksquashfsCmd + " -version")
		if out.Result == 0 {
			//fmt.Println(string(out.Output))
		} else {
			//fmt.Println(string(out.Error.Error()))
			ok = false
		}
	}
	return ok
}

func CheckSquashfsTools() bool {
	if CheckMkSquashfs() && CheckUnSquashfs() {
		return true
	} else {
		return false
	}
}

func CheckCommandExisted(paths []string, name string) (bool, string) {
	out := Cmd("where " + name)
	if out.Result == 0 {
		fds := bytes.Fields(out.Output)
		if len(fds) > 0 {
			return true, string(fds[0])
		}
	}
	for _, p := range Paths {
		file := filepath.Join(p, name)
		if FileExisted(file) {
			return true, file
		}
	}
	return false, ""
}

func Cmd(cmd string) *CmdResult {
	exitCode := 0
	out, err := exec.Command("cmd", "/C", cmd).Output()
	if err != nil {
		if p, ok := err.(*exec.ExitError); ok {
			ws := p.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = -1
		}
		return &CmdResult{
			Result: exitCode,
			Output: nil,
			Error:  err,
		}
	}
	return &CmdResult{
		Result: 0,
		Output: out,
		Error:  nil,
	}
}
