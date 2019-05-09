# aws-ecs-deploy

Simple AWS ECS Task Deployment

## Usage
```
Usage of deploy:
  -cluster string
        Cluster Name
  -container string
        Container Name
  -image string
        Image URL
  -service string
        Service Name
  -skip-verify
        Skip Service Verify
  -task-definition string
        Task Definition Name
  -wait-count int
        Wait Count (default 30)
  -wait-sleep int
        Wait Sleep (default 3)
```

## Environment Variables

| Name                  | Description                          |
|-----------------------|--------------------------------------|
| AWS_REGION            | Region                               |
| AWS_ACCESS_KEY_ID     | Access Key ID                        |
| AWS_SECRET_ACCESS_KEY | Secret Access Key                    |
| IMAGE                 | Image Name                           |
| CONTAINER             | Container Name                       | 
| CLUSTER               | Cluster Name                         |
| SERVICE               | Service Name                         |
| WAIT_SLEEP            | Service verify Interval (sec)        |
| WAIT_COUNT            | Maximum number of verify for service | 
| SKIP_VERIFY           | Skip service verify                  |

## Example

### Google Cloud Build

cloudbuild.yaml
```
steps:
- name: gcr.io/cloud-builders/docker
  args:
  - build
  - --tag=$_IMAGE
  - .
- name: gcr.io/cloud-builders/docker
  args:
  - push
  - $_IMAGE
- name: aplulu/aws-ecs-deploy
  args:
  - --image=$_IMAGE:$COMMIT_SHA
  - --task-definition=[TaskDefinitionName]
  - --container=[ECSTaskDefinitionContainerName]
  - --cluster=[ECSClusterName]
  - --service=[ECSServiceName]
  env:
  - AWS_REGION=ap-northeast-1
  - AWS_ACCESS_KEY_ID=[AccessKeyID]
  secretEnv:
  - AWS_SECRET_ACCESS_KEY
substitutions:
  _IMAGE: [ImageName]
timeout: 720s
secrets:
- kmsKeyName: projects/[ProjectID]/locations/global/keyRings/[KeyringName]/cryptoKeys/[KeyName]
  secretEnv:
    AWS_SECRET_ACCESS_KEY: 
```