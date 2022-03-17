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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindVaultPolicyBinding = "VaultPolicyBinding"
	ResourceVaultPolicyBinding     = "vaultpolicybinding"
	ResourceVaultPolicyBindings    = "vaultpolicybindings"
)

// VaultPolicyBinding binds a list of Vault server policies with Vault users authenticated by various auth methods.
// Currently VaultPolicyBinding only supports users authenticated via Kubernetes auth method.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vaultpolicybindings,singular=vaultpolicybinding,shortName=vpb,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VaultPolicyBinding struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VaultPolicyBindingSpec   `json:"spec,omitempty"`
	Status            VaultPolicyBindingStatus `json:"status,omitempty"`
}

// links: https://www.vaultproject.io/api/auth/kubernetes/index.html#parameters-1
type VaultPolicyBindingSpec struct {
	// VaultRef is the name of a AppBinding referencing to a Vault Server
	VaultRef core.LocalObjectReference `json:"vaultRef"`

	// VaultRoleName is the role name which will be bound of the policies
	// This defaults to following format: k8s.${cluster}.${metadata.namespace}.${metadata.name}
	// xref: https://www.vaultproject.io/api/auth/kubernetes/index.html#create-role
	// +optional
	VaultRoleName string `json:"vaultRoleName,omitempty"`

	// Policies is a list of Vault policy identifiers.
	Policies []PolicyIdentifier `json:"policies"`

	// SubjectRef refers to Vault users who will be granted policies.
	SubjectRef `json:"subjectRef"`
}

type PolicyIdentifier struct {
	// Name is a Vault server policy name. This name should be returned by `vault read sys/policy` command.
	// More info: https://www.vaultproject.io/docs/concepts/policies.html#listing-policies
	Name string `json:"name,omitempty"`

	// Ref is name of a VaultPolicy crd object. Actual vault policy name is spec.vaultRoleName field.
	// More info: https://www.vaultproject.io/docs/concepts/policies.html#listing-policies
	Ref string `json:"ref,omitempty"`
}

type SubjectRef struct {
	// Kubernetes refers to Vault users who are authenticated via Kubernetes auth method
	// More info: https://www.vaultproject.io/docs/auth/kubernetes.html#configuration
	Kubernetes *KubernetesSubjectRef `json:"kubernetes,omitempty"`
	// More info: https://www.vaultproject.io/docs/auth/approle#configuration
	AppRole *AppRoleSubjectRef `json:"appRole,omitempty"`
	// More info: https://www.vaultproject.io/api-docs/auth/ldap#configure-ldap
	LdapGroup *LdapGroupSubjectRef `json:"ldapGroup,omitempty"`
	LdapUser  *LdapUserSubjectRef  `json:"ldapUser,omitempty"`
	// More info: https://www.vaultproject.io/api-docs/auth/jwt#configure
	JWT *JWTSubjectRef `json:"jwt,omitempty"`
}

// More info: https://www.vaultproject.io/api/auth/kubernetes/index.html#create-role
type KubernetesSubjectRef struct {
	// Specifies the path where kubernetes auth is enabled
	// default : kubernetes
	// +optional
	Path string `json:"path,omitempty"`

	// Name of the role
	Name string `json:"name,omitempty"`

	// Specifies the names of the service account to bind with policy
	ServiceAccountNames []string `json:"serviceAccountNames"`

	// Specifies the namespaces of the service account
	ServiceAccountNamespaces []string `json:"serviceAccountNamespaces"`

	// Specifies the TTL period of tokens issued using this role in seconds.
	// +optional
	TTL string `json:"ttl,omitempty"`

	// Specifies the maximum allowed lifetime of tokens issued in seconds using this role.
	// +optional
	MaxTTL string `json:"maxTTL,omitempty"`

	// If set, indicates that the token generated using this role should never expire.
	// The token should be renewed within the duration specified by this value.
	// At each renewal, the token's TTL will be set to the value of this parameter.
	// +optional
	Period string `json:"period,omitempty"`
}

// More info: https://www.vaultproject.io/api-docs/auth/approle#create-update-approle
type AppRoleSubjectRef struct {
	// Specifies the path where approle auth is enabled
	// default : approle
	// +optional
	Path string `json:"path,omitempty"`

	// RoleName is the Name of the AppRole
	// This defaults to following format: k8s.${cluster}.${metadata.namespace}.${metadata.name}
	RoleName string `json:"roleName,omitempty"`

	// Require secret_id to be presented when logging in using this AppRole.
	BindSecretID bool `json:"bindSecretID"`

	// List of CIDR blocks; if set, specifies blocks of IP addresses which can perform the login operation.
	SecretIDBoundCidrs []string `json:"secretIdBoundCidrs,omitempty"`

	// Number of times any particular SecretID can be used to fetch a token from this AppRole, after which the SecretID will expire. A value of zero will allow unlimited uses.
	SecretIDNumUses int64 `json:"secretIdNumUses,omitempty"`

	// Duration in either an integer number of seconds (3600) or an integer time unit (60m) after which any SecretID expires.
	SecretIDTTL string `json:"secretIdTTL,omitempty"`

	// If set, the secret IDs generated using this role will be cluster local. This can only be set during role creation and once set, it can't be reset later.
	EnableLocalSecretIDs bool `json:"enableLocalSecretIds,omitempty"`

	// The incremental lifetime for generated tokens. This current value of this will be referenced at renewal time.
	TokenTTL int64 `json:"tokenTTL,omitempty"`

	// The maximum lifetime for generated tokens. This current value of this will be referenced at renewal time.
	TokenMaxTTL int64 `json:"tokenMaxTTL,omitempty"`

	// List of policies to encode onto generated tokens. Depending on the auth method, this list may be supplemented by user/group/other values.
	TokenPolicies []string `json:"tokenPolicies,omitempty"`

	// List of CIDR blocks; if set, specifies blocks of IP addresses which can authenticate successfully, and ties the resulting token to these blocks as well.
	TokenBoundCidrs []string `json:"tokenBoundCidrs,omitempty"`

	// If set, will encode an explicit max TTL onto the token. This is a hard cap even if token_ttl and token_max_ttl would otherwise allow a renewal.
	TokenExplicitMaxTTL int64 `json:"tokenExplicitMaxTTL,omitempty"`

	// If set, the default policy will not be set on generated tokens; otherwise it will be added to the policies set in token_policies.
	TokenNoDefaultPolicy bool `json:"tokenNoDefaultPolicy,omitempty"`

	// The maximum number of times a generated token may be used (within its lifetime); 0 means unlimited.
	TokenNumUses int64 `json:"tokenNumUses,omitempty"`

	// The period, if any, to set on the token.
	TokenPeriod int64 `json:"tokenPeriod,omitempty"`

	// The type of token that should be generated. Can be service, batch, or default to use the mount's tuned default (which unless changed will be service tokens). For token store roles, there are two additional possibilities: default-service and default-batch which specify the type to return unless the client requests a different type at generation time.
	TokenType string `json:"tokenType,omitempty"`
}

// More info: https://www.vaultproject.io/api-docs/auth/ldap#create-update-ldap-group
type LdapGroupSubjectRef struct {
	// Specifies the path where ldap groups auth is enabled
	// default : ldap/groups
	// +optional
	Path string `json:"path,omitempty"`

	// The name of the LDAP group
	Name string `json:"name"`

	// List of policies to encode onto generated tokens. Depending on the auth method, this list may be supplemented by user/group/other values.
	Policies []string `json:"policies,omitempty"`
}

// More info: https://www.vaultproject.io/api-docs/auth/ldap#create-update-ldap-user
type LdapUserSubjectRef struct {
	// Specifies the path where ldap groups auth is enabled
	// default : ldap/users
	// +optional
	Path string `json:"path,omitempty"`

	// The username of the LDAP user
	Username string `json:"username"`

	// List of policies to encode onto generated tokens. Depending on the auth method, this list may be supplemented by user/group/other values.
	Policies []string `json:"policies,omitempty"`

	// List of groups associated to the user.
	Groups []string `json:"groups,omitempty"`
}

// More info: https://www.vaultproject.io/api-docs/auth/jwt#create-role
type JWTSubjectRef struct {
	// Specifies the path where jwt/oidc auth is enabled
	// default : jwt
	// +optional
	Path string `json:"path,omitempty"`

	// Name of the role.
	// This defaults to following format: k8s.${cluster}.${metadata.namespace}.${metadata.name}
	Name string `json:"name,omitempty"`

	// Type of role, either "oidc" (default) or "jwt".
	RoleType string `json:"roleType,omitempty"`

	// List of aud claims to match against. Any match is sufficient. Required for "jwt" roles, optional for "oidc" roles.
	BoundAudiences []string `json:"boundAudiences,omitempty"`

	// The claim to use to uniquely identify the user; this will be used as the name for the Identity entity alias created due to a successful login. The claim value must be a string.
	UserClaim string `json:"userClaim"`

	// If set, requires that the sub claim matches this value.
	BoundSubject string `json:"boundSubject,omitempty"`

	// If set, a map of claims/values to match against. The expected value may be a single string or a list of strings. The interpretation of the bound claim values is configured with bound_claims_type.
	BoundClaims map[string]string `json:"boundClaims,omitempty"`

	// Configures the interpretation of the bound_claims values. If "string" (the default), the values will treated as string literals and must match exactly. If set to "glob", the values will be interpreted as globs, with * matching any number of characters.
	BoundClaimsType string `json:"boundClaimsType,omitempty"`

	// The claim to use to uniquely identify the set of groups to which the user belongs; this will be used as the names for the Identity group aliases created due to a successful login. The claim value must be a list of strings.
	GroupClaim string `json:"groupClaim,omitempty"`

	// If set, a map of claims (keys) to be copied to specified metadata fields (values).
	ClaimMappings map[string]string `json:"claimMappings,omitempty"`

	// If set, a list of OIDC scopes to be used with an OIDC role. The standard scope "openid" is automatically included and need not be specified.
	OIDCScopes []string `json:"oidcScopes,omitempty"`

	// The list of allowed values for redirect_uri during OIDC logins.
	AllowedRedirectUris []string `json:"allowedRedirectUris"`

	VerboseOIDCLogging bool `json:"verboseOidcLogging,omitempty"`

	// The incremental lifetime for generated tokens. This current value of this will be referenced at renewal time.
	TokenTTL int64 `json:"tokenTTL,omitempty"`

	// The maximum lifetime for generated tokens. This current value of this will be referenced at renewal time.
	TokenMaxTTL int64 `json:"tokenMaxTTL,omitempty"`

	// List of policies to encode onto generated tokens. Depending on the auth method, this list may be supplemented by user/group/other values.
	TokenPolicies []string `json:"tokenPolicies,omitempty"`

	// List of CIDR blocks; if set, specifies blocks of IP addresses which can authenticate successfully, and ties the resulting token to these blocks as well.
	TokenBoundCidrs []string `json:"tokenBoundCidrs,omitempty"`

	// If set, will encode an explicit max TTL onto the token. This is a hard cap even if token_ttl and token_max_ttl would otherwise allow a renewal.
	TokenExplicitMaxTTL int64 `json:"tokenExplicitMaxTTL,omitempty"`

	// If set, the default policy will not be set on generated tokens; otherwise it will be added to the policies set in token_policies.
	TokenNoDefaultPolicy bool `json:"tokenNoDefaultPolicy,omitempty"`

	// The maximum number of times a generated token may be used (within its lifetime); 0 means unlimited.
	TokenNumUses int64 `json:"tokenNumUses,omitempty"`

	// The period, if any, to set on the token.
	TokenPeriod int64 `json:"tokenPeriod,omitempty"`

	// The type of token that should be generated. Can be service, batch, or default to use the mount's tuned default (which unless changed will be service tokens). For token store roles, there are two additional possibilities: default-service and default-batch which specify the type to return unless the client requests a different type at generation time.
	TokenType string `json:"tokenType,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type VaultPolicyBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultPolicyBinding `json:"items,omitempty"`
}

// ServiceAccountReference contains name and namespace of the service account
type ServiceAccountReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// +kubebuilder:validation:Enum=Success;Failed
type PolicyBindingPhase string

const (
	PolicyBindingSuccess PolicyBindingPhase = "Success"
	PolicyBindingFailed  PolicyBindingPhase = "Failed"
)

type VaultPolicyBindingStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase indicates whether successfully bind the policy to service account in vault or not or in progress
	// +optional
	Phase PolicyBindingPhase `json:"phase,omitempty"`

	// Represents the latest available observations of a VaultPolicyBinding.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}
