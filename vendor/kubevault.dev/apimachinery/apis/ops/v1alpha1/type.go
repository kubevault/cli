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

// List of possible condition types for a ops request
const (
	AccessApproved               = "Approved"
	AccessDenied                 = "Denied"
	Failed                       = "Failed"
	NodeCreated                  = "NodeCreated"
	NodeDeleted                  = "NodeDeleted"
	NodeRestarted                = "NodeRestarted"
	PauseVaultServer             = "PauseVaultServer"
	Progressing                  = "Progressing"
	ResumeVaultServer            = "ResumeVaultServer"
	Successful                   = "Successful"
	Updating                     = "Updating"
	UpdateStatefulSets           = "UpdateStatefulSets"
	RestartNodes                 = "RestartNodes"
	TLSRemoved                   = "TLSRemoved"
	TLSAdded                     = "TLSAdded"
	TLSChanged                   = "TLSChanged"
	IssuingConditionUpdated      = "IssuingConditionUpdated"
	CertificateIssuingSuccessful = "CertificateIssuingSuccessful"
	TLSEnabling                  = "TLSEnabling"
	Restart                      = "Restart"
	RestartStatefulSet           = "RestartStatefulSet"
	CertificateSynced            = "CertificateSynced"
	Reconciled                   = "Reconciled"
	RestartStatefulSetPods       = "RestartStatefulSetPods"
)
