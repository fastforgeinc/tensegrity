# Tensegrity - K8s dependency and configuration orchestrator

## What is Tensegrity?
Tensegrity is Kubernetes controller and set of CRDs that allows to create Kubernetes native workloads such as
Deployments, StatefulSets and DaemonSets with dependencies between them by defining `produced` and `consumed` 
configuration keys and values. Tensegrity watches for those key and value changes,
and reconciles workloads if necessary to apply the new configuration.

In addition, Tensegrity allows to specify `delegates` to consume keys and values from other Namespaces
to build more complex development and deployment scenarios.

## Sneak peek

```yaml
apiVersion: k8s.tensegrity.fastforge.io/v1alpha1
kind: Deployment
metadata:
  name: api
spec:
  # native K8s Deployment spec here in addition to...
  delegates:
    # resolves in order production > staging > user-alice
    - kind: Namespace
      name: production
    - kind: Namespace
      name: staging
    - kind: Namespace
      name: user-alice
  consumes:
    # the Deployment consumes keys/values from postgres StatefulSet and maps them to env variables
    - apiVersion: k8s.tensegrity.fastforge.io/v1alpha1
      kind: deployment
      name: postgres
      maps:
        DATABASE_HOST: host
        DATABASE_PORT: port
    # the Deployment consumes keys/values from redis Deployment and maps them to env variables
    - apiVersion: k8s.tensegrity.fastforge.io/v1alpha1
      kind: Deployment
      name: redis
      maps:
        REDIS_HOST: host
        REDIS_PORT: port
  produces:
    # the Deployment produces its own host and port keys from related Ingress and Service
    - key: host
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      fieldPath: '{ .spec.rules[0].host }.{ .metadata.namespace }'
    - key: port
      apiVersion: v1
      kind: Service
      fieldPath: '{ .spec.ports[?(@.name=="http")].port }'
```

## Getting Started
Tensegrity relies on the cert-manager component to function correctly within the cluster. Specifically, it leverages the mutating webhook configuration and validation webhook configuration provided by cert-manager to ensure seamless integration and operation. However, _this dependency is planned for deprecation in future releases_, with efforts underway to eliminate the reliance on cert-manager, enabling a more streamlined and independent deployment process.

**Optional:**
If cert-manager is not already installed in your Kubernetes cluster, you can install it using the following command:
```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
```

This command applies the necessary configuration files to set up cert-manager in your cluster. Please ensure that your cluster meets the prerequisites for cert-manager installation and that you adjust the version (v1.12.0 in this example) as needed for compatibility with your environment.

### Installation Prerequisites
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- cert-manager v1.7+

**Tensegrity static install**

Install to the `tensegrity` namespace:
```shell
kubectl apply -f https://github.com/fastforgeinc/tensegrity/releases/latest/download/install.yaml
```

Follow the full getting started guide to walk through creating and then updating a Tensegrity objects.

### Development Prerequisites
- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=docker.io/<username>/tensegrity:latest
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=docker.io/<username>/tensegrity:latest
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the manifests/samples:

```sh
kubectl apply -k manifests/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k manifests/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=docker.io/<username>/tensegrity:latest
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://github.com/fastforgeinc/tensegrity/releases/latest/download/install.yaml
```

## Contributing
// TODO(user): TBD

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge, Inc.

Tensegrity is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see http://www.gnu.org/licenses/.
