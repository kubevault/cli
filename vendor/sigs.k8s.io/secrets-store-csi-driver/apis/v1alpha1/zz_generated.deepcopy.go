//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ByPodStatus) DeepCopyInto(out *ByPodStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ByPodStatus.
func (in *ByPodStatus) DeepCopy() *ByPodStatus {
	if in == nil {
		return nil
	}
	out := new(ByPodStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretObject) DeepCopyInto(out *SecretObject) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make([]*SecretObjectData, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(SecretObjectData)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretObject.
func (in *SecretObject) DeepCopy() *SecretObject {
	if in == nil {
		return nil
	}
	out := new(SecretObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretObjectData) DeepCopyInto(out *SecretObjectData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretObjectData.
func (in *SecretObjectData) DeepCopy() *SecretObjectData {
	if in == nil {
		return nil
	}
	out := new(SecretObjectData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClass) DeepCopyInto(out *SecretProviderClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClass.
func (in *SecretProviderClass) DeepCopy() *SecretProviderClass {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretProviderClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassList) DeepCopyInto(out *SecretProviderClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecretProviderClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassList.
func (in *SecretProviderClassList) DeepCopy() *SecretProviderClassList {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretProviderClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassObject) DeepCopyInto(out *SecretProviderClassObject) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassObject.
func (in *SecretProviderClassObject) DeepCopy() *SecretProviderClassObject {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassPodStatus) DeepCopyInto(out *SecretProviderClassPodStatus) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassPodStatus.
func (in *SecretProviderClassPodStatus) DeepCopy() *SecretProviderClassPodStatus {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassPodStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretProviderClassPodStatus) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassPodStatusList) DeepCopyInto(out *SecretProviderClassPodStatusList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecretProviderClassPodStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassPodStatusList.
func (in *SecretProviderClassPodStatusList) DeepCopy() *SecretProviderClassPodStatusList {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassPodStatusList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretProviderClassPodStatusList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassPodStatusStatus) DeepCopyInto(out *SecretProviderClassPodStatusStatus) {
	*out = *in
	if in.Objects != nil {
		in, out := &in.Objects, &out.Objects
		*out = make([]SecretProviderClassObject, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassPodStatusStatus.
func (in *SecretProviderClassPodStatusStatus) DeepCopy() *SecretProviderClassPodStatusStatus {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassPodStatusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassSpec) DeepCopyInto(out *SecretProviderClassSpec) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SecretObjects != nil {
		in, out := &in.SecretObjects, &out.SecretObjects
		*out = make([]*SecretObject, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(SecretObject)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassSpec.
func (in *SecretProviderClassSpec) DeepCopy() *SecretProviderClassSpec {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretProviderClassStatus) DeepCopyInto(out *SecretProviderClassStatus) {
	*out = *in
	if in.ByPod != nil {
		in, out := &in.ByPod, &out.ByPod
		*out = make([]*ByPodStatus, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ByPodStatus)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretProviderClassStatus.
func (in *SecretProviderClassStatus) DeepCopy() *SecretProviderClassStatus {
	if in == nil {
		return nil
	}
	out := new(SecretProviderClassStatus)
	in.DeepCopyInto(out)
	return out
}
