package containers

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
)

// GetFrontendDeployment returns the front end deployment
func (g *RgpDeployer) GetFrontendDeployment() *components.Deployment {

	deployment := components.NewDeployment(horizonapi.DeploymentConfig{
		Name:      "frontend-service",
		Namespace: g.Grspec.Namespace,
	})

	deployment.AddPod(g.GetFrontendPod())
	deployment.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "frontend-service",
	})

	deployment.AddMatchLabelsSelectors(map[string]string{
		"app":  "rgp",
		"name": "frontend-service",
	})

	return deployment
}

// GetFrontendPod returns the front end pod
func (g *RgpDeployer) GetFrontendPod() *components.Pod {

	pod := components.NewPod(horizonapi.PodConfig{
		Name: "frontend-servicer",
	})

	container, _ := g.GetFrontendContainer()

	pod.AddContainer(container)

	pod.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "frontend-service",
	})

	return pod
}

// GetFrontendContainer will return the container
func (g *RgpDeployer) GetFrontendContainer() (*components.Container, error) {
	container := components.NewContainer(horizonapi.ContainerConfig{
		Name:       "frontend-service",
		Image:      "gcr.io/snps-swip-staging/reporting-frontend-service:0.0.677",
		PullPolicy: horizonapi.PullAlways,
		//MinMem:     "500Mi",
		//MaxMem:     "",
		//MinCPU:     "250m",
		//MaxCPU:     "",
	})

	container.AddPort(horizonapi.PortConfig{
		ContainerPort: "8080",
		Protocol:      horizonapi.ProtocolTCP,
	})

	for _, v := range g.getFrontendEnvConfigs() {
		err := container.AddEnv(*v)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// GetFrontendService returns the front end service
func (g *RgpDeployer) GetFrontendService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "frontend-service",
		Namespace:     g.Grspec.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeDefault,
	})
	service.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "frontend-service",
	})
	service.AddSelectors(map[string]string{
		"name": "frontend-service",
	})
	service.AddPort(horizonapi.ServicePortConfig{Name: "80", Port: 80, Protocol: horizonapi.ProtocolTCP, TargetPort: "8080"})
	return service
}

func (g *RgpDeployer) getFrontendEnvConfigs() []*horizonapi.EnvConfig {
	var envs []*horizonapi.EnvConfig
	envs = append(envs, g.getSwipEnvConfigs()...)
	return envs
}
