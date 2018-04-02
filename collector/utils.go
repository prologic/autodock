package collector

import (
	"strings"

	engineClient "github.com/docker/docker/client"
	"github.com/prologic/autodock/client"
)

func (c *Collector) getDockerURL() string {
	return client.GetDockerURL(c.cfg.DockerURL)
}

func (c *Collector) getDockerClient() (*engineClient.Client, error) {
	return client.GetDockerClient(
		c.cfg.DockerURL,
		c.cfg.TLSCACert,
		c.cfg.TLSCert,
		c.cfg.TLSKey,
		c.cfg.AllowInsecure,
	)
}

// HACK: until we get a consumable endpoint from swarm we must parse the
// node list from /info
func parseSwarmNodes(driverStatus [][2]string) ([]*Node, error) {
	nodes := []*Node{}
	var node *Node
	nodeComplete := false
	name := ""
	addr := ""
	containers := ""
	reservedCPUs := ""
	reservedMemory := ""
	labels := []string{}
	for _, l := range driverStatus {
		if len(l) != 2 {
			continue
		}
		label := l[0]
		data := l[1]

		// cluster info label i.e. "Filters" or "Strategy"
		if strings.Index(label, "\u0008") > -1 {
			continue
		}

		if strings.Index(label, " └") == -1 {
			name = label
			addr = data
		}

		// node info like "Containers"
		switch label {
		case " └ Containers":
			containers = data
		case " └ Reserved CPUs":
			reservedCPUs = data
		case " └ Reserved Memory":
			reservedMemory = data
		case " └ Labels":
			lbls := strings.Split(data, ",")
			labels = lbls
			nodeComplete = true
		default:
			continue
		}

		if nodeComplete {
			node = &Node{
				Name:           name,
				Addr:           addr,
				Containers:     containers,
				ReservedCPUs:   reservedCPUs,
				ReservedMemory: reservedMemory,
				Labels:         labels,
			}
			nodes = append(nodes, node)

			// reset info
			name = ""
			addr = ""
			containers = ""
			reservedCPUs = ""
			reservedMemory = ""
			labels = []string{}
			nodeComplete = false
		}
	}

	return nodes, nil
}
