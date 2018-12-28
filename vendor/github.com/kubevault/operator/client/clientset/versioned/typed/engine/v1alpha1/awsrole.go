/*
Copyright 2018 The Kube Vault Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubevault/operator/apis/engine/v1alpha1"
	scheme "github.com/kubevault/operator/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AWSRolesGetter has a method to return a AWSRoleInterface.
// A group's client should implement this interface.
type AWSRolesGetter interface {
	AWSRoles(namespace string) AWSRoleInterface
}

// AWSRoleInterface has methods to work with AWSRole resources.
type AWSRoleInterface interface {
	Create(*v1alpha1.AWSRole) (*v1alpha1.AWSRole, error)
	Update(*v1alpha1.AWSRole) (*v1alpha1.AWSRole, error)
	UpdateStatus(*v1alpha1.AWSRole) (*v1alpha1.AWSRole, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AWSRole, error)
	List(opts v1.ListOptions) (*v1alpha1.AWSRoleList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AWSRole, err error)
	AWSRoleExpansion
}

// aWSRoles implements AWSRoleInterface
type aWSRoles struct {
	client rest.Interface
	ns     string
}

// newAWSRoles returns a AWSRoles
func newAWSRoles(c *EngineV1alpha1Client, namespace string) *aWSRoles {
	return &aWSRoles{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the aWSRole, and returns the corresponding aWSRole object, and an error if there is any.
func (c *aWSRoles) Get(name string, options v1.GetOptions) (result *v1alpha1.AWSRole, err error) {
	result = &v1alpha1.AWSRole{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awsroles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AWSRoles that match those selectors.
func (c *aWSRoles) List(opts v1.ListOptions) (result *v1alpha1.AWSRoleList, err error) {
	result = &v1alpha1.AWSRoleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awsroles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aWSRoles.
func (c *aWSRoles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("awsroles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a aWSRole and creates it.  Returns the server's representation of the aWSRole, and an error, if there is any.
func (c *aWSRoles) Create(aWSRole *v1alpha1.AWSRole) (result *v1alpha1.AWSRole, err error) {
	result = &v1alpha1.AWSRole{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("awsroles").
		Body(aWSRole).
		Do().
		Into(result)
	return
}

// Update takes the representation of a aWSRole and updates it. Returns the server's representation of the aWSRole, and an error, if there is any.
func (c *aWSRoles) Update(aWSRole *v1alpha1.AWSRole) (result *v1alpha1.AWSRole, err error) {
	result = &v1alpha1.AWSRole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("awsroles").
		Name(aWSRole.Name).
		Body(aWSRole).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *aWSRoles) UpdateStatus(aWSRole *v1alpha1.AWSRole) (result *v1alpha1.AWSRole, err error) {
	result = &v1alpha1.AWSRole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("awsroles").
		Name(aWSRole.Name).
		SubResource("status").
		Body(aWSRole).
		Do().
		Into(result)
	return
}

// Delete takes name of the aWSRole and deletes it. Returns an error if one occurs.
func (c *aWSRoles) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awsroles").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aWSRoles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awsroles").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched aWSRole.
func (c *aWSRoles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AWSRole, err error) {
	result = &v1alpha1.AWSRole{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("awsroles").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
