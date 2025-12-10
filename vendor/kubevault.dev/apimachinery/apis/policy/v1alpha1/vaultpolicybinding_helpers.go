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

package v1alpha1

import (
	"encoding/json"
	"fmt"

	"kubevault.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	clustermeta "kmodules.xyz/client-go/cluster"
	meta_util "kmodules.xyz/client-go/meta"
)

func (v VaultPolicyBinding) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceVaultPolicyBindings))
}

func (v VaultPolicyBinding) GetKey() string {
	return ResourceVaultPolicyBinding + "/" + v.Namespace + "/" + v.Name
}

func (v VaultPolicyBinding) PolicyBindingName() string {
	if v.Spec.VaultRoleName != "" {
		return v.Spec.VaultRoleName
	}

	cluster := "-"
	if clustermeta.ClusterName() != "" {
		cluster = clustermeta.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, v.Namespace, v.Name)
}

func (v VaultPolicyBinding) OffshootSelectors() map[string]string {
	return map[string]string{
		"app":                  "vault",
		"vault_policy_binding": v.Name,
	}
}

func (v VaultPolicyBinding) OffshootLabels() map[string]string {
	return meta_util.FilterKeys("kubevault.com", v.OffshootSelectors(), v.Labels)
}

func (v VaultPolicyBinding) IsValid() error {
	return nil
}

func (v *VaultPolicyBinding) SetDefaults() {
	if v == nil {
		return
	}

	if v.Spec.VaultRoleName == "" {
		v.Spec.VaultRoleName = v.PolicyBindingName()
	}

	if v.Spec.Kubernetes != nil {
		if v.Spec.Kubernetes.Path == "" {
			v.Spec.Kubernetes.Path = "kubernetes"
		}
		if v.Spec.Kubernetes.Name == "" {
			v.Spec.Kubernetes.Name = v.PolicyBindingName()
		}
	}

	if v.Spec.AppRole != nil {
		if v.Spec.AppRole.Path == "" {
			v.Spec.AppRole.Path = "approle"
		}
		if v.Spec.AppRole.RoleName == "" {
			v.Spec.AppRole.RoleName = v.PolicyBindingName()
		}
	}

	if v.Spec.LdapGroup != nil {
		if v.Spec.LdapGroup.Path == "" {
			v.Spec.LdapGroup.Path = "ldap"
		}
	}

	if v.Spec.LdapUser != nil {
		if v.Spec.LdapUser.Path == "" {
			v.Spec.LdapUser.Path = "ldap"
		}
	}

	if v.Spec.JWT != nil {
		if v.Spec.JWT.Path == "" {
			v.Spec.JWT.Path = "jwt"
		}
		if v.Spec.JWT.Name == "" {
			v.Spec.JWT.Name = v.PolicyBindingName()
		}
	}

	if v.Spec.OIDC != nil {
		if v.Spec.OIDC.Path == "" {
			v.Spec.OIDC.Path = "oidc"
		}
		if v.Spec.OIDC.Name == "" {
			v.Spec.OIDC.Name = v.PolicyBindingName()
		}
	}
}

func (v VaultPolicyBinding) GeneratePayload(i any) (map[string]any, error) {
	var err error
	payload := make(map[string]any)
	byte, err := json.Marshal(i)
	if err == nil {
		err = json.Unmarshal(byte, &payload)
	}
	return payload, err
}

func (v VaultPolicyBinding) GeneratePath(name, path, subPath string) string {
	return fmt.Sprintf("auth/%s/%s/%s", path, subPath, name)
}
