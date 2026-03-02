
This document assumes you've reviewed the Quick Start guide.

## Existing Secrets

To use a pre-existing secrets for your credentials, first create the credentials secret:

```bash
cat << EOF | kubectl create -f -
apiVersion: v1
kind: Secret
metadata:
    name: my-cosi-credentials
data:
    S3_ACCESS_KEY: myAccessKey
    S3_SECRET_KEY: mySecretAccessKey
EOF
```

Then, when installing the COSI driver, pass the name in via the values.yaml:
```
existingCredentialsSecret: 'my-cosi-credentials'
```

## Self-signed Certificates

Create a secret containing any self-signed certificates you wish to use:
```bash
kubectl create secret generic self-signed-certs --from-file /tmp/s3-region-1.pem --from-file /tmp/iam.pem
```

Then pass this in via the values.yaml:
```
existingS3CertificateSecret: 'self-signed-certs'
```
