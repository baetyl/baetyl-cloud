#!/bin/sh

set -e

OS=$(uname)
TOKEN="{{.Token}}"
CLOUD_ADDR="{{.CloudAddr}}"
SUDO=sudo

exec_cmd_nobail() {
  echo "+ $2 bash -c \"$1\""
  $2 bash -c "$1"
}

print_status() {
  echo "## $1"
}

url_safe_check() {
  if ! curl -Ifs $1 >/dev/null; then
    print_status "ERROR: $1 is invalid or Unreachable!"
  fi
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

install_docker() {
  TARGET_URL=http://get.daocloud.io/docker/
  url_safe_check ${TARGET_URL}
  exec_cmd_nobail "curl -sSL ${TARGET_URL} | $SUDO sh"

  if [[ ! -x "$(command -v docker)" ]]; then
      print_status "Install docker failed! Check the installing process for help..."
  fi

  if [[ ! -x "$(command -v systemctl)" ]]; then
      LSB_DIST=$(. /etc/os-release && echo "$ID" | tr '[:upper:]' '[:lower:]')
      case "$LSB_DIST" in
      ubuntu | debian | raspbian)
          exec_cmd_nobail "apt update && apt install --no-install-recommends -y systemd  >/dev/null 2>&1" $SUDO
          ;;
      centos)
          exec_cmd_nobail "yum install systemd -y  >/dev/null 2>&1" $SUDO
          ;;
      *)
          print_status "Your OS: $OS is not supported!"
          exit 0
          ;;
      esac
  fi

  exec_cmd_nobail "systemctl enable docker" $SUDO
  exec_cmd_nobail "systemctl start docker" $SUDO
}

check_and_get_kube() {
  if [ ! -x "$(check_cmd kubectl)" ]; then
    read -p "K8S/K3S is not installed yet, do you want us to install K3S for you? Yes/No (default: Yes):" IS_INSTALL_K3S
    if [ "$IS_INSTALL_K3S" = "n" -o "$IS_INSTALL_K3S" = "N" -o "$IS_INSTALL_K3S" = "no" -o "$IS_INSTALL_K3S" = "NO" ]; then
      echo "K3S is needed to run ${NAME}, this script will exit now..."
      exit 0
    fi

    if [ $OS = "Linux" ]; then
        read -p "K3S could run with containerd/docker, which do you want us to install for you? containerd for Yes, docker for No (default: Yes):" IS_INSTALL_CONTAINERD
        if [ "$IS_INSTALL_CONTAINERD" = "n" -o "$IS_INSTALL_CONTAINERD" = "N" -o "$IS_INSTALL_CONTAINERD" = "no" -o "$IS_INSTALL_CONTAINERD" = "NO" ]; then
          if [ ! -x "$(check_cmd docker)" ]; then
            install_docker
          else
            print_status "Docker already installed"
          fi
          export INSTALL_K3S_EXEC="--docker --write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666"
        else
          export INSTALL_K3S_EXEC="--write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666"
        fi

      exec_cmd_nobail "curl -sfL https://docs.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -"

      if [ ! -x "$(check_cmd kubectl)" ]; then
        print_status "Install k3s failed! Check the installing process for help..."
        exit 0
      fi

      exec_cmd_nobail "systemctl enable k3s" $SUDO
      exec_cmd_nobail "systemctl start k3s" $SUDO

    elif [ $OS = "Darwin" ]; then
      exec_cmd_nobail "curl -sfL https://get.k3s.io | sh -"
    else
      print_status "We are not supporting your system, this script will exit now..."
      exit 0
    fi
  fi
  exec_cmd_nobail "kubectl version" $SUDO
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
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/core-data" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/app-data" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/core-store" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/log/baetyl/core-log" $SUDO
    exec_cmd_nobail "mkdir -p -m 666 /var/lib/baetyl/core-page" $SUDO
    kube_apply "$CLOUD_ADDR/v1/active/baetyl-init.yml?token=$TOKEN&node=$KUBE_MASTER_NODE_NAME"
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
    kube_apply "$CLOUD_ADDR/v1/active/metrics.yml"
  fi
}

check_kube_res() {
  $SUDO kubectl get deployments -A | grep $1 | awk '{print $2}'
}

check_and_get_storage() {
  PROVISIONER=$(check_kube_res local-path-provisioner)
  if [ -z "$PROVISIONER" ]; then
    kube_apply "$$CLOUD_ADDR/v1/active/local-path-storage.yml"
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
  check_and_get_kube
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