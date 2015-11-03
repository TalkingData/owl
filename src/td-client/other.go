package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func NewUUID() string {
	var uuid string
	intfs, _ := net.Interfaces()
	for _, i := range intfs {
		uuid += i.HardwareAddr.String()
	}
	hasher := md5.New()
	hasher.Write([]byte(uuid))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetHostName() string {
	name, _ := os.Hostname()
	return strings.Trim(name, "\n")
}

func GetKernel() string {
	res, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return ""
	}
	return strings.Trim(string(res), "\n")
}

func GetOs() string {
	fd, err := os.Open("/etc/redhat-release")
	if err != nil {
		return ""
	}
	render := bufio.NewReader(fd)
	os, err := render.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.Trim(os, "\n")
}

func GetIDRACAddr() string {
	res, err := exec.Command("sh", "-c", "/usr/bin/ipmitool lan print | awk -F ':' '/IP Address.*[1-9]/ {print $2}'").Output()
	if err != nil {
		return ""
	}
	return strings.Trim(string(res), "\n")
}

func DownloadFile(url string) error {
	field := strings.Split(url, "/")
	filename := field[len(field)-1]
	res, err := http.Get(url)
	if err != nil {
		return err

	}
	if res.StatusCode != 200 {
		return fmt.Errorf("get the update file error, status code <%v>", res.StatusCode)
	}
	fd, err := os.Create("./update/" + filename)
	if nil != err {
		return err
	}
	defer fd.Close()
	n, err := io.Copy(fd, res.Body)
	if nil != err {
		return err
	}
	log.Info("download success, write %d byte", n)
	return nil
}

func Unzip(filename string) error {
	pth, err := exec.LookPath("tar")
	if nil != err {
		return fmt.Errorf("cannot find tar command")
	}
	_, err = exec.Command(pth, "xf", "./update/"+filename).Output()
	if nil != err {
		return err
	}
	return nil
}
