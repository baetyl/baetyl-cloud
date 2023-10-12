// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=baetyl-cloud/plugin/kube/apis/cloud
// +k8s:defaulter-gen=TypeMeta
// +groupName=cloud.baetyl.io

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: cloud.GroupName, Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Application{},
		&ApplicationList{},
		&Configuration{},
		&ConfigurationList{},
		&Node{},
		&NodeList{},
		&NodeReport{},
		&NodeReportList{},
		&NodeDesire{},
		&NodeDesireList{},
		&Secret{},
		&SecretList{},
	)
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
