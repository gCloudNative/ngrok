package util

import (
	"crypto/rand"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func ReadTokenFile(filename string) []string {
	tokens := make([]string, 0)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("error:", err)
		return tokens
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if utf8.RuneCountInString(line) == 0 {
			continue
		}
		tokens = append(tokens, strings.TrimSpace(line))
	}
	return tokens
}

func ExecCmd(command string) string {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stdout)

	if err != nil {
		log.Fatal("error:", err)
		return ""
	}
	return strings.TrimSpace(string(content))

	// f, err := exec.Command("ls", "/").Output()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(string(f))

}

func GetMACAddress() string {
	return ExecCmd("LANG=C ifconfig eth0 | awk '/HWaddr/{ print $5 }'")
}

func GetIPAddress() string {
	return ExecCmd("hostname -I | awk '{print $1}'")
}

func GetHostName() string {
	return ExecCmd("cat /proc/sys/kernel/hostname")
}

func GetHostInfo() (hostname string, platform string, kernel string, hostuuid string) {
	n, _ := host.Info()
	hostname = n.Hostname
	platform = fmt.Sprintf("%v-%v", n.Platform, n.PlatformVersion)
	kernel = n.KernelVersion
	hostuuid = n.HostID
	return hostname, platform, kernel, hostuuid
}

func GetCpuInfo() (vendor string, model string, cores string) {
	c, _ := cpu.Info()
	cpunum := len(c)
	cores = fmt.Sprintf("%v*%v", cpunum, c[0].Cores)
	vendor = c[0].VendorID
	model = c[0].ModelName
	// mhz = c[0].Mhz
	// return fmt.Sprintf("%v-%v-%v-%v-%v", corenum, model, cores, mhz, cores)
	return vendor, model, cores
}

func GetCoreNum() string {
	return ExecCmd("cat /proc/cpuinfo | grep processor | wc -l")
}

func GetCpuMHz() string {
	return ExecCmd("cat /proc/cpuinfo | grep 'cpu MHz' | sort -u | awk '{print $NF}' | head -1")
}

func GetMemStat() string {
	// cat /proc/meminfo  | grep MemTotal
	// total       used       free     shared    buffers     cached
	// return ExecCmd("free -m | grep Mem | awk -F':' '{print $NF}'")
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	return fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f", v.Total, v.Free, v.UsedPercent)
}

func GetDiskStat() string {
	// fdisk -l | grep Device -A 1
	return ExecCmd("fdisk -l | grep Disk | grep dev | grep -v docker")
}

//Generate generates a random MAC address or returns an error
func Generate() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return mac
}
