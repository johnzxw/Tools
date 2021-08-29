package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gitFile = "C:/Program Files/Git/bin/git.exe"
var upDir = ""

// InitConfig 对filePath和explodeString进行初始化
func InitConfig() {
	gitFileTmp := flag.String("G", gitFile, " git 文件路径")
	upDirTmp := flag.String("D", upDir, " 需要更新的目录")
	flag.Parse()
	gitFile = *gitFileTmp
	upDir = *upDirTmp

}
func init() {
	InitConfig()
}
func main() {
	if upDir == "" {
		upDir, _ = GetCurrentPath()
	}
	fmt.Println("current dir:" + upDir)
	//获取当前目录下的文件或目录名(包含路径)
	filepathNames, err := filepath.Glob(filepath.Join(upDir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	for i := range filepathNames {
		path := filepathNames[i]
		if path == "." || path == ".." {
			continue
		}
		if isDir(path) && isDir(path+"/.git") {
			mess, err := CmdAndChangeDirToShow(path, gitFile, []string{"pull"})
			if err != nil {
				fmt.Println(mess, err.Error())
			} else {
				fmt.Println(mess)
			}
		}
	}
}
func CmdAndChangeDirToShow(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	f, e := os.Stat(commandName)
	if e != nil {
		fmt.Println(f, e.Error())
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	fmt.Println("path", dir, "Cmd", cmd.Args)

	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err

}
func isDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\"`)
	}
	return path[0 : i+1], nil
}
