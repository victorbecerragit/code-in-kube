# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create an crd object with kubebuilder.



```text


🚀 What Are CRDs in Kubernetes?

Custom Resource Definitions (CRDs) let you extend Kubernetes by defining your own resource types—like Website, DatabaseBackup, or KafkaCluster. Once defined, these behave like native Kubernetes objects (e.g., Pod, Service) and can be managed using kubectl.

🛠️ What Does “Scaffold CRDs” Mean?

Scaffolding CRDs means automatically generating the YAML files, Go code, and controller logic needed to create and manage custom resources. This is often done using frameworks like:

| Tool | Purpose | 
| Kubebuilder | Generates CRD definitions, controllers, and webhooks using Go | 
| Operator SDK | Helps build Kubernetes Operators with CRDs using Go, Helm, or Ansible | 
| Skaffold | Manages the build/deploy workflow, but not CRD scaffolding directly | 


So while Skaffold itself doesn’t scaffold CRDs, it’s often used alongside Kubebuilder or Operator SDK to develop and deploy CRD-based applications


CRD schema#

```yaml

apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name must be in the form: <plural>.<group>
  name: <name>
spec:
  group: <group name>
  conversion: #optional
    # Specifies how custom resources are converted between versions
    # can be None or Webhook
    strategy: None
  names: # Specify the resource and kind names for the custom resource
    categories: # optional
    # List of categories this custom resource belongs to (e.g. 'all')
    - <mycategory>
    kind: <Uppercase name>
    listKind: <Uppercase list kind, defaulted to be kindList>
    plural: <lowercase plural name>
    shortNames: # optional
    # List of strings as short names
    - <alias1>
    singular: <lowercase singular name, defaulted to be lowercase kind>
  scope: Namespaced  # Namespaced or cluster scope
  versions: # List of all API versions
  - name: v1alpha1
    schema: # Optional
      openAPIV3Schema: # OpenAPI v3 schema to use for validation and pruning
        description: HelmChart is the Schema for the helm chart
        properties: # Describe all the fields
          ...
        required: # Mark required fields
        - ...
        type: object
    served: true
    storage: true
    subresources: # Optional
      status: {} # To enable the status subresource (optional)
      scale: # Optional
        specReplicasPath: <JSON path for the replica field, such as `spec.replicas`>
        statusReplicasPath: <JSON path for the replica number in the status>
        labelSelectorPath: <JSON path that corresponds to Scale `status.selector`>
    additionalPrinterColumns: # Optional
    # Specify additional columns returned in Table output. Used by kubectl
    - description: The phase of this custom resource # Example
      jsonPath: .status.phase
      name: STATUS
      type: string
    - jsonPath: .metadata.creationTimestamp # Example
      name: AGE
      type: date
```

📦 Example Workflow
- Use Kubebuilder to scaffold a CRD and controller

```bash
kubebuilder init --domain example.com --repo github.com/example/my-operator
kubebuilder create api --group web --version v1 --kind Website

```

🧰 Step 1: Install Kubebuilder
Make sure you have:
- Go (v1.20+ recommended)
- Kubebuilder installed:

```bash
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/

```

Create a scaffold API#
Below is the development environment in which we can create a scaffold project for CRD

- foo_types.go
....
        // Foo is the Schema for the foos API
        type Foo struct {
            metav1.TypeMeta   `json:",inline"`
            metav1.ObjectMeta `json:"metadata,omitempty"`
            Spec   FooSpec   `json:"spec,omitempty"`
            Status FooStatus `json:"status,omitempty"`
        }

Create a project#
First, we create a new folder and initialize the project. We use educative.io as the domain name.

```bash
mkdir -p ./projects/pwk
cd ./projects/pwk
kubebuilder init --domain educative.io --repo educative.io/pwk
```

Create an API#

```bash
kubebuilder create api --group apps --version v1beta1 --kind Foo --controller=false --make --resource
```

Now, we can see the Foo struct in api/v1beta1/foo_types.go, and we can insert more fields there. After we finish our new API, let’s run the following commands in the terminal above to automatically generate the manifests (e.g. CRDs, CRs, etc). Let’s see the magic!

```bash
make manifests
```

The CRD manifest is already generated at config/crd/bases/apps.educative.io_foos.yaml. Below is the scaffold project view, where we can see lots of template files that would help us easily extend Kubernetes APIs with CRDs.

```bash
|-- Dockerfile
|-- Makefile
|-- PROJECT
|-- README.md
|-- api
|   `-- v1beta1
|       |-- foo_types.go
|       |-- groupversion_info.go
|       `-- zz_generated.deepcopy.go
|-- bin
|   `-- controller-gen
|-- config
|   |-- crd
|   |   |-- kustomization.yaml
|   |   |-- kustomizeconfig.yaml
|   |   `-- patches
|   |       |-- cainjection_in_foos.yaml
|   |       `-- webhook_in_foos.yaml
|   |-- default
|   |   |-- kustomization.yaml
|   |   |-- manager_auth_proxy_patch.yaml
|   |   `-- manager_config_patch.yaml
|   |-- manager
|   |   |-- controller_manager_config.yaml
|   |   |-- kustomization.yaml
|   |   `-- manager.yaml
|   |-- prometheus
|   |   |-- kustomization.yaml
|   |   `-- monitor.yaml
|   |-- rbac
|   |   |-- auth_proxy_client_clusterrole.yaml
|   |   |-- auth_proxy_role.yaml
|   |   |-- auth_proxy_role_binding.yaml
|   |   |-- auth_proxy_service.yaml
|   |   |-- foo_editor_role.yaml
|   |   |-- foo_viewer_role.yaml
|   |   |-- kustomization.yaml
|   |   |-- leader_election_role.yaml
|   |   |-- leader_election_role_binding.yaml
|   |   |-- role_binding.yaml
|   |   `-- service_account.yaml
|   `-- samples
|       `-- apps_v1beta1_foo.yaml
|-- go.mod
|-- go.sum
|-- hack
|   `-- boilerplate.go.txt
`-- main.go

12 directories, 36 files

```

Test it out#

```bash
make install

kubectl get crds

kubectl apply -f config/samples/

$ kubectl get foo -A
