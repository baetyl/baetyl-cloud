set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=$GOPATH/src/k8s.io/code-generator

bash $CODEGEN_PKG/generate-groups.sh \
  "deepcopy,client" \
  github.com/baetyl/baetyl-cloud/v2/plugin/kube/client \
  github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis \
  cloud:v1alpha1 \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

# 注意 code-generator 库的版本要和 go.mod 中引用的 k8s 库保持一致
# 期望生成的函数列表 deepcopy,defaulter,client,lister,informer
# 生成代码的目标目录，格式为 go.mod 中的项目名称+输出路径
# CRD 所在目录，格式为 go.mod 中的项目名称+apis目录
# CRD的group name和version