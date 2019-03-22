#cloud-config

# set locale
locale: en_GB.UTF-8

# set timezone
timezone: Etc/UTC

# Set hostname and FQDN
hostname: ${hostname}
fqdn: ${fqdn}

# set root password
chpasswd:
  list: |
    root:linux
    opensuse:linux
  expire: False

ssh_authorized_keys:
${authorized_keys}

# need to disable gpg checks because the cloud image has an untrusted repo
zypper:
  repos:
    - id: caasp
      name: caasp
      baseurl: ${repo_baseurl}
      enabled: 1
      autorefresh: 1
      gpgcheck: 0
  config:
    gpgcheck: "off"
    solver.onlyRequires: "true"
    download.use_deltarpm: "true"

# need to remove the standard docker packages that are pre-installed on the
# cloud image because they conflict with the kubic- ones that are pulled by
# the kubernetes packages
packages:
  - kubernetes-kubeadm
  - kubernetes-kubelet
  - kubectl
  - cni-plugins
  - "-docker"
  - "-containerd"
  - "-docker-runc"
  - "-docker-libnetwork"

bootcmd:
  - ip link set dev eth0 mtu 1400

final_message: "The system is finally up, after $UPTIME seconds"
