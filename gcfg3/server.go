package gcfg3

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
)

type Empty struct{}

type Probe struct {
	Name   string
	Script string
}

type ProbeData struct {
	Name   string
	Result string
}

type Gcfg3Client struct {
	Client *rpc.Client
}

type ConfigEntry interface {
	Install() error
}

type FileEntry struct {
	Path     string
	Contents string
}

func (e FileEntry) Install() error {
	tempFile, err := ioutil.TempFile("", "gcfg3")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	buf := bytes.NewBufferString(e.Contents)
	_, err = buf.WriteTo(tempFile)
	if err != nil {
		return err
	}
	tempFile.Close()

	err = os.Rename(tempFile.Name(), e.Path)
	if err != nil {
		return err
	}

	return nil
}

func NewGcfg3Client(addr string) (Gcfg3Client, error) {
	client, err := rpc.DialHTTP("tcp", addr)
	return Gcfg3Client{client}, err
}

func (c Gcfg3Client) GetProbes() ([]Probe, error) {
	var probes []Probe
	err := c.Client.Call("Gcfg3Server.GetProbes", Empty{}, &probes)
	return probes, err
}

func (c Gcfg3Client) SendProbeData(d []ProbeData) error {
	var r Empty
	err := c.Client.Call("Gcfg3Server.SendProbeData", d, &r)
	return err
}

func (c Gcfg3Client) GetEntries() ([]ConfigEntry, error) {
	var entries []ConfigEntry
	err := c.Client.Call("Gcfg3Server.GetEntries", Empty{}, &entries)
	return entries, err
}

type Gcfg3Server struct {
	Root string
}

func (s Gcfg3Server) GetProbes(e Empty, p *[]Probe) error {
	probeDir := filepath.Join(s.Root, "probes")
	probeContents, err := ioutil.ReadDir(probeDir)
	if err != nil {
		return err
	}
	for _, entry := range probeContents {
		if entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(probeDir, entry.Name())
		contents, err := ioutil.ReadFile(entryPath)
		if err != nil {
			continue
		}

		*p = append(*p, Probe{entry.Name(), string(contents)})
	}
	return nil
}

func (s Gcfg3Server) GetEntries(n Empty, e *[]ConfigEntry) error {
	file := FileEntry{"/tmp/foo", "contents"}
	*e = append(*e, file)
	return nil
}

func (s Gcfg3Server) SendProbeData(d []ProbeData, e *Empty) error {
	for _, datum := range d {
		log.Printf("%s: %s", datum.Name, datum.Result)
	}
	return nil
}
