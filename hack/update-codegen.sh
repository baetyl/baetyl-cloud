GENERATOR_PATH=$GOPATH/src/k8s.io/code-generator

${GENERATOR_PATH}/generate-groups.sh "deepcopy,client" \
  github.com/baetyl/baetyl-cloud/plugin/kube/client github.com/baetyl/baetyl-cloud/plugin/kube/apis \
  "cloud:v1alpha1"
