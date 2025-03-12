
# Getting Started

## Dependencies

`kubectl` - v1.23+ recommended

`helm` - v3.5+ recommended

## Install

1. Apply manifests for COSI API and Controller
```
kubectl apply -k github.com/kubernetes-sigs/container-object-storage-interface-api
kubectl apply -k github.com/kubernetes-sigs/container-object-storage-interface-controller
```

2. Clone this repository.

3. Edit [values.yaml](../helm/cloudian-cosi/values.yaml) to reflect your Hyperstore configuration.

4. Run helm install, e.g. from the root of this repository:
```
helm install cloudian-cosi helm/cloudian-cosi -f helm/cloudian-cosi/values.yaml
```

## Create a Bucket

### Create a BucketClass

Create a yaml file containing your preferred bucket class configuration, e.g. `bucketclass.yaml`:
```
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketClass
metadata:
  name: cloudian-bucketclass
driverName: cloudian-cosi-driver
deletionPolicy: Delete
```

Apply the BucketClass Resource:
```
kubectl apply -f bucketclass.yaml
```

### Create a BucketClaim

Create a yaml file with your preferred bucket claim configuration, e.g. `bucketclaim.yaml`:
```
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketClaim
metadata:
  name: cloudian-bucketclaim
spec:
  bucketClassName: cloudian-bucketclass
  protocols:
    - S3
```
Apply the BucketClaim Resource:
```
kubectl apply -f bucketclaim.yaml
```

A `bucket` kubernetes resource object will then be created as well as an actual bucket in the HyperStore backend.

## Delete a Bucket

To delete a bucket, simply remove the corresponding BucketClaim, e.g:
```
kubectl delete bucketclaim cloudian-bucketclaim
```
This will delete the `bucket` kubernetes resource object. Whether the actual bucket in HyperStore is deleted depends
on the `deletionPolicy` set in the BucketClass.

## Grant Bucket Access

### Create a BucketAccessClass

Create a yaml file containing your preferred bucket access class configuration, e.g. `bucketaccessclass.yaml`:
```
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketAccessClass
metadata:
  name: cloudian-bucketaccessclass                                  
driverName: cloudian-cosi-driver           
authenticationType: IAM
```
Apply the BucketAccessClass Resource:
```
kubectl apply -f bucketaccessclass.yaml
```

### Create a BucketAccess

Create a yaml file with your preferred bucket access configuration, e.g. `bucketaccess.yaml`:
```
apiVersion: cosi.storage.k8s.io/v1alpha1
kind: BucketAccess
metadata:
  name: cloudian-bucketaccess
spec:
  bucketClaimName: cloudian-bucketclaim
  bucketAccessClassName: cloudian-bucketaccessclass
  credentialsSecretName: cloudian-secret
  serviceAccountName: cloudian-serviceaccount
```
Apply the BucketAccess Resource:
```
kubectl apply -f bucketaccess.yaml
```

Then depending on the `authenticationType` specified in the BucketAccessClass, access the bucket will be granted by
either IAM or keys and the credentials will be stored in a secret with the name specified in the BucketAccess.

## Revoke Bucket Access

To revoke access to a bucket, simply remove the corresponding BucketAccess, e.g:
```
kubectl delete bucketclaim cloudian-bucketclaim
```
This will delete the secret containing the credentials to access the bucket, and ensure that those credentials are
no longer valid.
