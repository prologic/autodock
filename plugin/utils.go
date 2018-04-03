package plugin

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	log "github.com/sirupsen/logrus"
)

func serviceUpdate(client *client.Client, name string, force bool) error {
	service, err := getService(client, name)
	if err != nil {
		log.Errorf("unable to get service %s: %s", name, err)
		return err
	}

	/*
		serviceSpec := swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: name,
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: service.Spec.TaskTemplate.ContainerSpec.Image,
				},
			},
		}
	*/

	if force {
		// oldForceUpdate := service.Spec.TaskTemplate.ForceUpdate
		service.Spec.TaskTemplate.ForceUpdate++
	}

	_, err = client.ServiceUpdate(
		context.Background(),
		service.ID,
		swarm.Version{Index: service.Version.Index},
		//serviceSpec,
		service.Spec,
		types.ServiceUpdateOptions{QueryRegistry: false},
	)

	if err != nil {
		log.Errorf("error updating services %s: %s", name, err)
		return err
	}

	return nil
}

// getService queries the docker API for a service with the name as provided
// serviceName is the desired name
func getService(client *client.Client, name string) (swarm.Service, error) {
	args := filters.NewArgs(filters.Arg("name", name))
	services, err := client.ServiceList(context.Background(), types.ServiceListOptions{Filters: args})

	if services == nil || len(services) == 0 {
		return swarm.Service{}, errors.New("No matching service found for name " + name)
	}

	return services[0], err
}
