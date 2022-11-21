package models

import "time"

type KubernetesMessage struct {
	Date float64 `json:"date"`
	//Log        time.Time  `json:"log"`
	Log        string     `json:"log"`
	Kubernetes Kubernetes `json:"kubernetes"`
}
type Labels struct {
	Component string `json:"component"`
	Tier      string `json:"tier"`
}
type Annotations struct {
	KubernetesIoConfigHash              string    `json:"kubernetes.io/config.hash"`
	KubernetesIoConfigMirror            string    `json:"kubernetes.io/config.mirror"`
	KubernetesIoConfigSeen              time.Time `json:"kubernetes.io/config.seen"`
	KubernetesIoConfigSource            string    `json:"kubernetes.io/config.source"`
	SeccompSecurityAlphaKubernetesIoPod string    `json:"seccomp.security.alpha.kubernetes.io/pod"`
}
type Kubernetes struct {
	PodName        string      `json:"pod_name"`
	NamespaceName  string      `json:"namespace_name"`
	PodID          string      `json:"pod_id"`
	Labels         Labels      `json:"labels"`
	Annotations    Annotations `json:"annotations"`
	Host           string      `json:"host"`
	ContainerName  string      `json:"container_name"`
	DockerID       string      `json:"docker_id"`
	ContainerHash  string      `json:"container_hash"`
	ContainerImage string      `json:"container_image"`
}
