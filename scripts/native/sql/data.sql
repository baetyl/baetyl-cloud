INSERT INTO `baetyl_property` (`name`, `value`) VALUES
('sync-server-address', 'https://host.docker.internal:9005'),
('init-server-address', 'https://0.0.0.0:9003'),

('baetyl-image', 'hub.baidubce.com/baetyl/baetyl:v2.1.0'),
('baetyl-function-image', 'hub.baidubce.com/baetyl/function:v2.1.0'),
('baetyl-broker-image', 'hub.baidubce.com/baetyl/broker:v2.1.0'),
('baetyl-function-runtime-python3','hub.baidubce.com/baetyl/function-python:3.6-v2.1.0'),
('baetyl-function-runtime-nodejs10','hub.baidubce.com/baetyl/function-node:10.19-v2.1.0'),
('baetyl-function-runtime-sql','hub.baidubce.com/baetyl/function-sql:v2.1.0'),
('baetyl-version-latest', 'v2.1.0'),

('command-docker-installation', 'curl -sSL https://get.daocloud.io/docker | sh'),
('command-k3s-installation-containerd', 'curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="--write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666" sh -'),
('command-k3s-installation-docker', 'curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="--docker --write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666" sh -'),

('baetyl-init-command', 'curl -skfL \'{{GetProperty \"init-server-address\"}}/v1/init/baetyl-install.sh?token={{.Token}}&mode={{.mode}}&initApplyYaml={{.InitApplyYaml}}\' -osetup.sh && sh setup.sh');