/*
Copyright The Kubernetes Authors.

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
	"context"
	"time"

	v1alpha1 "github.com/muffin-rice/scheduler-plugins/pkg/apis/scheduling/v1alpha1"
	scheme "github.com/muffin-rice/scheduler-plugins/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ElasticQuotasGetter has a method to return a ElasticQuotaInterface.
// A group's client should implement this interface.
type ElasticQuotasGetter interface {
	ElasticQuotas(namespace string) ElasticQuotaInterface
}

// ElasticQuotaInterface has methods to work with ElasticQuota resources.
type ElasticQuotaInterface interface {
	Create(ctx context.Context, elasticQuota *v1alpha1.ElasticQuota, opts v1.CreateOptions) (*v1alpha1.ElasticQuota, error)
	Update(ctx context.Context, elasticQuota *v1alpha1.ElasticQuota, opts v1.UpdateOptions) (*v1alpha1.ElasticQuota, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ElasticQuota, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ElasticQuotaList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ElasticQuota, err error)
	ElasticQuotaExpansion
}

// elasticQuotas implements ElasticQuotaInterface
type elasticQuotas struct {
	client rest.Interface
	ns     string
}

// newElasticQuotas returns a ElasticQuotas
func newElasticQuotas(c *SchedulingV1alpha1Client, namespace string) *elasticQuotas {
	return &elasticQuotas{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the elasticQuota, and returns the corresponding elasticQuota object, and an error if there is any.
func (c *elasticQuotas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ElasticQuota, err error) {
	result = &v1alpha1.ElasticQuota{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticquotas").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ElasticQuotas that match those selectors.
func (c *elasticQuotas) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ElasticQuotaList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ElasticQuotaList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested elasticQuotas.
func (c *elasticQuotas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("elasticquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a elasticQuota and creates it.  Returns the server's representation of the elasticQuota, and an error, if there is any.
func (c *elasticQuotas) Create(ctx context.Context, elasticQuota *v1alpha1.ElasticQuota, opts v1.CreateOptions) (result *v1alpha1.ElasticQuota, err error) {
	result = &v1alpha1.ElasticQuota{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("elasticquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(elasticQuota).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a elasticQuota and updates it. Returns the server's representation of the elasticQuota, and an error, if there is any.
func (c *elasticQuotas) Update(ctx context.Context, elasticQuota *v1alpha1.ElasticQuota, opts v1.UpdateOptions) (result *v1alpha1.ElasticQuota, err error) {
	result = &v1alpha1.ElasticQuota{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("elasticquotas").
		Name(elasticQuota.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(elasticQuota).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the elasticQuota and deletes it. Returns an error if one occurs.
func (c *elasticQuotas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticquotas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *elasticQuotas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticquotas").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched elasticQuota.
func (c *elasticQuotas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ElasticQuota, err error) {
	result = &v1alpha1.ElasticQuota{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("elasticquotas").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
