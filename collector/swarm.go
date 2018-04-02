package collector

import (
	"golang.org/x/net/context"
)

type Node struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Addr           string   `json:"addr,omitempty"`
	Containers     string   `json:"containers,omitempty"`
	ReservedCPUs   string   `json:"reserved_cpus,omitempty"`
	ReservedMemory string   `json:"reserved_memory,omitempty"`
	Labels         []string `json:"labels,omitempty"`
}

func (c *Collector) getSwarmNodes() ([]*Node, error) {
	client, err := c.getDockerClient()
	if err != nil {
		return nil, err
	}

	info, err := client.Info(context.Background())
	if err != nil {
		return nil, err
	}

	nodes, err := parseSwarmNodes(info.DriverStatus)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
