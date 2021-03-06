/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

import (
	"errors"
	"fmt"
)

// Container is the definition for a container type in marathon
type Container struct {
	Type    string    `json:"type,omitempty"`
	Docker  *Docker   `json:"docker,omitempty"`
	Volumes *[]Volume `json:"volumes,omitempty"`
}

// PortMapping is the portmapping structure between container and mesos
type PortMapping struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort"`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol"`
}

// Parameters is the parameters to pass to the docker client when creating the container
type Parameters struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Volume is the docker volume details associated to the container
type Volume struct {
	ContainerPath string `json:"containerPath,omitempty"`
	HostPath      string `json:"hostPath,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

// Docker is the docker definition from a marathon application
type Docker struct {
	ForcePullImage *bool          `json:"forcePullImage,omitempty"`
	Image          string         `json:"image,omitempty"`
	Network        string         `json:"network,omitempty"`
	Parameters     *[]Parameters  `json:"parameters,omitempty"`
	PortMappings   *[]PortMapping `json:"portMappings,omitempty"`
	Privileged     *bool          `json:"privileged,omitempty"`
}

// Volume attachs a volume to the container
//		host_path:			the path on the docker host to map
//		container_path:		the path inside the container to map the host volume
//		mode:				the mode to map the container
func (container *Container) Volume(hostPath, containerPath, mode string) *Container {
	if container.Volumes == nil {
		container.EmptyVolumes()
	}

	volumes := *container.Volumes
	volumes = append(volumes, Volume{
		ContainerPath: containerPath,
		HostPath:      hostPath,
		Mode:          mode,
	})

	container.Volumes = &volumes

	return container
}

// EmptyVolumes explicitly empties the volumes -- use this if you need to empty
// volumes of an application that already has volumes set (setting volumes to nil will
// keep the current value)
func (container *Container) EmptyVolumes() *Container {
	container.Volumes = &[]Volume{}
	return container
}

// NewDockerContainer creates a default docker container for you
func NewDockerContainer() *Container {
	container := &Container{}
	container.Type = "DOCKER"
	container.Docker = &Docker{}

	return container
}

// SetForcePullImage sets whether the docker image should always be force pulled before
// starting an instance
//		forcePull:			true / false
func (docker *Docker) SetForcePullImage(forcePull bool) *Docker {
	docker.ForcePullImage = &forcePull

	return docker
}

// SetPrivileged sets whether the docker image should be started
// with privilege turned on
//		priv:			true / false
func (docker *Docker) SetPrivileged(priv bool) *Docker {
	docker.Privileged = &priv

	return docker
}

// Container sets the image of the container
//		image:			the image name you are using
func (docker *Docker) Container(image string) *Docker {
	docker.Image = image
	return docker
}

// Bridged sets the networking mode to bridged
func (docker *Docker) Bridged() *Docker {
	docker.Network = "HOST"
	return docker
}

// Expose sets the container to expose the following TCP ports
//		ports:			the TCP ports the container is exposing
func (docker *Docker) Expose(ports ...int) *Docker {
	for _, port := range ports {
		docker.ExposePort(port, 0, 0, "tcp")
	}
	return docker
}

// ExposeUDP sets the container to expose the following UDP ports
//		ports:			the UDP ports the container is exposing
func (docker *Docker) ExposeUDP(ports ...int) *Docker {
	for _, port := range ports {
		docker.ExposePort(port, 0, 0, "udp")
	}
	return docker
}

// ExposePort exposes an port in the container
//		containerPort:			the container port which is being exposed
//		hostPort:						the host port we should expose it on
//		servicePort:				check the marathon documentation
//		protocol:						the protocol to use TCP, UDP
func (docker *Docker) ExposePort(containerPort, hostPort, servicePort int, protocol string) *Docker {
	if docker.PortMappings == nil {
		docker.EmptyPortMappings()
	}

	portMappings := *docker.PortMappings
	portMappings = append(portMappings, PortMapping{
		ContainerPort: containerPort,
		HostPort:      hostPort,
		ServicePort:   servicePort,
		Protocol:      protocol})
	docker.PortMappings = &portMappings

	return docker
}

// EmptyPortMappings explicitly empties the port mappings -- use this if you need to empty
// port mappings of an application that already has port mappings set (setting port mappings to nil will
// keep the current value)
func (docker *Docker) EmptyPortMappings() *Docker {
	docker.PortMappings = &[]PortMapping{}
	return docker
}

// AddParameter adds a parameter to the docker execution line when creating the container
//		key:			the name of the option to add
//		value:		the value of the option
func (docker *Docker) AddParameter(key string, value string) *Docker {
	if docker.Parameters == nil {
		docker.EmptyParameters()
	}

	parameters := *docker.Parameters
	parameters = append(parameters, Parameters{
		Key:   key,
		Value: value})

	docker.Parameters = &parameters

	return docker
}

// EmptyParameters explicitly empties the parameters -- use this if you need to empty
// parameters of an application that already has parameters set (setting parameters to nil will
// keep the current value)
func (docker *Docker) EmptyParameters() *Docker {
	docker.Parameters = &[]Parameters{}
	return docker
}

// ServicePortIndex finds the service port index of the exposed port
//		port:			the port you are looking for
func (docker *Docker) ServicePortIndex(port int) (int, error) {
	if docker.PortMappings == nil || len(*docker.PortMappings) == 0 {
		return 0, errors.New("The docker does not contain any port mappings to search")
	}

	// step: iterate and find the port
	for index, containerPort := range *docker.PortMappings {
		if containerPort.ContainerPort == port {
			return index, nil
		}
	}

	// step: we didn't find the port in the mappings
	return 0, fmt.Errorf("The container port required was not found in the container port mappings")
}
