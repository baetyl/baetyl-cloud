#!/bin/sh

ADDR="{{GetProperty "init-server-address"}}"
DEPLOYYML="{{.InitApplyYaml}}"
DB_PATH='{{.DBPath}}'
TOKEN="{{.Token}}"
MODE='{{.Mode}}'
SUDO=sudo

exec_cmd_nobail() {	
  echo "+ $2 bash -c \"$1\""	
  $2 bash -c "$1"	
}

print_status() {	
  echo "## $1"	
}

dbfile_clean() {
  INIT_DBFILE=$DB_PATH/store/core.db
  CORE_DBFILE=$DB_PATH/init/store/core.db
  if [ -f $INIT_DBFILE ]; then
    exec_cmd_nobail "rm -rf $INIT_DBFILE" $SUDO
    print_status "init db file deleted"
  fi
  if [ -f $CORE_DBFILE ]; then
    exec_cmd_nobail "rm -rf $CORE_DBFILE" $SUDO
    print_status "core db file deleted"
  fi
}

kube_clean() {
  exec_cmd_nobail "kubectl delete clusterrolebinding baetyl-edge-system-rbac --ignore-not-found=true" $SUDO
  exec_cmd_nobail "kubectl delete ns baetyl-edge-system --ignore-not-found=true" $SUDO
}

download_tool=curl
check_download_tool() {
	if [ -x "$(command -v wget)" ]; then
		download_tool=wget
	fi
}

kube_apply() {	
  TempFile=$(mktemp temp.XXXXXX)
  if [ "$download_tool" = "wget" ]; then
    exec_cmd_nobail "wget --no-check-certificate -O $TempFile \"$1\"" $SUDO
  else
    exec_cmd_nobail "curl -skfL \"$1\" >$TempFile" $SUDO
  fi
  exec_cmd_nobail "kubectl apply -f $TempFile" $SUDO
  exec_cmd_nobail "rm -f $TempFile 2>/dev/null" $SUDO
}

install_baetyl() {
  dbfile_clean
  if [ $MODE = "kube" ]; then
    print_status "baetyl install in k8s mode"
    kube_clean
    kube_apply "$ADDR/v1/init/$DEPLOYYML?token=$TOKEN"
  elif [ $MODE = "native" ]; then
    print_status "baetyl install in native mode"
    exec_cmd_nobail "baetyl delete" $SUDO
    exec_cmd_nobail "baetyl apply -f '$ADDR/v1/init/$DEPLOYYML?token=$TOKEN' --skip-verify=true" $SUDO
  else
    print_status "Not supported install mode $MODE"
    exit 0
  fi
}

install_baetyl
