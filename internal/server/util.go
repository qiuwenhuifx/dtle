package server

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/hashicorp/serf/serf"
)

// ensurePath is used to make sure a path exists
func ensurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// RuntimeStats is used to return various runtime information
func RuntimeStats() map[string]string {
	return map[string]string{
		"kernel.name": runtime.GOOS,
		"arch":        runtime.GOARCH,
		"version":     runtime.Version(),
		"max_procs":   strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
		"goroutines":  strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		"cpu_count":   strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}
}

// serverParts is used to return the parts of a server role
type serverParts struct {
	Name       string
	Region     string
	Datacenter string
	Port       int
	Bootstrap  bool
	Expect     int
	Addr       net.Addr
}

func (s *serverParts) String() string {
	return fmt.Sprintf("%s (Addr: %s) (DC: %s)",
		s.Name, s.Addr, s.Datacenter)
}

// Returns if a member is a Udup server. Returns a boolean,
// and a struct with the various important components
func isUdupServer(m serf.Member) (bool, *serverParts) {
	if m.Tags["role"] != "server" {
		return false, nil
	}

	region := m.Tags["region"]
	datacenter := m.Tags["dc"]
	_, bootstrap := m.Tags["bootstrap"]

	expect := 0
	expect_str, ok := m.Tags["expect"]
	var err error
	if ok {
		expect, err = strconv.Atoi(expect_str)
		if err != nil {
			return false, nil
		}
	}

	port_str := m.Tags["port"]
	port, err := strconv.Atoi(port_str)
	if err != nil {
		return false, nil
	}

	addr := &net.TCPAddr{IP: m.Addr, Port: port}
	parts := &serverParts{
		Name:       m.Name,
		Region:     region,
		Datacenter: datacenter,
		Port:       port,
		Bootstrap:  bootstrap,
		Expect:     expect,
		Addr:       addr,
	}
	return true, parts
}

// shuffleStrings randomly shuffles the list of strings
func shuffleStrings(list []string) {
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

// maxUint64 returns the maximum value
func maxUint64(a, b uint64) uint64 {
	if a >= b {
		return a
	}
	return b
}