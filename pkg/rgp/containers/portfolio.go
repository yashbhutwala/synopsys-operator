package containers

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetPortfolioDeployment returns the portfolio deployment
func (g *RgpDeployer) GetPortfolioDeployment() *components.Deployment {
	deployment := components.NewDeployment(horizonapi.DeploymentConfig{
		Name:      "rp-portfolio-service",
		Namespace: g.Grspec.Namespace,
	})

	deployment.AddPod(g.GetPortfolioPod())
	deployment.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "rp-portfolio-service",
	})

	deployment.AddMatchLabelsSelectors(map[string]string{
		"app":  "rgp",
		"name": "rp-portfolio-service",
	})
	return deployment
}

// GetPortfolioPod returns the portfolio pod
func (g *RgpDeployer) GetPortfolioPod() *components.Pod {

	pod := components.NewPod(horizonapi.PodConfig{
		Name: "rp-portfolio-service",
	})

	container, _ := g.getPortfolioContainer()

	pod.AddContainer(container)
	for _, v := range g.getPortfolioVolumes() {
		pod.AddVolume(v)
	}

	pod.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "rp-portfolio-service",
	})

	return pod
}

func (g *RgpDeployer) getPortfolioContainer() (*components.Container, error) {
	container := components.NewContainer(horizonapi.ContainerConfig{
		Name:       "rp-portfolio-service",
		Image:      "gcr.io/snps-swip-staging/reporting-rp-portfolio-service:latest",
		PullPolicy: horizonapi.PullIfNotPresent,
		//MinMem:     "500Mi",
		//MaxMem:     "",
		//MinCPU:     "250m",
		//MaxCPU:     "",
	})

	container.AddPort(horizonapi.PortConfig{
		ContainerPort: "60289",
		Protocol:      horizonapi.ProtocolTCP,
	})

	for _, v := range g.getPortfolioVolumeMounts() {
		err := container.AddVolumeMount(*v)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range g.getPortfolioEnvConfigs() {
		err := container.AddEnv(*v)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// GetPortfolioService returns the portfolio service
func (g *RgpDeployer) GetPortfolioService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "rp-portfolio-service",
		Namespace:     g.Grspec.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeDefault,
	})
	service.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "rp-portfolio-service",
	})
	service.AddSelectors(map[string]string{
		"name": "rp-portfolio-service",
	})
	service.AddPort(horizonapi.ServicePortConfig{Name: "60289", Port: 60289, Protocol: horizonapi.ProtocolTCP, TargetPort: "60289"})
	service.AddPort(horizonapi.ServicePortConfig{Name: "admin", Port: 8081, Protocol: horizonapi.ProtocolTCP, TargetPort: "8081"})
	return service
}

func (g *RgpDeployer) getPortfolioVolumes() []*components.Volume {
	var volumes []*components.Volume

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-cacert",
		MapOrSecretName: "vault-ca-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_cacrt": {KeyOrPath: "tls.crt", Mode: util.IntToInt32(420)},
		},
	}))

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-client-key",
		MapOrSecretName: "auth-client-tls-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_client_key": {KeyOrPath: "tls.key", Mode: util.IntToInt32(420)},
		},
	}))

	volumes = append(volumes, components.NewSecretVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "vault-client-cert",
		MapOrSecretName: "auth-client-tls-certificate",
		Items: map[string]horizonapi.KeyAndMode{
			"vault_client_cert": {KeyOrPath: "tls.crt", Mode: util.IntToInt32(420)},
		},
	}))

	return volumes
}

func (g *RgpDeployer) getPortfolioVolumeMounts() []*horizonapi.VolumeMountConfig {
	var volumeMounts []*horizonapi.VolumeMountConfig
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-cacert", MountPath: "/mnt/vault/ca"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-client-key", MountPath: "/mnt/vault/key"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-client-cert", MountPath: "/mnt/vault/cert"})

	return volumeMounts
}

func (g *RgpDeployer) getPortfolioEnvConfigs() []*horizonapi.EnvConfig {
	var envs []*horizonapi.EnvConfig
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "SWIP_VAULT_ADDRESS", KeyOrVal: "https://vault:8200"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CACERT", KeyOrVal: "/mnt/vault/ca/vault_cacrt"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_KEY", KeyOrVal: "/mnt/vault/key/vault_client_key"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_CERT", KeyOrVal: "/mnt/vault/cert/vault_client_cert"})

	envs = append(envs, g.getCommonEnvConfigs()...)
	envs = append(envs, g.getSwipEnvConfigs()...)
	envs = append(envs, g.getPostgresEnvConfigs()...)

	return envs
}