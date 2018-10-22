# authx
Nalej Authentication with support for JWT

## Cluster management requirements

To deploy Authx in a Kubernetes cluster, it requires that the cluster contains a specific secret.

```
apiVersion: v1
kind: Secret
metadata:
  name: authx-secret
  namespace: nalej
type: Opaque
data:
  secret: [your_secret]
```