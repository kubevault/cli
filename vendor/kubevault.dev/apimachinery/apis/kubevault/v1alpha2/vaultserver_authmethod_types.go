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
	core "k8s.io/api/core/v1"
)

// +kubebuilder:validation:Enum=kubernetes;aws;gcp;userpass;cert;azure;jwt;oidc
type AuthMethodType string

const (
	AuthTypeKubernetes AuthMethodType = "kubernetes"
	AuthTypeAws        AuthMethodType = "aws"
	AuthTypeGcp        AuthMethodType = "gcp"
	AuthTypeUserPass   AuthMethodType = "userpass"
	AuthTypeCert       AuthMethodType = "cert"
	AuthTypeAzure      AuthMethodType = "azure"
	AuthTypeJWT        AuthMethodType = "jwt"
	AuthTypeOIDC       AuthMethodType = "oidc"
)

// AuthMethod contains the information to enable vault auth method
// links: https://www.vaultproject.io/api/system/auth.html
type AuthMethod struct {
	//  Specifies the name of the authentication method type, such as "github" or "token".
	Type AuthMethodType `json:"type"`

	// Specifies the path in which to enable the auth method.
	// Default value is the same as the 'type'
	Path string `json:"path"`

	// Specifies a human-friendly description of the auth method.
	// +optional
	Description string `json:"description,omitempty"`

	// Kubernetes auth config
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`

	// OIDC auth config
	OIDCConfig *JWTOIDCConfig `json:"oidcConfig,omitempty"`

	// JWT auth config
	JWTConfig *JWTOIDCConfig `json:"jwtConfig,omitempty"`

	// Specifies the name of the auth plugin to use based from the name in the plugin catalog.
	// Applies only to plugin methods.
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// Specifies if the auth method is a local only. Local auth methods are not replicated nor (if a secondary) removed by replication.
	// +optional
	Local bool `json:"local,omitempty"`
}

// +kubebuilder:validation:Enum=EnableSucceeded;EnableFailed;DisableSucceeded;DisableFailed
type AuthMethodEnableDisableStatus string

const (
	AuthMethodEnableSucceeded  AuthMethodEnableDisableStatus = "EnableSucceeded"
	AuthMethodEnableFailed     AuthMethodEnableDisableStatus = "EnableFailed"
	AuthMethodDisableSucceeded AuthMethodEnableDisableStatus = "DisableSucceeded"
	AuthMethodDisableFailed    AuthMethodEnableDisableStatus = "DisableFailed"
)

// AuthMethodStatus specifies the status of the auth method maintained by the auth method controller
type AuthMethodStatus struct {
	//  Specifies the name of the authentication method type, such as "github" or "token".
	Type AuthMethodType `json:"type"`

	// Specifies the path in which to enable the auth method.
	Path string `json:"path"`

	// Specifies whether auth method is enabled or not
	Status AuthMethodEnableDisableStatus `json:"status"`

	// Specifies the reason why failed to enable auth method
	// +optional
	Reason string `json:"reason,omitempty"`
}

type KubernetesConfig struct {
	// The default lease duration, specified as a string duration like "5s" or "30m".
	// +optional
	DefaultLeaseTTL string `json:"defaultLeaseTTL,omitempty"`

	// The maximum lease duration, specified as a string duration like "5s" or "30m".
	// +optional
	MaxLeaseTTL string `json:"maxLeaseTTL,omitempty"`

	// The name of the plugin in the plugin catalog to use.
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of keys that will not be HMAC'd by audit devices in the request data object.
	// +optional
	AuditNonHMACRequestKeys []string `json:"auditNonHMACRequestKeys,omitempty"`

	// List of keys that will not be HMAC'd by audit devices in the response data object.
	// +optional
	AuditNonHMACResponseKeys []string `json:"auditNonHMACResponseKeys,omitempty"`

	// Speficies whether to show this mount in the UI-specific listing endpoint.
	// +optional
	ListingVisibility string `json:"listingVisibility,omitempty"`

	// List of headers to whitelist and pass from the request to the backend.
	// +optional
	PassthroughRequestHeaders []string `json:"passthroughRequestHeaders,omitempty"`
}

type JWTOIDCConfig struct {
	// The default lease duration, specified as a string duration like "5s" or "30m".
	// +optional
	DefaultLeaseTTL string `json:"defaultLeaseTTL,omitempty"`

	// The maximum lease duration, specified as a string duration like "5s" or "30m".
	// +optional
	MaxLeaseTTL string `json:"maxLeaseTTL,omitempty"`

	// The name of the plugin in the plugin catalog to use.
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of keys that will not be HMAC'd by audit devices in the request data object.
	// +optional
	AuditNonHMACRequestKeys []string `json:"auditNonHMACRequestKeys,omitempty"`

	// List of keys that will not be HMAC'd by audit devices in the response data object.
	// +optional
	AuditNonHMACResponseKeys []string `json:"auditNonHMACResponseKeys,omitempty"`

	// Speficies whether to show this mount in the UI-specific listing endpoint.
	// +optional
	ListingVisibility string `json:"listingVisibility,omitempty"`

	// List of headers to whitelist and pass from the request to the backend.
	// +optional
	PassthroughRequestHeaders []string `json:"passthroughRequestHeaders,omitempty"`

	// CredentialSecretRef
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// TLSSecretRef
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// common configuration parameters
	// The OIDC Discovery URL, without any .well-known component (base path). Cannot be used with "jwks_url" or "jwt_validation_pubkeys".
	// +optional
	OIDCDiscoveryURL string `json:"oidcDiscoveryURL,omitempty"`

	// The OAuth Client ID from the provider for OIDC roles.
	// +optional
	OIDCClientID string `json:"oidcClientID,omitempty"`

	// The response mode to be used in the OAuth2 request. Allowed values are "query" and "form_post". Defaults to "query".
	// If using Vault namespaces, and oidc_response_mode is "form_post", then "namespace_in_state" should be set to false.
	// +optional
	OIDCResponseMode string `json:"oidcResponseMode,omitempty"`

	// (comma-separated string, or array of strings: <optional>) - The response types to request.
	// Allowed values are "code" and "id_token". Defaults to "code". Note: "id_token" may only be used if "oidc_response_mode" is set to "form_post".
	// +optional
	OIDCResponseTypes string `json:"oidcResponseTypes,omitempty"`

	// The default role to use if none is provided during login
	// +optional
	DefaultRole string `json:"defaultRole,omitempty"`

	// Configuration options for provider-specific handling.
	// Providers with specific handling include: Azure, Google. The options are described in each provider's section in OIDC Provider Setup.
	// +optional
	ProviderConfig map[string]string `json:"providerConfig,omitempty"`

	// JWKS URL to use to authenticate signatures. Cannot be used with "oidc_discovery_url" or "jwt_validation_pubkeys".
	// +optional
	JWKSURL string `json:"jwksURL,omitempty"`

	// (comma-separated string, or array of strings: <optional>)
	// A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with "jwks_url" or "oidc_discovery_url".
	JWTValidationPubkeys []string `json:"jwtValidationPubkeys,omitempty"`

	// (comma-separated string, or array of strings: <optional>)
	// A list of supported signing algorithms. Defaults to [RS256] for OIDC roles. Defaults to all available algorithms for JWT roles.
	// +optional
	JWTSupportedAlgs []string `json:"jwtSupportedAlgs,omitempty"`

	// The value against which to match the iss claim in a JWT.
	// +optional
	BoundIssuer string `json:"boundIssuer,omitempty"`
}
