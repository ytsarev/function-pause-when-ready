package main

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"PauseWhenReady": {
			reason: "Expected paused annotation to be set when resource is ready",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "test"},
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
							},
						},
					},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "test", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true",
										   "crossplane.io/paused": "true"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
		"UnpauseWhenDeleting": {
			reason: "Expected paused annotation to be set to false when resource is getting deleted",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "test"},
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr",
									"deletionTimestamp": "2023-12-04T23:58:19Z"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
							},
						},
					},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "test", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed",
										"annotations": {
										   "fn.crossplane.io/pause-when-ready": "true",
										   "crossplane.io/paused": "false"
										}
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
