package main

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired composed resources from %T", req))
		return rsp, nil
	}

	for name, dr := range desired {
		log := f.log.WithValues("desired-resource-name", name)

		annotation := dr.Resource.GetAnnotations()["fn.crossplane.io/pause-when-ready"]
		if annotation == "true" && dr.Ready == resource.ReadyTrue {
			log.Debug("Detected the resource to pause when ready:")
			setAnnotations := dr.Resource.GetAnnotations()
			deletionTimestamp := dr.Resource.GetDeletionTimestamp()
			if deletionTimestamp == nil {
				setAnnotations["crossplane.io/paused"] = "true"
			} else {
				setAnnotations["crossplane.io/paused"] = "false"
			}
			dr.Resource.SetAnnotations(setAnnotations)
		}
	}

	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources from %T", req))
		return rsp, nil
	}

	return rsp, nil
}
