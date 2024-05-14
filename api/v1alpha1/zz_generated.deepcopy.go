//go:build !ignore_autogenerated

/*
Copyright 2024 FastForge Inc. support@fastforge.com.

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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsumeSpec) DeepCopyInto(out *ConsumeSpec) {
	*out = *in
	if in.Maps != nil {
		in, out := &in.Maps, &out.Maps
		*out = make(KeyMap, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsumeSpec.
func (in *ConsumeSpec) DeepCopy() *ConsumeSpec {
	if in == nil {
		return nil
	}
	out := new(ConsumeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaemonSet) DeepCopyInto(out *DaemonSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaemonSet.
func (in *DaemonSet) DeepCopy() *DaemonSet {
	if in == nil {
		return nil
	}
	out := new(DaemonSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DaemonSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaemonSetList) DeepCopyInto(out *DaemonSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DaemonSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaemonSetList.
func (in *DaemonSetList) DeepCopy() *DaemonSetList {
	if in == nil {
		return nil
	}
	out := new(DaemonSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DaemonSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaemonSetSpec) DeepCopyInto(out *DaemonSetSpec) {
	*out = *in
	in.DaemonSetSpec.DeepCopyInto(&out.DaemonSetSpec)
	in.WorkloadSpec.DeepCopyInto(&out.WorkloadSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaemonSetSpec.
func (in *DaemonSetSpec) DeepCopy() *DaemonSetSpec {
	if in == nil {
		return nil
	}
	out := new(DaemonSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaemonSetStatus) DeepCopyInto(out *DaemonSetStatus) {
	*out = *in
	in.DaemonSetStatus.DeepCopyInto(&out.DaemonSetStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaemonSetStatus.
func (in *DaemonSetStatus) DeepCopy() *DaemonSetStatus {
	if in == nil {
		return nil
	}
	out := new(DaemonSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Deployment) DeepCopyInto(out *Deployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Deployment.
func (in *Deployment) DeepCopy() *Deployment {
	if in == nil {
		return nil
	}
	out := new(Deployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Deployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentList) DeepCopyInto(out *DeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Deployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentList.
func (in *DeploymentList) DeepCopy() *DeploymentList {
	if in == nil {
		return nil
	}
	out := new(DeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentSpec) DeepCopyInto(out *DeploymentSpec) {
	*out = *in
	in.DeploymentSpec.DeepCopyInto(&out.DeploymentSpec)
	in.WorkloadSpec.DeepCopyInto(&out.WorkloadSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentSpec.
func (in *DeploymentSpec) DeepCopy() *DeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(DeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentStatus) DeepCopyInto(out *DeploymentStatus) {
	*out = *in
	in.DeploymentStatus.DeepCopyInto(&out.DeploymentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentStatus.
func (in *DeploymentStatus) DeepCopy() *DeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(DeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in KeyMap) DeepCopyInto(out *KeyMap) {
	{
		in := &in
		*out = make(KeyMap, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyMap.
func (in KeyMap) DeepCopy() KeyMap {
	if in == nil {
		return nil
	}
	out := new(KeyMap)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProduceFromKindSpec) DeepCopyInto(out *ProduceFromKindSpec) {
	*out = *in
	if in.Kind != nil {
		in, out := &in.Kind, &out.Kind
		*out = new(string)
		**out = **in
	}
	if in.Expression != nil {
		in, out := &in.Expression, &out.Expression
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProduceFromKindSpec.
func (in *ProduceFromKindSpec) DeepCopy() *ProduceFromKindSpec {
	if in == nil {
		return nil
	}
	out := new(ProduceFromKindSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProduceFromValueSpec) DeepCopyInto(out *ProduceFromValueSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProduceFromValueSpec.
func (in *ProduceFromValueSpec) DeepCopy() *ProduceFromValueSpec {
	if in == nil {
		return nil
	}
	out := new(ProduceFromValueSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProduceSpec) DeepCopyInto(out *ProduceSpec) {
	*out = *in
	in.ProduceFromKindSpec.DeepCopyInto(&out.ProduceFromKindSpec)
	out.ProduceFromValueSpec = in.ProduceFromValueSpec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProduceSpec.
func (in *ProduceSpec) DeepCopy() *ProduceSpec {
	if in == nil {
		return nil
	}
	out := new(ProduceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulSet) DeepCopyInto(out *StatefulSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulSet.
func (in *StatefulSet) DeepCopy() *StatefulSet {
	if in == nil {
		return nil
	}
	out := new(StatefulSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StatefulSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulSetList) DeepCopyInto(out *StatefulSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StatefulSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulSetList.
func (in *StatefulSetList) DeepCopy() *StatefulSetList {
	if in == nil {
		return nil
	}
	out := new(StatefulSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StatefulSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulSetSpec) DeepCopyInto(out *StatefulSetSpec) {
	*out = *in
	in.StatefulSetSpec.DeepCopyInto(&out.StatefulSetSpec)
	in.WorkloadSpec.DeepCopyInto(&out.WorkloadSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulSetSpec.
func (in *StatefulSetSpec) DeepCopy() *StatefulSetSpec {
	if in == nil {
		return nil
	}
	out := new(StatefulSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulSetStatus) DeepCopyInto(out *StatefulSetStatus) {
	*out = *in
	in.StatefulSetStatus.DeepCopyInto(&out.StatefulSetStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulSetStatus.
func (in *StatefulSetStatus) DeepCopy() *StatefulSetStatus {
	if in == nil {
		return nil
	}
	out := new(StatefulSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadSpec) DeepCopyInto(out *WorkloadSpec) {
	*out = *in
	if in.Delegates != nil {
		in, out := &in.Delegates, &out.Delegates
		*out = make([]Delegate, len(*in))
		copy(*out, *in)
	}
	if in.Consumes != nil {
		in, out := &in.Consumes, &out.Consumes
		*out = make(map[WorkloadName]ConsumeSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Produces != nil {
		in, out := &in.Produces, &out.Produces
		*out = make(map[Key]ProduceSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadSpec.
func (in *WorkloadSpec) DeepCopy() *WorkloadSpec {
	if in == nil {
		return nil
	}
	out := new(WorkloadSpec)
	in.DeepCopyInto(out)
	return out
}
