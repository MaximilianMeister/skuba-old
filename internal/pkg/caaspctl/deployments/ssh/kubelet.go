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

package ssh

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"suse.com/caaspctl/internal/pkg/caaspctl/deployments/ssh/assets"
	"suse.com/caaspctl/pkg/caaspctl"
)

type KubeadmConfiguration struct {
	ImageRepository string
}

func renderTemplate(templateContents string, kubeadmConfiguration KubeadmConfiguration) string {
	template, err := template.New("").Parse(templateContents)
	if err != nil {
		log.Fatal("could not parse template")
	}
	var rendered bytes.Buffer
	if err := template.Execute(&rendered, kubeadmConfiguration); err != nil {
		log.Fatal("could not render configuration")
	}
	return rendered.String()
}

func init() {
	stateMap["kubelet.configure"] = kubeletConfigure()
	stateMap["kubelet.enable"] = kubeletEnable()
}

func kubeletConfigure() Runner {
	kubeadmConfiguration := KubeadmConfiguration{ImageRepository: caaspctl.ImageRepository}
	return func(t *Target, data interface{}) error {
		osRelease, err := t.target.OSRelease()
		if err != nil {
			return err
		}
		if strings.Contains(osRelease["ID_LIKE"], "suse") {
			if err := t.UploadFileContents("/usr/lib/systemd/system/kubelet.service", assets.KubeletService); err != nil {
				return err
			}
			if err := t.UploadFileContents("/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf", renderTemplate(assets.KubeadmService, kubeadmConfiguration)); err != nil {
				return err
			}
			if err := t.UploadFileContents("/etc/sysconfig/kubelet", assets.KubeletSysconfig); err != nil {
				return err
			}
		} else {
			if err := t.UploadFileContents("/lib/systemd/system/kubelet.service", assets.KubeletService); err != nil {
				return err
			}
			if err := t.UploadFileContents("/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", assets.KubeadmService); err != nil {
				return err
			}
		}
		_, _, err = t.ssh("systemctl", "daemon-reload")
		return err
	}
}

func kubeletEnable() Runner {
	return func(t *Target, data interface{}) error {
		_, _, err := t.ssh("systemctl", "enable", "kubelet")
		return err
	}
}
