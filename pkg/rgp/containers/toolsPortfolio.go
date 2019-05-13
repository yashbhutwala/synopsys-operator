package containers

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetToolsPortfolioDeployment returns the tools portfolio deployment
func (g *RgpDeployer) GetToolsPortfolioDeployment() *components.Deployment {
	deployment := components.NewDeployment(horizonapi.DeploymentConfig{
		Name:      "tools-portfolio-service",
		Namespace: g.Grspec.Namespace,
	})

	deployment.AddPod(g.GetToolsPortfolioPod())
	deployment.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "tools-portfolio-service",
	})

	deployment.AddMatchLabelsSelectors(map[string]string{
		"app":  "rgp",
		"name": "tools-portfolio-service",
	})
	return deployment
}

// GetToolsPortfolioPod returns the tools portfolio pod
func (g *RgpDeployer) GetToolsPortfolioPod() *components.Pod {

	pod := components.NewPod(horizonapi.PodConfig{
		Name: "tools-portfolio-service",
	})

	container, _ := g.getToolPortfolioContainer()

	pod.AddContainer(container)
	for _, v := range g.getToolsPortfolioVolumes() {
		pod.AddVolume(v)
	}

	pod.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "tools-portfolio-service",
	})

	return pod
}

func (g *RgpDeployer) getToolPortfolioContainer() (*components.Container, error) {
	container := components.NewContainer(horizonapi.ContainerConfig{
		Name:       "tools-portfolio-service",
		Image:      "gcr.io/snps-swip-staging/reporting-tools-portfolio-service:0.0.998",
		PullPolicy: horizonapi.PullAlways,
		//MinMem:     "1Gi",
		//MaxMem:     "",
		//MinCPU:     "250m",
		//MaxCPU:     "",
	})

	container.AddPort(horizonapi.PortConfig{
		ContainerPort: "60281",
		Protocol:      horizonapi.ProtocolTCP,
	})

	for _, v := range g.getToolsPortfolioVolumeMounts() {
		err := container.AddVolumeMount(*v)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range g.getToolsPortfolioEnvConfigs() {
		err := container.AddEnv(*v)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// GetToolsPortfolioService returns the tools portfolio service
func (g *RgpDeployer) GetToolsPortfolioService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "tools-portfolio-service",
		Namespace:     g.Grspec.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeDefault,
	})
	service.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "tools-portfolio-service",
	})
	service.AddSelectors(map[string]string{
		"name": "tools-portfolio-service",
	})
	service.AddPort(horizonapi.ServicePortConfig{Name: "60289", Port: 60281, Protocol: horizonapi.ProtocolTCP, TargetPort: "60281"})
	service.AddPort(horizonapi.ServicePortConfig{Name: "admin", Port: 8081, Protocol: horizonapi.ProtocolTCP, TargetPort: "8081"})
	return service
}

func (g *RgpDeployer) getToolsPortfolioVolumes() []*components.Volume {
	var volumes []*components.Volume

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-cacert",
		MapOrSecretName: "vault-ca-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_cacrt": {KeyOrPath: "tls.crt", Mode: util.IntToInt32(420)},
		},
	}))

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-server-key",
		MapOrSecretName: "auth-server-tls-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_server_key": {KeyOrPath: "tls.key", Mode: util.IntToInt32(420)},
		},
	}))

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-server-cert",
		MapOrSecretName: "auth-server-tls-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_server_cert": {KeyOrPath: "tls.crt", Mode: util.IntToInt32(420)},
		},
	}))

	return volumes
}

func (g *RgpDeployer) getToolsPortfolioVolumeMounts() []*horizonapi.VolumeMountConfig {
	var volumeMounts []*horizonapi.VolumeMountConfig
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-cacert", MountPath: "/mnt/vault/ca"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-server-key", MountPath: "/mnt/vault/key"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-server-cert", MountPath: "/mnt/vault/cert"})

	return volumeMounts
}

func (g *RgpDeployer) getToolsPortfolioEnvConfigs() []*horizonapi.EnvConfig {
	var envs []*horizonapi.EnvConfig
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "SWIP_VAULT_ADDRESS", KeyOrVal: "https://vault:8200"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CACERT", KeyOrVal: "/mnt/vault/ca/vault_cacrt"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_KEY", KeyOrVal: "/mnt/vault/key/vault_server_key"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_CERT", KeyOrVal: "/mnt/vault/cert/vault_server_cert"})

	envs = append(envs, g.getCommonEnvConfigs()...)
	envs = append(envs, g.getSwipEnvConfigs()...)
	envs = append(envs, g.getPostgresEnvConfigs()...)

	return envs
}
