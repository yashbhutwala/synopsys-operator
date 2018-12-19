/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownershia. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package sample-component

import (
	"fmt"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
)

// sampleComponentDeployment creates a new deployment for sample-component
func (a *SpecConfig) sampleComponentDeployment() (*components.Deployment, error) {
	replicas := int32(1)
	deployment := components.NewDeployment(horizonapi.DeploymentConfig{
		Replicas:  &replicas,
		Name:      "sample-component",
		Namespace: a.config.Namespace,
	})
	deployment.AddMatchLabelsSelectors(map[string]string{"app": "sample-component", "tier": "sample-component"})

	pod, err := a.sampleComponentPod()
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %v", err)
	}
	deployment.AddPod(pod)

	return deployment, nil
}

func (a *SpecConfig) sampleComponentPod() (*components.Pod, error) {
	pod := components.NewPod(horizonapi.PodConfig{
		Name: "sample-component",
	})
	pod.AddLabels(map[string]string{"app": "sample-component", "tier": "sample-component"})

	pod.AddContainer(a.sampleComponentContainer())

	vol, err := a.sampleComponentVolume()
	if err != nil {
		return nil, fmt.Errorf("error creating volumes: %v", err)
	}
	pod.AddVolume(vol)

	return pod, nil
}

func (a *SpecConfig) sampleComponentContainer() *components.Container {
	// This will prevent it from working on openshift without a privileged service account.  Remove once the
	// chowns are removed in the image
	user := int64(0)
	container := components.NewContainer(horizonapi.ContainerConfig{
		Name:   "sample-component",
		Image:  fmt.Sprintf("%s/%s/%s:%s", a.config.Registry, a.config.ImagePath, a.config.SampleComponentImageName, a.config.SampleComponentImageVersion),
		MinMem: a.config.SampleComponentMemory,
		UID:    &user,
	})

	container.AddPort(horizonapi.PortConfig{
		ContainerPort: "8443",
		Protocol:      horizonapi.ProtocolTCP,
	})

	container.AddVolumeMount(horizonapi.VolumeMountConfig{
		Name:      "dir-sample-component",
		MountPath: "/opt/blackduck/sample-component/sample-component-config",
	})

	container.AddEnv(horizonapi.EnvConfig{
		Type:     horizonapi.EnvFromConfigMap,
		FromName: "sample-component",
	})

	container.AddLivenessProbe(horizonapi.ProbeConfig{
		ActionConfig: horizonapi.ActionConfig{
			Command: []string{"/usr/local/bin/docker-healthcheck.sh", "https://localhost:8443/sample-component/api/about"},
		},
		Delay:           240,
		Timeout:         10,
		Interval:        30,
		MinCountFailure: 5,
	})

	return container
}

func (a *SpecConfig) sampleComponentVolume() (*components.Volume, error) {
	vol, err := components.NewEmptyDirVolume(horizonapi.EmptyDirVolumeConfig{
		VolumeName: "dir-sample-component",
		Medium:     horizonapi.StorageMediumDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create empty dir volume: %v", err)
	}

	return vol, nil
}

// sampleComponentService creates a service for sample-component
func (a *SpecConfig) sampleComponentService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "sample-component",
		Namespace:     a.config.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeNodePort,
	})

	service.AddPort(horizonapi.ServicePortConfig{
		Port:       8443,
		TargetPort: "8443",
		Protocol:   horizonapi.ProtocolTCP,
		Name:       "8443-tcp",
	})

	service.AddSelectors(map[string]string{"app": "sample-component"})

	return service
}

// sampleComponentExposedService creates a loadBalancer service for sample-component
func (a *SpecConfig) sampleComponentExposedService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "sample-component-lb",
		Namespace:     a.config.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeLoadBalancer,
	})

	service.AddPort(horizonapi.ServicePortConfig{
		Port:       8443,
		TargetPort: "8443",
		Protocol:   horizonapi.ProtocolTCP,
		Name:       "8443-tcp",
	})

	service.AddSelectors(map[string]string{"app": "sample-component"})

	return service
}

// sampleComponentConfigMap creates a config map for sample-component
func (a *SpecConfig) sampleComponentConfigMap() *components.ConfigMap {
	configMap := components.NewConfigMap(horizonapi.ConfigMapConfig{
		Name:      "sample-component",
		Namespace: a.config.Namespace,
	})

	configMap.AddData(map[string]string{
		"SAMPLE_COMPONENT_SERVER_PORT":         fmt.Sprintf("%d", *a.config.Port),
	})

	return configMap
}
