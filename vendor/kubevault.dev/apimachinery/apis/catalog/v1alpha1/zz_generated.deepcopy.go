//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersion) DeepCopyInto(out *VaultServerVersion) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersion.
func (in *VaultServerVersion) DeepCopy() *VaultServerVersion {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VaultServerVersion) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionExporter) DeepCopyInto(out *VaultServerVersionExporter) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionExporter.
func (in *VaultServerVersionExporter) DeepCopy() *VaultServerVersionExporter {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionExporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionInitContainer) DeepCopyInto(out *VaultServerVersionInitContainer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionInitContainer.
func (in *VaultServerVersionInitContainer) DeepCopy() *VaultServerVersionInitContainer {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionInitContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionList) DeepCopyInto(out *VaultServerVersionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VaultServerVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionList.
func (in *VaultServerVersionList) DeepCopy() *VaultServerVersionList {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VaultServerVersionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionSpec) DeepCopyInto(out *VaultServerVersionSpec) {
	*out = *in
	out.Vault = in.Vault
	out.Unsealer = in.Unsealer
	out.InitContainer = in.InitContainer
	out.Exporter = in.Exporter
	in.Stash.DeepCopyInto(&out.Stash)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionSpec.
func (in *VaultServerVersionSpec) DeepCopy() *VaultServerVersionSpec {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionUnsealer) DeepCopyInto(out *VaultServerVersionUnsealer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionUnsealer.
func (in *VaultServerVersionUnsealer) DeepCopy() *VaultServerVersionUnsealer {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionUnsealer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerVersionVault) DeepCopyInto(out *VaultServerVersionVault) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerVersionVault.
func (in *VaultServerVersionVault) DeepCopy() *VaultServerVersionVault {
	if in == nil {
		return nil
	}
	out := new(VaultServerVersionVault)
	in.DeepCopyInto(out)
	return out
}
