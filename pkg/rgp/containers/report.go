package containers

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetReportDeployment returns the report deployment
func (g *RgpDeployer) GetReportDeployment() *components.Deployment {
	deployment := components.NewDeployment(horizonapi.DeploymentConfig{
		Name:      "report-service",
		Namespace: g.Grspec.Namespace,
	})

	deployment.AddPod(g.GetReportPod())
	deployment.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "report-service",
	})

	deployment.AddMatchLabelsSelectors(map[string]string{
		"app":  "rgp",
		"name": "report-service",
	})
	return deployment
}

// GetReportPod returns the report pod
func (g *RgpDeployer) GetReportPod() *components.Pod {

	pod := components.NewPod(horizonapi.PodConfig{
		Name: "report-service",
	})

	container, _ := g.getReportContainer()

	pod.AddContainer(container)
	for _, v := range g.getReportVolumes() {
		pod.AddVolume(v)
	}

	pod.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "report-service",
	})

	return pod
}

// getAuthServersContainer returns the auth server pod
func (g *RgpDeployer) getReportContainer() (*components.Container, error) {
	container := components.NewContainer(horizonapi.ContainerConfig{
		Name:       "report-service",
		Image:      "gcr.io/snps-swip-staging/reporting-report-service:0.0.456",
		PullPolicy: horizonapi.PullAlways,
		//MinMem:     "1Gi",
		//MaxMem:     "",
		//MinCPU:     "250m",
		//MaxCPU:     "",
	})

	container.AddPort(horizonapi.PortConfig{
		ContainerPort: "7979",
		Protocol:      horizonapi.ProtocolTCP,
	})

	for _, v := range g.getReportVolumeMounts() {
		err := container.AddVolumeMount(*v)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range g.getReportEnvConfigs() {
		err := container.AddEnv(*v)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// GetReportService returns the report service
func (g *RgpDeployer) GetReportService() *components.Service {
	service := components.NewService(horizonapi.ServiceConfig{
		Name:          "report-service",
		Namespace:     g.Grspec.Namespace,
		IPServiceType: horizonapi.ClusterIPServiceTypeDefault,
	})
	service.AddLabels(map[string]string{
		"app":  "rgp",
		"name": "report-service",
	})
	service.AddSelectors(map[string]string{
		"name": "report-service",
	})
	service.AddPort(horizonapi.ServicePortConfig{Name: "7979", Port: 7979, Protocol: horizonapi.ProtocolTCP, TargetPort: "7979"})
	service.AddPort(horizonapi.ServicePortConfig{Name: "admin", Port: 8081, Protocol: horizonapi.ProtocolTCP, TargetPort: "8081"})
	return service
}

func (g *RgpDeployer) getReportVolumes() []*components.Volume {
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

func (g *RgpDeployer) getReportVolumeMounts() []*horizonapi.VolumeMountConfig {
	var volumeMounts []*horizonapi.VolumeMountConfig
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-cacert", MountPath: "/mnt/vault/ca"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-client-key", MountPath: "/mnt/vault/key"})
	volumeMounts = append(volumeMounts, &horizonapi.VolumeMountConfig{Name: "vault-client-cert", MountPath: "/mnt/vault/cert"})

	return volumeMounts
}

func (g *RgpDeployer) getReportEnvConfigs() []*horizonapi.EnvConfig {
	var envs []*horizonapi.EnvConfig
	//envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvFromPodIP, NameOrPrefix: "POD_IP"})
	//envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLUSTER_ADDR", KeyOrVal: "https://$(POD_IP):8201"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "SWIP_VAULT_ADDRESS", KeyOrVal: "https://vault:8200"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CACERT", KeyOrVal: "/mnt/vault/ca/vault_cacrt"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_KEY", KeyOrVal: "/mnt/vault/key/vault_client_key"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "VAULT_CLIENT_CERT", KeyOrVal: "/mnt/vault/cert/vault_client_cert"})

	envs = append(envs, g.getCommonEnvConfigs()...)
	envs = append(envs, g.getSwipEnvConfigs()...)
	envs = append(envs, g.getPostgresEnvConfigs()...)

	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "MINIO_HOST", KeyOrVal: "minio"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "MINIO_PORT", KeyOrVal: "9000"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "MINIO_BUCKET", KeyOrVal: "reports"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "MINIO_REGION", KeyOrVal: "us-central1"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvFromSecret, NameOrPrefix: "MINIO_ACCESS_KEY", KeyOrVal: "access_key", FromName: "minio-keys"})
	envs = append(envs, &horizonapi.EnvConfig{Type: horizonapi.EnvFromSecret, NameOrPrefix: "MINIO_SECRET_KEY", KeyOrVal: "secret_key", FromName: "minio-keys"})

	return envs
}
