package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	mf "github.com/baetyl/baetyl-cloud/v2/mock/facade"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	gmodels "github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
	"github.com/baetyl/baetyl-go/v2/context"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	testAppDeploy = `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:latest
          name: nginx
          ports:
          - containerPort: 80
          volumeMounts:
          - name: common-cm
            mountPath: /etc/config
          - name: dcell
            mountPath: /etc/secret
          - name: cache-volume
            mountPath: /cache
          - name: test-volume
            mountPath: /test-hp
      imagePullSecrets:
        - name: myregistrykey
      volumes:
        - name: common-cm
          configMap:
            name: common-cm
        - name: dcell
          secret:
            secretName: dcell
        - name: cache-volume
          emptyDir: {}
        - name: test-volume
          hostPath:
            path: /var/lib/baetyl
            type: Directory`
	updateAppDeploy = `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:latest
          name: nginx
          ports:
          - containerPort: 8080
          volumeMounts:
          - name: common-cm
            mountPath: /etc/config
          - name: dcell
            mountPath: /etc/secret
          - name: cache-volume
            mountPath: /cache
          - name: test-volume
            mountPath: /test-hp
      imagePullSecrets:
        - name: myregistrykey
      volumes:
        - name: common-cm
          configMap:
            name: common-cm
        - name: dcell
          secret:
            secretName: dcell
        - name: cache-volume
          emptyDir: {}
        - name: test-volume
          hostPath:
            path: /var/lib/baetyl
            type: Directory`
	testAppDs = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    app: nginx
  name: nginx
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:latest
          name: nginx
          resources:
            limits:
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
          ports:
          - containerPort: 80
          volumeMounts:
          - name: common-cm
            mountPath: /etc/config
          - name: dcell
            mountPath: /etc/secret
          - name: cache-volume
            mountPath: /cache
          - name: test-volume
            mountPath: /test-hp
      imagePullSecrets:
        - name: myregistrykey
      volumes:
        - name: common-cm
          configMap:
            name: common-cm
        - name: dcell
          secret:
            secretName: dcell
        - name: cache-volume
          emptyDir: {}
        - name: test-volume
          hostPath:
            path: /var/lib/baetyl
            type: Directory`
	updateAppDs = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    app: nginx
  name: nginx
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:latest
          name: nginx
          resources:
            limits:
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
          ports:
          - containerPort: 8080
          volumeMounts:
          - name: common-cm
            mountPath: /etc/config
          - name: dcell
            mountPath: /etc/secret
          - name: cache-volume
            mountPath: /cache
          - name: test-volume
            mountPath: /test-hp
      imagePullSecrets:
        - name: myregistrykey
      volumes:
        - name: common-cm
          configMap:
            name: common-cm
        - name: dcell
          secret:
            secretName: dcell
        - name: cache-volume
          emptyDir: {}
        - name: test-volume
          hostPath:
            path: /var/lib/baetyl
            type: Directory`
	testAppJob = `
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
  labels:
    app: pi
spec:
  backoffLimit: 6
  completions: 1
  parallelism: 1
  template:
    metadata:
      name: pi
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never`
	updateAppJob = `
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
  labels:
    app: pi
spec:
  backoffLimit: 6
  completions: 1
  parallelism: 1
  template:
    metadata:
      name: pi
    spec:
      containers:
      - name: pi
        image: perl:latest
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never`
	testService = `
apiVersion: v1
kind: Service
metadata:
  labels:
    svc: nginx
  name: nginx-svc
  namespace: default
spec:
  ports:
    - name: web
      port: 80
      targetPort: 80
      nodePort: 8080
  selector:
    app: nginx
  type: NodePort`
	updateService = `
apiVersion: v1
kind: Service
metadata:
  labels:
    svc: nginx
  name: nginx-svc
  namespace: default
spec:
  ports:
    - name: web
      port: 80
      targetPort: 80
      nodePort: 9090
  selector:
    app: nginx
  type: NodePort`
)

func TestAPI_CreateDeployApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create deploy app
	deployApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}
	cfg := &specV1.Configuration{
		Name: "common-cm",
	}
	secret := &specV1.Secret{
		Name: "dcell",
	}
	registry := &specV1.Secret{
		Name: "myregistrykey",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	sApp.EXPECT().Get("default", "nginx", "").Return(nil, nil).Times(1)
	sConfig.EXPECT().Get("default", "common-cm", "").Return(cfg, nil).Times(2)
	sSecret.EXPECT().Get("default", "dcell", "").Return(secret, nil).Times(4)
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(registry, nil).Times(4)
	sFacade.EXPECT().CreateApp("default", nil, gomock.Any(), nil).Return(deployApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(testAppDeploy))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal(body, &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(testAppDeploy))
	deploy, _ := resources[0].(*appv1.Deployment)
	resapp, err := api.generateDeployApp("default", deploy)
	resapp.Services[0].Resources.Limits = nil
	resapp.Services[0].Resources.Requests = nil
	appView, _ := api.ToApplicationView(resapp)
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_CreateDaemonSetApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create ds app
	dsApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{
					Limits:   map[string]string{"memory": "200Mi"},
					Requests: map[string]string{"cpu": "100m", "memory": "200Mi"},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Workload: "daemonset",
	}
	cfg := &specV1.Configuration{
		Name: "common-cm",
	}
	secret := &specV1.Secret{
		Name: "dcell",
	}
	registry := &specV1.Secret{
		Name: "myregistrykey",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	sApp.EXPECT().Get("default", "nginx", "").Return(nil, nil).Times(1)
	sConfig.EXPECT().Get("default", "common-cm", "").Return(cfg, nil).Times(2)
	sSecret.EXPECT().Get("default", "dcell", "").Return(secret, nil).Times(4)
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(registry, nil).Times(4)
	sFacade.EXPECT().CreateApp("default", nil, gomock.Any(), nil).Return(dsApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(testAppDs))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal([]byte(body), &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(testAppDs))
	ds, _ := resources[0].(*appv1.DaemonSet)
	resapp, err := api.generateDaemonSetApp("default", ds)
	appView, _ := api.ToApplicationView(resapp)
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_CreateJobApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create ds app
	jobApp := &specV1.Application{
		Name:      "pi",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "pi",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
		JobConfig: &specV1.AppJobConfig{
			RestartPolicy: "Never",
			BackoffLimit:  6,
			Parallelism:   1,
			Completions:   1,
		},
		Workload: "job",
	}

	sApp.EXPECT().Get("default", "pi", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateApp("default", nil, gomock.Any(), nil).Return(jobApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(testAppJob))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal([]byte(body), &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(testAppJob))
	job, _ := resources[0].(*batchv1.Job)
	resapp, err := api.generateJobApp("default", job)
	resapp.Services[0].Resources = nil
	appView, _ := api.ToApplicationView(resapp)
	appView.Volumes = nil
	appView.Registries = nil
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_UpdateDeployApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: update common kv config
	expectApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}
	updateApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 8080,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}
	cfg := &specV1.Configuration{
		Name: "common-cm",
	}
	secret := &specV1.Secret{
		Name: "dcell",
	}
	registry := &specV1.Secret{
		Name: "myregistrykey",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}

	sApp.EXPECT().Get("default", "nginx", "").Return(expectApp, nil).Times(1)
	sConfig.EXPECT().Get("default", "common-cm", "").Return(cfg, nil).Times(2)
	sSecret.EXPECT().Get("default", "dcell", "").Return(secret, nil).Times(4)
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(registry, nil).Times(4)
	sFacade.EXPECT().UpdateApp("default", expectApp, gomock.Any(), nil).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(updateAppDeploy))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal(body, &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(updateAppDeploy))
	deploy, _ := resources[0].(*appv1.Deployment)
	resapp, err := api.generateDeployApp("default", deploy)
	resapp.Services[0].Resources.Limits = nil
	resapp.Services[0].Resources.Requests = nil
	appView, _ := api.ToApplicationView(resapp)
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_UpdateDsApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: update common kv config
	dsApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{
					Limits:   map[string]string{"memory": "200Mi"},
					Requests: map[string]string{"cpu": "100m", "memory": "200Mi"},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Workload: "daemonset",
	}
	updateApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 8080,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{
					Limits:   map[string]string{"memory": "200Mi"},
					Requests: map[string]string{"cpu": "100m", "memory": "200Mi"},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
			{
				Name: "myregistrykey",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "myregistrykey",
					},
				},
			},
		},
		Workload: "daemonset",
	}
	cfg := &specV1.Configuration{
		Name: "common-cm",
	}
	secret := &specV1.Secret{
		Name: "dcell",
	}
	registry := &specV1.Secret{
		Name: "myregistrykey",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}

	sApp.EXPECT().Get("default", "nginx", "").Return(dsApp, nil).Times(1)
	sConfig.EXPECT().Get("default", "common-cm", "").Return(cfg, nil).Times(2)
	sSecret.EXPECT().Get("default", "dcell", "").Return(secret, nil).Times(4)
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(registry, nil).Times(4)
	sFacade.EXPECT().UpdateApp("default", dsApp, gomock.Any(), nil).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(updateAppDs))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal(body, &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(updateAppDs))
	ds, _ := resources[0].(*appv1.DaemonSet)
	resapp, err := api.generateDaemonSetApp("default", ds)
	appView, _ := api.ToApplicationView(resapp)
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_UpdateJobApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create ds app
	jobApp := &specV1.Application{
		Name:      "pi",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "pi",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
		JobConfig: &specV1.AppJobConfig{
			RestartPolicy: "Never",
			BackoffLimit:  6,
			Parallelism:   1,
			Completions:   1,
		},
		Workload: "job",
	}
	updateApp := &specV1.Application{
		Name:      "pi",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "pi",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:    "pi",
				Image:   "perl:latest",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
		JobConfig: &specV1.AppJobConfig{
			RestartPolicy: "Never",
			BackoffLimit:  6,
			Parallelism:   1,
			Completions:   1,
		},
		Workload: "job",
	}

	sApp.EXPECT().Get("default", "pi", "").Return(jobApp, nil).Times(1)
	sFacade.EXPECT().UpdateApp("default", jobApp, gomock.Any(), nil).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(updateAppJob))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	var resApp []gmodels.ApplicationView
	body, err := ioutil.ReadAll(re.Body)
	assert.NilError(t, err)
	err = json.Unmarshal(body, &resApp)
	assert.NilError(t, err)

	resources := api.parseK8SYaml([]byte(updateAppJob))
	job, _ := resources[0].(*batchv1.Job)
	resapp, err := api.generateJobApp("default", job)
	resapp.Services[0].Resources = nil
	appView, _ := api.ToApplicationView(resapp)
	appView.Registries = nil
	appView.Volumes = nil
	assert.DeepEqual(t, &resApp[0], appView)
}

func TestAPI_DeleteApp(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Index = sIndex
	api.Facade = sFacade

	// good case: delete kv config
	deployApp := &specV1.Application{
		Name:      "deployApp",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}
	sApp.EXPECT().Get("default", "nginx", "").Return(deployApp, nil).Times(1)
	sFacade.EXPECT().DeleteApp("default", "deployApp", deployApp).Return(nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "app.yaml")
	io.Copy(fw, strings.NewReader(testAppDeploy))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml/delete", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_CreateService(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create deploy app
	appList := &gmodels.ApplicationList{
		Total: 1,
		Items: []gmodels.AppItem{
			{
				Name: "nginx",
			},
		},
	}
	deployApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	updateApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						NodePort:      8080,
						ServiceType:   string(corev1.ServiceTypeNodePort),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	sApp.EXPECT().List("default", gomock.Any()).Return(appList, nil).Times(1)
	sApp.EXPECT().Get("default", "nginx", "").Return(deployApp, nil).Times(1)
	sApp.EXPECT().Update(nil, "default", updateApp).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "service.yaml")
	io.Copy(fw, strings.NewReader(testService))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateService(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create deploy app
	appList := &gmodels.ApplicationList{
		Total: 1,
		Items: []gmodels.AppItem{
			{
				Name: "nginx",
			},
		},
	}
	deployApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	updateApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						NodePort:      9090,
						ServiceType:   string(corev1.ServiceTypeNodePort),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	sApp.EXPECT().List("default", gomock.Any()).Return(appList, nil).Times(1)
	sApp.EXPECT().Get("default", "nginx", "").Return(deployApp, nil).Times(1)
	sApp.EXPECT().Update(nil, "default", updateApp).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "service.yaml")
	io.Copy(fw, strings.NewReader(updateService))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_DeleteService(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create deploy app
	appList := &gmodels.ApplicationList{
		Total: 1,
		Items: []gmodels.AppItem{
			{
				Name: "nginx",
			},
		},
	}
	deployApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						NodePort:      9090,
						ServiceType:   string(corev1.ServiceTypeNodePort),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	updateApp := &specV1.Application{
		Name:      "nginx",
		Namespace: "default",
		Labels: map[string]string{
			"app":             "nginx",
			"baetyl-app-mode": "kube",
		},
		Type: common.ContainerApp,
		Mode: context.RunModeKube,
		Services: []specV1.Service{
			{
				Name:  "nginx",
				Image: "nginx:latest",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "common-cm",
						MountPath: "/etc/config",
					},
					{
						Name:      "dcell",
						MountPath: "/etc/secret",
					},
					{
						Name:      "cache-volume",
						MountPath: "/cache",
					},
					{
						Name:      "test-volume",
						MountPath: "/test-hp",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						ContainerPort: 80,
						ServiceType:   string(corev1.ServiceTypeClusterIP),
					},
				},
				Resources: &specV1.Resources{},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "common-cm",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "common-cm",
					},
				},
			},
			{
				Name: "dcell",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "dcell",
					},
				},
			},
			{
				Name: "cache-volume",
				VolumeSource: specV1.VolumeSource{
					EmptyDir: &specV1.EmptyDirVolumeSource{
						Medium:    "",
						SizeLimit: "",
					},
				},
			},
			{
				Name: "test-volume",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "/var/lib/baetyl",
						Type: "Directory",
					},
				},
			},
		},
		Replica:  1,
		Workload: "deployment",
	}

	sApp.EXPECT().List("default", gomock.Any()).Return(appList, nil).Times(1)
	sApp.EXPECT().Get("default", "nginx", "").Return(deployApp, nil).Times(1)
	sApp.EXPECT().Update(nil, "default", updateApp).Return(updateApp, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "service.yaml")
	io.Copy(fw, strings.NewReader(testService))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml/delete", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}
