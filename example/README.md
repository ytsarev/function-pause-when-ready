# Example manifests

You can run your function locally and test it using `crossplane beta render`
with these example manifests.

```shell
# Run the function locally
$ go run . --insecure --debug
```

```shell
# Then, in another terminal, call it with these example manifests
$ crossplane beta render example/xr.yaml example/composition.yaml example/functions.yaml -o example/observed.yaml -r
---
apiVersion: test.fn.upbound.io/v1alpha1
kind: AnyXR
metadata:
  name: any-xr
---
apiVersion: s3.aws.upbound.io/v1beta1
kind: Bucket
metadata:
  annotations:
    crossplane.io/composition-resource-name: bucket
    crossplane.io/paused: "true"
    fn.crossplane.io/pause-when-ready: "true"
  generateName: any-xr-
  labels:
    crossplane.io/composite: any-xr
  ownerReferences:
  - apiVersion: test.fn.upbound.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: AnyXR
    name: any-xr
    uid: ""
spec:
  forProvider:
    region: us-east-2
```
