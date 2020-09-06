package service

import (
	"fmt"
	"testing"

	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

var params = map[string]interface{}{
	"Namespace":           "ns-1",
	"NodeName":            "node-name-1",
	"CoreAppName":         "core-app-1",
	"CoreConfName":        "core-conf-name-1",
	"CoreConfVersion":     "core-conf-version-1",
	"FunctionAppName":     "func-app-name-1",
	"FunctionConfName":    "func-conf-name-1",
	"FunctionConfVersion": "func-conf-version-1",
	"EdgeNamespace":       common.DefaultBaetylEdgeNamespace,
	"EdgeSystemNamespace": common.DefaultBaetylEdgeSystemNamespace,
	"KubeNodeName":        "kube-node-1",
	"NodeCertName":        "node-cert-name-1",
	"NodeCertVersion":     "node-cert-version-1",
	"NodeCertPem":         "---node cert pem---",
	"NodeCertKey":         "---node cert key---",
	"NodeCertCa":          "---node cert ca---",
}

func TestTemplateServiceImpl_UnmarshalTemplate(t *testing.T) {
	mocks := InitMockEnvironment(t)
	defer mocks.Close()

	funcs := map[string]interface{}{
		"GetProperty": func(in string) string {
			return fmt.Sprintf("out-%s", in)
		},
	}
	sTemplate, err := NewTemplateService(mocks.conf, funcs)

	assert.NoError(t, err)
	assert.NotNil(t, sTemplate)

	tests := []struct {
		name    string
		out     interface{}
		want    string
		wantErr bool
	}{
		{
			name: "baetyl-core-app.yml",
			out:  &v1.Application{},
			want: `name: core-app-1
type: container
labels:
  baetyl-cloud-system: "true"
namespace: ns-1
selector: baetyl-node-name=node-name-1
services:
- name: baetyl-core
  image: out-baetyl-image
  replica: 1
  volumeMounts:
  - name: core-conf
    mountPath: /etc/baetyl
    readOnly: true
  - name: node-cert
    mountPath: /var/lib/baetyl/node
  - name: core-store-path
    mountPath: /var/lib/baetyl/store
  - name: object-download-path
    mountPath: /var/lib/baetyl/object
  - name: host-root-path
    mountPath: /var/lib/baetyl/host
  ports:
  - hostPort: 30050
    containerPort: 80
    protocol: TCP
  args:
  - core
volumes:
- name: core-conf
  config:
    name: core-conf-name-1
    version: core-conf-version-1
- name: node-cert
  secret:
    name: node-cert-name-1
    version: node-cert-version-1
- name: core-store-path
  hostPath:
    path: /var/lib/baetyl/store
- name: object-download-path
  hostPath:
    path: /var/lib/baetyl/object
- name: host-root-path
  hostPath:
    path: /var/lib/baetyl/host
system: true
`,
		},
		{
			name: "baetyl-core-conf.yml",
			out:  &v1.Configuration{},
			want: `name: core-conf-name-1
namespace: ns-1
labels:
  baetyl-app-name: core-app-1
  baetyl-cloud-system: "true"
  baetyl-node-name: node-name-1
data:
  conf.yml: |-
    node:
      ca: var/lib/baetyl/node/ca.pem
      key: var/lib/baetyl/node/client.key
      cert: var/lib/baetyl/node/client.pem
    httplink:
      address: "out-sync-server-address"
      insecureSkipVerify: true
    logger:
      level: debug
system: true
`,
		},
		{
			name: "baetyl-function-app.yml",
			out:  &v1.Application{},
			want: `name: func-app-name-1
type: container
labels:
  baetyl-cloud-system: "true"
namespace: ns-1
selector: baetyl-node-name=node-name-1
services:
- name: baetyl-function
  image: out-baetyl-function-image
  replica: 1
  volumeMounts:
  - name: func-conf
    mountPath: /etc/baetyl
    readOnly: true
  ports:
  - containerPort: 80
    protocol: TCP
volumes:
- name: func-conf
  config:
    name: func-conf-name-1
    version: func-conf-version-1
system: true
`,
		},
		{
			name: "baetyl-function-conf.yml",
			out:  &v1.Configuration{},
			want: `name: func-conf-name-1
namespace: ns-1
labels:
  baetyl-app-name: func-app-name-1
  baetyl-cloud-system: "true"
  baetyl-node-name: node-name-1
data:
  conf.yml: |-
    logger:
      level: debug
system: true
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sTemplate.UnmarshalTemplate(tt.name, params, tt.out)
			assert.NoError(t, err)
			data, err := yaml.Marshal(tt.out)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func TestTemplateServiceImpl_ParseTemplate(t *testing.T) {
	mocks := InitMockEnvironment(t)
	defer mocks.Close()

	funcs := map[string]interface{}{
		"GetProperty": func(in string) string {
			return fmt.Sprintf("out-%s", in)
		},
	}
	sTemplate, err := NewTemplateService(mocks.conf, funcs)

	assert.NoError(t, err)
	assert.NotNil(t, sTemplate)

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "baetyl-init-deployment.yml",
			want: `---
apiVersion: v1
kind: Namespace
metadata:
  name: baetyl-edge-system

---
apiVersion: v1
kind: Namespace
metadata:
  name: baetyl-edge

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: baetyl-edge-system-service-account
  namespace: baetyl-edge-system

---
# elevation of authority
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: baetyl-edge-system-rbac
subjects:
  - kind: ServiceAccount
    name: baetyl-edge-system-service-account
    namespace: baetyl-edge-system
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io

---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer

---
apiVersion: v1
kind: Secret
metadata:
  name: node-cert-name-1
  namespace: baetyl-edge-system
type: Opaque
data:
  client.pem: '---node cert pem---'
  client.key: '---node cert key---'
  ca.pem: '---node cert ca---'

---
# baetyl-init configmap
apiVersion: v1
kind: ConfigMap
metadata:
  name: baetyl-init-config
  namespace: baetyl-edge-system
data:
  conf.yml: |-
    node:
      ca: var/lib/baetyl/node/ca.pem
      key: var/lib/baetyl/node/client.key
      cert: var/lib/baetyl/node/client.pem
    httplink:
      address: out-sync-server-address
      insecureSkipVerify: true
    logger:
      level: debug

---
# baetyl-init deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: baetyl-init
  namespace: baetyl-edge-system
  labels:
    baetyl-app-name: baetyl-init
    baetyl-service-name: baetyl-init
spec:
  selector:
    matchLabels:
      baetyl-service-name: baetyl-init
  replicas: 1
  template:
    metadata:
      labels:
        baetyl-app-name: baetyl-init
        baetyl-service-name: baetyl-init
    spec:
      nodeName: kube-node-1
      serviceAccountName: baetyl-edge-system-service-account
      containers:
        - name: baetyl-init
          image: out-baetyl-image
          imagePullPolicy: IfNotPresent
          args:
            - init
          env:
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - name: init-conf
              mountPath: /etc/baetyl
            - name: core-store-path
              mountPath: /var/lib/baetyl/store
            - name: object-download-path
              mountPath: /var/lib/baetyl/object
            - name: host-root-path
              mountPath: /var/lib/baetyl/host
            - name: node-cert
              mountPath: var/lib/baetyl/node
      volumes:
        - name: init-conf
          configMap:
            name: baetyl-init-config
        - name: core-store-path
          hostPath:
            path: /var/lib/baetyl/store
        - name: object-download-path
          hostPath:
            path: /var/lib/baetyl/object
        - name: host-root-path
          hostPath:
            path: /var/lib/baetyl/host
        - name: node-cert
          secret:
            secretName: node-cert-name-1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := sTemplate.ParseTemplate(tt.name, params)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}
