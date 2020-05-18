/*
Copyright The KubeVault Authors.

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

package util

import (
	"encoding/json"
	"fmt"

	api "kubevault.dev/operator/apis/engine/v1alpha1"
	cs "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchAWSAccessKeyRequest(c cs.EngineV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.AWSAccessKeyRequest) *api.AWSAccessKeyRequest) (*api.AWSAccessKeyRequest, kutil.VerbType, error) {
	cur, err := c.AWSAccessKeyRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating AWSAccessKeyRequest %s/%s.", meta.Namespace, meta.Name)
		out, err := c.AWSAccessKeyRequests(meta.Namespace).Create(transform(&api.AWSAccessKeyRequest{
			TypeMeta: metav1.TypeMeta{
				Kind:       api.ResourceKindAWSAccessKeyRequest,
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchAWSAccessKeyRequest(c, cur, transform)
}

func PatchAWSAccessKeyRequest(c cs.EngineV1alpha1Interface, cur *api.AWSAccessKeyRequest, transform func(*api.AWSAccessKeyRequest) *api.AWSAccessKeyRequest) (*api.AWSAccessKeyRequest, kutil.VerbType, error) {
	return PatchAWSAccessKeyRequestObject(c, cur, transform(cur.DeepCopy()))
}

func PatchAWSAccessKeyRequestObject(c cs.EngineV1alpha1Interface, cur, mod *api.AWSAccessKeyRequest) (*api.AWSAccessKeyRequest, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching AWSAccessKeyRequest %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.AWSAccessKeyRequests(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateAWSAccessKeyRequest(c cs.EngineV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.AWSAccessKeyRequest) *api.AWSAccessKeyRequest) (result *api.AWSAccessKeyRequest, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.AWSAccessKeyRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.AWSAccessKeyRequests(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update AWSAccessKeyRequest %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update AWSAccessKeyRequest %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdateAWSAccessKeyRequestStatus(
	c cs.EngineV1alpha1Interface,
	meta metav1.ObjectMeta,
	transform func(*api.AWSAccessKeyRequestStatus) *api.AWSAccessKeyRequestStatus,
) (result *api.AWSAccessKeyRequest, err error) {
	apply := func(x *api.AWSAccessKeyRequest) *api.AWSAccessKeyRequest {
		return &api.AWSAccessKeyRequest{
			TypeMeta:   x.TypeMeta,
			ObjectMeta: x.ObjectMeta,
			Spec:       x.Spec,
			Status:     *transform(x.Status.DeepCopy()),
		}
	}

	attempt := 0
	cur, err := c.AWSAccessKeyRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		var e2 error
		result, e2 = c.AWSAccessKeyRequests(meta.Namespace).UpdateStatus(apply(cur))
		if kerr.IsConflict(e2) {
			latest, e3 := c.AWSAccessKeyRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			switch {
			case e3 == nil:
				cur = latest
				return false, nil
			case kutil.IsRequestRetryable(e3):
				return false, nil
			default:
				return false, e3
			}
		} else if err != nil && !kutil.IsRequestRetryable(e2) {
			return false, e2
		}
		return e2 == nil, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update status of AWSAccessKeyRequest %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
