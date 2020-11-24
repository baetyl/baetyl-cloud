package service

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/context"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var params = map[string]interface{}{
	"Namespace":                  "ns-1",
	"NodeName":                   "node-name-1",
	"CoreAppName":                "core-app-1",
	"CoreConfName":               "core-conf-name-1",
	"CoreConfVersion":            "core-conf-version-1",
	"FunctionAppName":            "func-app-name-1",
	"FunctionConfName":           "func-conf-name-1",
	"FunctionConfVersion":        "func-conf-version-1",
	"EdgeNamespace":              context.EdgeNamespace(),
	"EdgeSystemNamespace":        context.EdgeSystemNamespace(),
	"KubeNodeName":               "kube-node-1",
	"NodeCertName":               "node-cert-name-1",
	"NodeCertVersion":            "node-cert-version-1",
	"NodeCertPem":                "---node cert pem---",
	"NodeCertKey":                "---node cert key---",
	"NodeCertCa":                 "---node cert ca---",
	context.KeyBaetylHostPathLib: "{{." + context.KeyBaetylHostPathLib + "}}",
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
    path: '{{.BAETYL_HOST_PATH_LIB}}/store'
- name: object-download-path
  hostPath:
    path: '{{.BAETYL_HOST_PATH_LIB}}/object'
- name: host-root-path
  hostPath:
    path: '{{.BAETYL_HOST_PATH_LIB}}/host'
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
    server:
      key: var/lib/baetyl/system/certs/key.pem
      cert: var/lib/baetyl/system/certs/crt.pem
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
			want: `kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:aggregated-metrics-reader
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
  - apiGroups: ["metrics.k8s.io"]
    resources: ["pods", "nodes"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: metrics-server:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
  - kind: ServiceAccount
    name: metrics-server
    namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: metrics-server-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
  - kind: ServiceAccount
    name: metrics-server
    namespace: kube-system
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.metrics.k8s.io
spec:
  service:
    name: metrics-server
    namespace: kube-system
  group: metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-server
  namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    k8s-app: metrics-server
spec:
  selector:
    matchLabels:
      k8s-app: metrics-server
  template:
    metadata:
      name: metrics-server
      labels:
        k8s-app: metrics-server
    spec:
      serviceAccountName: metrics-server
      volumes:
        # mount in tmp so we can safely use from-scratch images and/or read-only containers
        - name: tmp-dir
          emptyDir: {}
      containers:
        - name: metrics-server
          image: 'rancher/metrics-server:v0.3.6'
          imagePullPolicy: IfNotPresent
          command:
            - /metrics-server
            - --kubelet-insecure-tls
            - --kubelet-preferred-address-types=InternalDNS,InternalIP,ExternalDNS,ExternalIP,Hostname
          volumeMounts:
            - name: tmp-dir
              mountPath: /tmp
---
apiVersion: v1
kind: Service
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    kubernetes.io/name: "Metrics-server"
    kubernetes.io/cluster-service: "true"
spec:
  selector:
    k8s-app: metrics-server
  ports:
    - port: 443
      protocol: TCP
      targetPort: 443
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:metrics-server
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - nodes
      - nodes/stats
      - namespaces
    verbs:
      - get
      - list
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
  - kind: ServiceAccount
    name: metrics-server
    namespace: kube-system
---
apiVersion: v1
kind: Namespace
metadata:
  name: local-path-storage
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: local-path-provisioner-service-account
  namespace: local-path-storage
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: local-path-provisioner-role
rules:
  - apiGroups: [""]
    resources: ["nodes", "persistentvolumeclaims"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["endpoints", "persistentvolumes", "pods"]
    verbs: ["*"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: local-path-provisioner-bind
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: local-path-provisioner-role
subjects:
  - kind: ServiceAccount
    name: local-path-provisioner-service-account
    namespace: local-path-storage
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: local-path-provisioner
  namespace: local-path-storage
spec:
  replicas: 1
  selector:
    matchLabels:
      app: local-path-provisioner
  template:
    metadata:
      labels:
        app: local-path-provisioner
    spec:
      serviceAccountName: local-path-provisioner-service-account
      containers:
        - name: local-path-provisioner
          image: rancher/local-path-provisioner:v0.0.12
          imagePullPolicy: IfNotPresent
          command:
            - local-path-provisioner
            - --debug
            - start
            - --config
            - /etc/config/config.json
          volumeMounts:
            - name: config-volume
              mountPath: /etc/config/
          env:
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
      volumes:
        - name: config-volume
          configMap:
            name: local-path-config
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-path
provisioner: rancher.io/local-path
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Delete
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: local-path-config
  namespace: local-path-storage
data:
  config.json: |-
    {
            "nodePathMap":[
            {
                    "node":"DEFAULT_PATH_FOR_NON_LISTED_NODES",
                    "paths":["/opt/local-path-provisioner"]
            }]
    }
---
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
      encoding: console

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
              mountPath: /var/lib/baetyl/node
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
			res := string(data)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, res)
		})
	}
}
