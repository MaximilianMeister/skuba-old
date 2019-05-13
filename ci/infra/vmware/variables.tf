variable "VSPHERE_SERVER" {}
variable "VSPHERE_USER" {}
variable "VSPHERE_PASSWORD" {}
variable "VSPHERE_ALLOW_UNVERIFIED_SSL" {}
variable "template_name" {}
variable "stack_name" {}
variable "vsphere_datastore" {}
variable "vsphere_datacenter" {}
variable "vsphere_network" {}
variable "vsphere_resource_pool" {}

variable "authorized_keys" {
  type        = "list"
  default     = []
  description = "ssh keys to inject into all the nodes"
}

variable "repositories" {
  type        = "list"
  default     = []
  description = "Urls of the repositories to mount via cloud-init"
}

variable "ntp_servers" {
  type        = "list"
  default     = ["0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org"]
  description = "list of ntp servers to configure"
}

variable "packages" {
  type        = "list"
  default     = []
  description = "list of additional packages to install"
}

variable "username" {
  default     = "sles"
  description = "Username for the cluster nodes"
}

variable "password" {
  default     = "sles"
  description = "Password for the cluster nodes"
}

variable "masters" {
  default     = 1
  description = "Number of master nodes"
}

variable "workers" {
  default     = 1
  description = "Number of worker nodes"
}

variable "load-balancers" {
  default     = 1
  description = "Number of load-balancer nodes"
}

variable "worker_cpus" {
  default     = 4
  description = "Number of CPUs used on worker node"
}

variable "worker_memory" {
  default     = 8192
  description = "Amount of memory used on worker node"
}

variable "master_cpus" {
  default     = 4
  description = "Number of CPUs used on master node"
}

variable "master_memory" {
  default     = 8192
  description = "Amount of memory used on master node"
}

variable "lb_cpus" {
  default     = 1
  description = "Number of CPUs used on load-balancer node"
}

variable "lb_memory" {
  default     = 2048
  description = "Amount of memory used on load-balancer node"
}

#### To be moved to separate vsphere.tf? ####

provider "vsphere" {
  vsphere_server       = "${var.VSPHERE_SERVER}"
  user                 = "${var.VSPHERE_USER}"
  password             = "${var.VSPHERE_PASSWORD}"
  allow_unverified_ssl = "${var.VSPHERE_ALLOW_UNVERIFIED_SSL}"
}

data "vsphere_resource_pool" "pool" {
  name          = "${var.vsphere_resource_pool}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "datastore" {
  name          = "${var.vsphere_datastore}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

data "vsphere_network" "network" {
  name          = "${var.vsphere_network}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_virtual_machine" "template" {
  name          = "${var.template_name}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}
