# Exercises from Programming with Kubernetes (educative.io)

- Webhook authentication and authorization.
- Admission Controller - Validating webhook admissions.
- CRD (Custom Resources Definition ) - Example

Webhook Authentication by token.

| Application | Description |
|-------------|-------------|
| [webhook-setup](lab-01-webhook-auhtN/) | Sample code to create a simple http server (webhook server) that authenticate a user (mock sample)  |
| [webhook-certificates](lab-02-webhook-authN-cert) | https webhook w/certificates  |
| [webook-deploy](lab-03-webhook-authN-deploy/) | Kube deploy for webhook authentication app |


Webhook Authorization by entity or user.

| Application | Description |
|-------------|-------------|
| [webhook-setup](lab-04-webhook-authZ/) | Sample code to create a simple http server (webhook server) that authorizate a user (mock sample)  |
| [webhook-certificates](lab-05-webhook-authZ-cert) | https webhook w/certificates  |
| [webook-deploy](lab-06-webhook-authZ-deploy/) | Kube deploy for webhook authorization app |


Admission Controller - Validating, Mutating webhook admissions

| Application | Description |
|-------------|-------------|
| [admission-webhook-setup](lab-07-admission-webhook/) | Sample code to create an validating-admission-webhook |
| [mutating-admission-webhook](lab-08-mutating-admission-webhook/) | Sample code to create an mutating-admission-webhook |

CRD - Example 

| Application | Description |
|-------------|-------------|
| [crd-setup](lab-09-crd-example/) | Sample code to create an crd object |
