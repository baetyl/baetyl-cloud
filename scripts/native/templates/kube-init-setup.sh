#!/bin/sh

set -e

OS=$(uname)
TOKEN="{{.Token}}"
ADDR="{{GetProperty "init-server-address"}}"
DEPLOYYML="{{.DeploymentYml}}"
SUDO=sudo

exec_cmd_nobail() {
  echo "+ $2 bash -c \"$1\""
  $2 bash -c "$1"
}

print_status() {
  echo "## $1"
}

check_cmd() {
  command -v $1 | awk '{print}'
}

get_dependencies() {
  PRE_INSTALL_PKGS=""

  if [ ! -x "$(check_cmd curl)" ]; then
    PRE_INSTALL_PKGS="${PRE_INSTALL_PKGS} curl"
  fi

  if [ "X${PRE_INSTALL_PKGS}" != "X" ]; then
    case "$OS" in
    Linux)
      LSB_DIST=$(. /etc/os-release && echo "$ID" | tr '[:upper:]' '[:lower:]')
      case "$LSB_DIST" in
      ubuntu | debian | raspbian)
        exec_cmd_nobail "apt update && apt install --no-install-recommends -y ${PRE_INSTALL_PKGS} >/dev/null 2>&1" $SUDO
        ;;
      centos)
        exec_cmd_nobail "yum install ${PRE_INSTALL_PKGS} -y >/dev/null 2>&1" $SUDO
        ;;
      *)
        print_status "Your OS is not supported!"
        ;;
      esac
      ;;
    Darwin)
      print_status "You must install ${PRE_INSTALL_PKGS} to continue..."
      exit 0
      ;;
    *)
      print_status "Your OS: $OS is not supported!"
      exit 0
      ;;
    esac
  fi
}

get_kube_master() {
  $SUDO kubectl get no | awk '/master/ || /controlplane/' | awk '{print $1}'
}

check_baetyl_namespace() {
  $SUDO kubectl get ns | grep 'baetyl-edge-system' | awk '{print $1}'
}

check_and_install_baetyl() {
  BAETYL_NAMESPACE=$(check_baetyl_namespace)
  if [ ! -z "$BAETYL_NAMESPACE" ]; then
    read -p "The namespace 'baetyl-edge-system' already exists, do you want to clean up old applications by deleting this namespace? Yes/No (default: Yes):" IS_DELETE_NS
    if [ "$IS_DELETE_NS" = "n" -o "$IS_DELETE_NS" = "N" -o "$IS_DELETE_NS" = "no" -o "$IS_DELETE_NS" = "NO" ]; then
      echo "baetyl-init is not install, this script will exit now..."
      exit 0
    else
      rbac=$($SUDO kubectl get clusterrolebinding | grep baetyl-edge-system-rbac | awk '{print $1}')
      if [ -n "$rbac" ]; then
        exec_cmd_nobail "kubectl delete clusterrolebinding baetyl-edge-system-rbac" $SUDO
      fi
      exec_cmd_nobail "kubectl delete namespace baetyl-edge-system" $SUDO
    fi
  fi

  KUBE_MASTER_NODE_NAME=$(get_kube_master)
  if [ ! -z "$KUBE_MASTER_NODE_NAME" ]; then
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/host" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/object" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/store" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/log" $SUDO
    kube_apply "$ADDR/v1/init/$DEPLOYYML?token=$TOKEN&node=$KUBE_MASTER_NODE_NAME"
  else
    print_status "Can not get kubernetes master or controlplane node, this script will exit now..."
  fi
}

kube_apply() {
  TempFile=$(mktemp temp.XXXXXX)
  exec_cmd_nobail "curl -skfL \"$1\" >$TempFile" $SUDO
  exec_cmd_nobail "kubectl apply -f $TempFile" $SUDO
  exec_cmd_nobail "rm -f $TempFile 2>/dev/null" $SUDO
}

check_and_get_metrics() {
  METRICS=$(check_kube_res metrics-server)
  if [ -z "$METRICS" ]; then
    kube_apply "$ADDR/v1/init/kube-api-metrics.yml?token=$TOKEN"
  fi
}

check_kube_res() {
  $SUDO kubectl get deployments -A | grep $1 | awk '{print $2}'
}

check_and_get_storage() {
  PROVISIONER=$(check_kube_res local-path-provisioner)
  if [ -z "$PROVISIONER" ]; then
    kube_apply "$ADDR/v1/init/kube-local-path-storage.yml?token=$TOKEN"
  fi
}

check_user() {
  if [ $(id -u) -eq 0 ]; then
    SUDO=
  fi
}

uninstall() {
  exec_cmd_nobail "kubectl delete ns baetyl-edge baetyl-edge-system" $SUDO
}

install() {
  check_user
  get_dependencies
  check_and_get_metrics
  check_and_get_storage
  check_and_install_baetyl
}

case C"$1" in
C)
  install
  ;;
Cuninstall)
  uninstall
  ;;
Cinstall)
  install
  ;;
C*)
  Usage: setup.sh { install | uninstall }
  ;;
esac

echo "Done!"
exit 0