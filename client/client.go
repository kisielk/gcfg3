package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kisielk/gcfg3/gcfg3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCommand(c string, timeout time.Duration, args ...string) (result string, err error) {
	cmd := exec.Command(c, args...)
	o := make(chan []byte)
	e := make(chan error)
	go func() {
		result, err := cmd.Output()
		if err != nil {
			e <- err
		}
		o <- result
	}()

	select {
	case result := <-o:
		return string(result), nil
	case err := <-e:
		return "", err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", errors.New("process timed out")
	}
	return "", nil
}

func runScript(s string, timeout time.Duration) (result string, err error) {
	if !strings.HasPrefix(s, "#!") {
		return "", errors.New("script buffer must start with a shebang (#!)")
	}

	tempFile, err := ioutil.TempFile("", "gcfg3")
	if err != nil {
		return
	}
	defer os.Remove(tempFile.Name())

	buf := bytes.NewBufferString(s)
	shebang, err := buf.ReadString('\n')
	if err != nil {
		return
	}
	interpreter := strings.TrimSpace(shebang)[2:]

	_, err = buf.WriteTo(tempFile)
	if err != nil {
		return
	}

	return runCommand(interpreter, timeout, tempFile.Name())
}

func processResult(r string) (pd []gcfg3.ProbeData, err error) {
	r = strings.TrimSpace(r)
	lines := strings.Split(r, "\n")
	for _, l := range lines {
		sl := strings.SplitN(l, ":", 2)
		if len(sl) != 2 {
			err = errors.New(fmt.Sprintf("unparseable probe result: %s", l))
			return
		}
		pd = append(pd, gcfg3.ProbeData{sl[0], sl[1]})
	}
	return
}

func main() {
	client, err := gcfg3.NewGcfg3Client("localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	probes, err := client.GetProbes()
	if err != nil {
		log.Fatal("GetProbes:", err)
	}

	var probeData []gcfg3.ProbeData
	for _, probe := range probes {
		r, err := runScript(probe.Script, 10*time.Second)
		if err != nil {
			log.Printf("error in probe '%s': %s\n", probe.Name, err)
			continue
		}

		pd, err := processResult(r)
		if err != nil {
			log.Printf("error processing result from probe '%s': %s\n", probe.Name, err)
			continue
		}

		probeData = append(probeData, pd...)
	}
	client.SendProbeData(probeData)

	entries, err := client.GetEntries()
	for _, entry := range entries {
		log.Printf("%s", entry)
	}
}
