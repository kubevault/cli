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
	"unsafe"

	"kubevault.dev/apimachinery/apis/kubevault/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha1_MySQLSpec_To_v1alpha2_MySQLSpec(in *MySQLSpec, out *v1alpha2.MySQLSpec, s conversion.Scope) error {
	out.Address = in.Address
	out.Database = in.Database
	out.Table = in.Table
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.UserCredentialSecret,
	}
	out.TLSSecretRef = &core.LocalObjectReference{
		Name: in.TLSCASecret,
	}
	// WARNING: in.UserCredentialSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.TLSCASecret requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha2_MySQLSpec_To_v1alpha1_MySQLSpec(in *v1alpha2.MySQLSpec, out *MySQLSpec, s conversion.Scope) error {
	out.Address = in.Address
	out.Database = in.Database
	out.Table = in.Table
	out.UserCredentialSecret = in.CredentialSecretRef.Name
	out.TLSCASecret = in.TLSSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	// WARNING: in.TLSSecretRef requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	// WARNING: in.DatabaseRef requires manual conversion: does not exist in peer-type
	// WARNING: in.PlaintextCredentialTransmission requires manual conversion: does not exist in peer-type
	// WARNING: in.MaxIdleConnection requires manual conversion: does not exist in peer-type
	// WARNING: in.MaxConnectionLifetime requires manual conversion: does not exist in peer-type
	// WARNING: in.HAEnabled requires manual conversion: does not exist in peer-type
	// WARNING: in.LockTable requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_PostgreSQLSpec_To_v1alpha2_PostgreSQLSpec(in *PostgreSQLSpec, out *v1alpha2.PostgreSQLSpec, s conversion.Scope) error {
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.ConnectionURLSecret,
	}
	// WARNING: in.ConnectionURLSecret requires manual conversion: does not exist in peer-type
	out.Table = in.Table
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha2_PostgreSQLSpec_To_v1alpha1_PostgreSQLSpec(in *v1alpha2.PostgreSQLSpec, out *PostgreSQLSpec, s conversion.Scope) error {
	out.ConnectionURLSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.Table = in.Table
	out.MaxParallel = in.MaxParallel
	// WARNING: in.DatabaseRef requires manual conversion: does not exist in peer-type
	// WARNING: in.SSLMode requires manual conversion: does not exist in peer-type
	// WARNING: in.MaxIdleConnection requires manual conversion: does not exist in peer-type
	// WARNING: in.HAEnabled requires manual conversion: does not exist in peer-type
	// WARNING: in.HaTable requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_AwsKmsSsmSpec_To_v1alpha2_AwsKmsSsmSpec(in *AwsKmsSsmSpec, out *v1alpha2.AwsKmsSsmSpec, s conversion.Scope) error {
	out.KmsKeyID = in.KmsKeyID
	out.SsmKeyPrefix = in.SsmKeyPrefix
	out.Region = in.Region
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	out.Endpoint = in.Endpoint
	return nil
}

func Convert_v1alpha2_AwsKmsSsmSpec_To_v1alpha1_AwsKmsSsmSpec(in *v1alpha2.AwsKmsSsmSpec, out *AwsKmsSsmSpec, s conversion.Scope) error {
	out.KmsKeyID = in.KmsKeyID
	out.SsmKeyPrefix = in.SsmKeyPrefix
	out.Region = in.Region
	out.CredentialSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.Endpoint = in.Endpoint
	return nil
}

func Convert_v1alpha1_AzureKeyVault_To_v1alpha2_AzureKeyVault(in *AzureKeyVault, out *v1alpha2.AzureKeyVault, s conversion.Scope) error {
	out.VaultBaseURL = in.VaultBaseURL
	out.Cloud = in.Cloud
	out.TenantID = in.TenantID
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.AADClientSecret,
	}
	out.TLSSecretRef = &core.LocalObjectReference{
		Name: in.ClientCertSecret,
	}
	// WARNING: in.ClientCertSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.AADClientSecret requires manual conversion: does not exist in peer-type
	out.UseManagedIdentity = in.UseManagedIdentity
	return nil
}

func Convert_v1alpha2_AzureKeyVault_To_v1alpha1_AzureKeyVault(in *v1alpha2.AzureKeyVault, out *AzureKeyVault, s conversion.Scope) error {
	out.VaultBaseURL = in.VaultBaseURL
	out.Cloud = in.Cloud
	out.TenantID = in.TenantID
	out.ClientCertSecret = in.TLSSecretRef.Name
	out.AADClientSecret = in.CredentialSecretRef.Name
	// WARNING: in.TLSSecretRef requires manual conversion: does not exist in peer-type
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.UseManagedIdentity = in.UseManagedIdentity
	return nil
}

func Convert_v1alpha1_AzureSpec_To_v1alpha2_AzureSpec(in *AzureSpec, out *v1alpha2.AzureSpec, s conversion.Scope) error {
	out.AccountName = in.AccountName
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.AccountKeySecret,
	}
	// WARNING: in.AccountKeySecret requires manual conversion: does not exist in peer-type
	out.Container = in.Container
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha2_AzureSpec_To_v1alpha1_AzureSpec(in *v1alpha2.AzureSpec, out *AzureSpec, s conversion.Scope) error {
	out.AccountName = in.AccountName
	out.AccountKeySecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.Container = in.Container
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha1_GoogleKmsGcsSpec_To_v1alpha2_GoogleKmsGcsSpec(in *GoogleKmsGcsSpec, out *v1alpha2.GoogleKmsGcsSpec, s conversion.Scope) error {
	out.KmsCryptoKey = in.KmsCryptoKey
	out.KmsKeyRing = in.KmsKeyRing
	out.KmsLocation = in.KmsLocation
	out.KmsProject = in.KmsProject
	out.Bucket = in.Bucket
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_GoogleKmsGcsSpec_To_v1alpha1_GoogleKmsGcsSpec(in *v1alpha2.GoogleKmsGcsSpec, out *GoogleKmsGcsSpec, s conversion.Scope) error {
	out.KmsCryptoKey = in.KmsCryptoKey
	out.KmsKeyRing = in.KmsKeyRing
	out.KmsLocation = in.KmsLocation
	out.KmsProject = in.KmsProject
	out.Bucket = in.Bucket
	out.CredentialSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_GcsSpec_To_v1alpha2_GcsSpec(in *GcsSpec, out *v1alpha2.GcsSpec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.ChunkSize = in.ChunkSize
	out.MaxParallel = in.MaxParallel
	out.HAEnabled = in.HAEnabled
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_GcsSpec_To_v1alpha1_GcsSpec(in *v1alpha2.GcsSpec, out *GcsSpec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.ChunkSize = in.ChunkSize
	out.MaxParallel = in.MaxParallel
	out.HAEnabled = in.HAEnabled
	out.CredentialSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_EtcdSpec_To_v1alpha2_EtcdSpec(in *EtcdSpec, out *v1alpha2.EtcdSpec, s conversion.Scope) error {
	out.Address = in.Address
	out.EtcdApi = in.EtcdApi
	out.HAEnable = in.HAEnable
	out.Path = in.Path
	out.Sync = in.Sync
	out.DiscoverySrv = in.DiscoverySrv
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecretName,
	}
	out.TLSSecretRef = &core.LocalObjectReference{
		Name: in.TLSSecretName,
	}
	// WARNING: in.CredentialSecretName requires manual conversion: does not exist in peer-type
	// WARNING: in.TLSSecretName requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_EtcdSpec_To_v1alpha1_EtcdSpec(in *v1alpha2.EtcdSpec, out *EtcdSpec, s conversion.Scope) error {
	out.Address = in.Address
	out.EtcdApi = in.EtcdApi
	out.HAEnable = in.HAEnable
	out.Path = in.Path
	out.Sync = in.Sync
	out.DiscoverySrv = in.DiscoverySrv
	out.CredentialSecretName = in.CredentialSecretRef.Name
	out.TLSSecretName = in.TLSSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	// WARNING: in.TLSSecretRef requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_DynamoDBSpec_To_v1alpha2_DynamoDBSpec(in *DynamoDBSpec, out *v1alpha2.DynamoDBSpec, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.Region = in.Region
	out.HaEnabled = in.HaEnabled
	out.ReadCapacity = in.ReadCapacity
	out.WriteCapacity = in.WriteCapacity
	out.Table = in.Table
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.SessionTokenSecret requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha2_DynamoDBSpec_To_v1alpha1_DynamoDBSpec(in *v1alpha2.DynamoDBSpec, out *DynamoDBSpec, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.Region = in.Region
	out.HaEnabled = in.HaEnabled
	out.ReadCapacity = in.ReadCapacity
	out.WriteCapacity = in.WriteCapacity
	out.Table = in.Table
	out.CredentialSecret = in.CredentialSecretRef.Name
	out.SessionTokenSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha1_RaftSpec_To_v1alpha2_RaftSpec(in *RaftSpec, out *v1alpha2.RaftSpec, s conversion.Scope) error {
	// WARNING: in.Path requires manual conversion: does not exist in peer-type
	out.PerformanceMultiplier = in.PerformanceMultiplier
	out.TrailingLogs = (*int64)(unsafe.Pointer(in.TrailingLogs))
	out.SnapshotThreshold = (*int64)(unsafe.Pointer(in.SnapshotThreshold))
	out.MaxEntrySize = (*int64)(unsafe.Pointer(in.MaxEntrySize))
	out.AutopilotReconcileInterval = in.AutopilotReconcileInterval
	out.Storage = in.Storage
	return nil
}

func Convert_v1alpha1_SwiftSpec_To_v1alpha2_SwiftSpec(in *SwiftSpec, out *v1alpha2.SwiftSpec, s conversion.Scope) error {
	out.AuthURL = in.AuthURL
	out.Container = in.Container
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	out.Tenant = in.Tenant
	out.Region = in.Region
	out.TenantID = in.TenantID
	out.Domain = in.Domain
	out.ProjectDomain = in.ProjectDomain
	out.TrustID = in.TrustID
	out.StorageURL = in.StorageURL
	// WARNING: in.AuthTokenSecret requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha2_SwiftSpec_To_v1alpha1_SwiftSpec(in *v1alpha2.SwiftSpec, out *SwiftSpec, s conversion.Scope) error {
	out.AuthURL = in.AuthURL
	out.Container = in.Container
	out.CredentialSecret = in.CredentialSecretRef.Name
	out.AuthTokenSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.Tenant = in.Tenant
	out.Region = in.Region
	out.TenantID = in.TenantID
	out.Domain = in.Domain
	out.ProjectDomain = in.ProjectDomain
	out.TrustID = in.TrustID
	out.StorageURL = in.StorageURL
	out.MaxParallel = in.MaxParallel
	return nil
}

func Convert_v1alpha1_S3Spec_To_v1alpha2_S3Spec(in *S3Spec, out *v1alpha2.S3Spec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Endpoint = in.Endpoint
	out.Region = in.Region
	out.CredentialSecretRef = &core.LocalObjectReference{
		Name: in.CredentialSecret,
	}
	// WARNING: in.CredentialSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.SessionTokenSecret requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	out.ForcePathStyle = in.ForcePathStyle
	out.DisableSSL = in.DisableSSL
	return nil
}

func Convert_v1alpha2_S3Spec_To_v1alpha1_S3Spec(in *v1alpha2.S3Spec, out *S3Spec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Endpoint = in.Endpoint
	out.Region = in.Region
	out.CredentialSecret = in.CredentialSecretRef.Name
	out.SessionTokenSecret = in.CredentialSecretRef.Name
	// WARNING: in.CredentialSecretRef requires manual conversion: does not exist in peer-type
	out.MaxParallel = in.MaxParallel
	out.ForcePathStyle = in.ForcePathStyle
	out.DisableSSL = in.DisableSSL
	return nil
}
