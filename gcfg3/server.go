package gcfg3

import (
	"log"
	"net/rpc"
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

func NewGcfg3Client(addr string) (Gcfg3Client, error) {
	client, err := rpc.DialHTTP("tcp", addr)
	return Gcfg3Client{client}, err
}

func (c Gcfg3Client) GetProbes() ([]Probe, error) {
	var probes []Probe
	err := c.Client.Call("Gcfg3Server.GetProbes", Empty{}, &probes)
	return probes, err
}

func (c Gcfg3Client) SendProbeData(d []Probe) error {
	err := c.Client.Call("Gcfg3Server.SendProbeData", d, Empty{})
	return err
}

type Gcfg3Server struct{}

func (s Gcfg3Server) GetProbes(e Empty, p *[]Probe) error {
	hostname := Probe{"hostname", "#!/bin/bash\nhostname\n"}
	*p = append(*p, hostname)
	return nil
}

func (s Gcfg3Server) SendProbeData(d []ProbeData, e *Empty) error {
	for _, datum := range d {
		log.Printf("%s: %s", datum.Name, datum.Result)
	}
	return nil
}
