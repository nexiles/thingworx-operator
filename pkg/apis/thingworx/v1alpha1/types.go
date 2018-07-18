package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"github.com/sirupsen/logrus"
)


const (
	defaultBaseImage = "nexiles/thingworx"
	// version format is "<our-version>-<upstream-version>"
	defaultVersion = "0.3.0-8.3.0"
)

type ClusterPhase string

const (
	ClusterPhaseInitial ClusterPhase = ""
	ClusterPhaseRunning              = "Running"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ThingworxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Thingworx `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Thingworx struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ThingworxSpec   `json:"spec"`
	Status            ThingworxStatus `json:"status,omitempty"`
}

type ThingworxSpec struct {
	// Number of nodes to deploy for a Vault deployment.
	// Default: 1.
	Nodes int32 `json:"nodes,omitempty"`

	// Base image to use for a Vault deployment.
	BaseImage string `json:"baseImage"`

	// Version of Vault to be deployed.
	Version string `json:"version"`

	// Pod defines the policy for pods owned by vault operator.
	// This field cannot be updated once the CR is created.
	Pod *PodPolicy `json:"pod,omitempty"`

	// Name of the ConfigMap for Vault's configuration
	// If this is empty, operator will create a default config for Vault.
	// If this is not empty, operator will create a new config overwriting
	// the "storage", "listener" sections in orignal config.
	ConfigMapName string `json:"configMapName"`
}

func (twx *Thingworx) SetDefaults() (changed bool) {
	logrus.Debugf("SetDefault(): %s", twx.GetName())
	changed = false

	if twx.Spec.Nodes == 0 {
		twx.Spec.Nodes = 1
		changed = true
	}

	if len(twx.Spec.BaseImage) == 0 {
		twx.Spec.BaseImage = defaultBaseImage
		changed = true
	}

	if len(twx.Spec.Version) == 0 {
		twx.Spec.Version = defaultVersion
		changed = true
	}
	return
}

// PodPolicy defines the policy for pods owned by vault operator.
type PodPolicy struct {
	// Resources is the resource requirements for the containers.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

type ThingworxStatus struct {
	// Phase indicates the state this Thingworx cluster jumps in.
	// Phase goes as one way as below:
	//   Initial -> Running
	Phase ClusterPhase `json:"phase"`

	// Initialized indicates if the Thingworx service is initialized.
	Initialized bool `json:"initialized"`

	// ServiceName is the LB service for accessing vault nodes.
	ServiceName string `json:"serviceName,omitempty"`

	// ClientPort is the port for vault client to access.
	// It's the same on client LB service and vault nodes.
	ClientPort int `json:"clientPort,omitempty"`

	// PodNames of updated Thingworx nodes. Updated means the Thingworx container image version
	// matches the spec's version.
	UpdatedNodes []string `json:"updatedNodes,omitempty"`
}
