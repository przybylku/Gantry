package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunContainer(ctx context.Context, cli *client.Client, imageName string, containerName string, internalPort string) (string, error) {
	
	
	rule := fmt.Sprintf("Host(`%s.localhost`)", containerName)

	containerConfig := &container.Config{
		Image: imageName,
		Labels: map[string]string{
			"managed-by": "gantry",
			"created-for": containerName,
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.%s.rule", containerName):rule,
			fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", containerName): internalPort,
		},
		ExposedPorts: nat.PortSet{
			nat.Port(internalPort + "/tcp"): struct{}{},
		},
	}
	
	

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"gantry-net": {},
		},
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("error with building container %w", err)
	}
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting container %w", err)
	}
	return resp.ID, nil
}