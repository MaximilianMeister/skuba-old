/*
 * Copyright (c) 2019 SUSE LLC. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package kubernetes

type Component string

const (
	Etcd    Component = "etcd"
	CoreDNS Component = "coredns"
)

type ControlPlaneComponentsVersion struct {
	EtcdVersion    string
	CoreDNSVersion string
}

type ComponentsVersion struct {
	KubeletVersion string
}

type KubernetesVersion struct {
	ControlPlaneComponentsVersion ControlPlaneComponentsVersion
	ComponentsVersion             ComponentsVersion
}

const (
	// TBD depending on the exact versions we publish on the registry
	CurrentVersion = "v1.13.3"
)

var (
	Versions = map[string]KubernetesVersion{
		// TBD depending on the exact versions we publish on the registry
		"v1.13.4": KubernetesVersion{
			ControlPlaneComponentsVersion: ControlPlaneComponentsVersion{
				EtcdVersion:    "3.2.24",
				CoreDNSVersion: "1.3.1",
			},
			ComponentsVersion: ComponentsVersion{
				KubeletVersion: "1.13.4",
			},
		},
		"v1.13.3": KubernetesVersion{
			ControlPlaneComponentsVersion: ControlPlaneComponentsVersion{
				EtcdVersion:    "3.2.24",
				CoreDNSVersion: "1.3.1",
			},
			ComponentsVersion: ComponentsVersion{
				KubeletVersion: "1.13.3",
			},
		},
	}
)

func CurrentComponentVersion(component Component) string {
	currentKubernetesVersion := Versions[CurrentVersion]
	switch component {
	case Etcd:
		return currentKubernetesVersion.ControlPlaneComponentsVersion.EtcdVersion
	case CoreDNS:
		return currentKubernetesVersion.ControlPlaneComponentsVersion.CoreDNSVersion
	}
	return ""
}
