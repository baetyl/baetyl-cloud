package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Node customize node resource definition
type Node struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NodeSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Node `json:"items"`
}

type NodeSpec struct {
	Attributes map[string]interface{}   `json:"attributes,omitempty"`
	DesireRef  *v1.LocalObjectReference `json:"desireRef,omitempty"`
	ReportRef  *v1.LocalObjectReference `json:"reportRef,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            ReportStatus `json:"status,omitempty"`
}

type ReportStatus struct {
	Report []byte `json:"report,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeReport `json:"items"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeDesire struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DesireSpec `json:"spec,omitempty"`
}

type DesireSpec struct {
	Desire []byte `json:"desire,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeDesireList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeDesire `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application application resource definition
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

// ApplicationSpec application specification
type ApplicationSpec struct {
	// Selector label selector of application
	Selector string `json:"selector,omitempty"`
	// specifies the service information of the application
	Services []Service `json:"services,omitempty"`
	// specifies the storage volume information of the application
	Volumes []Volume `json:"volumes,omitempty"`
	// the registry's credentials
	Registries []*ObjectReference `json:"registries,omitempty"`
	// specifies the system attr of the application
	System bool `json:"system,omitempty"`
	// specifies the type of the application
	Type string `json:"type,omitempty"`
}

type Service struct {
	v1.Container
	// specifies the hostname of the service
	Hostname string `json:"hostname,omitempty"`
	// specifies the device of the service
	Devices []Device `json:"devices,omitempty"`
	// specifies the number of instances started
	Replica int `json:"replica,omitempty"`
	// specifies resource limits for a single instance of the service,  only for Docker container mode
	Resources *Resources `json:"resources,omitempty"`
	// specifies runtime to use, only for Docker container mode
	Runtime string `json:"runtime,omitempty"`
	// labels
	Labels map[string]string `json:"labels,omitempty"`
	// specifies the security context of service
	SecurityContext *SecurityContext `json:"security,omitempty"`
	// specifies host network mode of service
	HostNetwork bool `json:"hostNetwork,omitempty"`
	// specifies function config of service
	FunctionConfig *ServiceFunctionConfig `json:"functionConfig,omitempty"`
}

type SecurityContext struct {
	Privileged bool `json:"privileged,omitempty"`
}

// VolumeDevice device volume config
type Device struct {
	DevicePath  string `json:"devicePath,omitempty"`
	Policy      string `json:"policy,omitempty"`
	Description string `json:"description,omitempty"`
}

type Volume struct {
	Name     string                `json:"name,omitempty"`
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty"`
	Config   *ObjectReference      `json:"config,omitempty"`
	Secret   *ObjectReference      `json:"secret,omitempty"`
}

// ObjectReference the reference of other source
type ObjectReference struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type HostPathVolumeSource struct {
	Path string `json:"path,omitempty"`
}

// Resources resources config
type Resources struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

type ServiceFunctionConfig struct {
	Name    string `json:"name,omitempty" validate:"resourceName"`
	Runtime string `json:"runtime,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application application resource definition
type Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Data              map[string][]byte `json:"data,omitempty"`
	System            bool              `json:"system,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secret `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Configuration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Data              map[string]string `json:"data,omitempty"`
	System            bool              `json:"system,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Configuration `json:"items"`
}
