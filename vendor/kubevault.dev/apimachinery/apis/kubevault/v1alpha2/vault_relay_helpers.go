/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"fmt"

	"kubevault.dev/apimachinery/apis/kubevault"
	"kubevault.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (*VaultRelay) Hub() {}

func (vr VaultRelay) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceVaultRelays))
}

func (vr VaultRelay) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourceVaultRelays, kubevault.GroupName)
}

func (vr VaultRelay) GetKey() string {
	return vr.Namespace + "/" + vr.Name
}

func (vr VaultRelay) OffshootName() string {
	return vr.Name
}

func (vr VaultRelay) ServiceAccountName() string {
	return vr.Name
}

func (vr VaultRelay) PodName() string {
	return meta_util.NameWithSuffix(vr.Name, "relay")
}

func (vr VaultRelay) AppBindingName() string {
	return meta_util.NameWithSuffix(vr.Name, "hub-vault")
}

func (vr VaultRelay) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      vr.ResourceFQN(),
		meta_util.InstanceLabelKey:  vr.Name,
		meta_util.ManagedByLabelKey: kubevault.GroupName,
	}
}

func (vr VaultRelay) OffshootLabels() map[string]string {
	return meta_util.FilterKeys("kubevault.com", vr.OffshootSelectors(), vr.Labels)
}

func (vr VaultRelay) IsValid() error {
	return nil
}

// GetGRPCPort returns the gRPC port with default value
func (vr VaultRelay) GetGRPCPort() int32 {
	if vr.Spec.HubVaultRef.GRPCPort == 0 {
		return 50053
	}
	return vr.Spec.HubVaultRef.GRPCPort
}

// GetImage returns the spoke-relay image with default value
func (vr VaultRelay) GetImage() string {
	if vr.Spec.Image == "" {
		return "ghcr.io/kubevault/spoke-relay:latest"
	}
	return vr.Spec.Image
}

// SetDefaults sets default values for VaultRelay
func (vr *VaultRelay) SetDefaults() {
	if vr.Spec.HubVaultRef.GRPCPort == 0 {
		vr.Spec.HubVaultRef.GRPCPort = 50053
	}

	if vr.Spec.Reconnect == nil {
		vr.Spec.Reconnect = &ReconnectConfig{
			Enabled:           true,
			BackoffSeconds:    5,
			MaxBackoffSeconds: 300,
		}
	}
}
