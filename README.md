# Downsampling InfluxDB Metrics

Downsampling is necessary to reduce the storage burden and search optimization for InfluxDB queries.

- downsample historic data
- write to different cluster than source
- job gets deployed on K8

## Deployment

```
metrics-downsampling-job go/compile       compile go programs
metrics-downsampling-job docker/tags      list the existing tagged images
metrics-downsampling-job docker/build     build and tag the Docker image. vars:tag
metrics-downsampling-job docker/push      push the Docker image to ECR. vars:tag
metrics-downsampling-job helm/install     Deploy the stack into kubernetes. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
metrics-downsampling-job helm/delete      delete stack from reference. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
metrics-downsampling-job helm/reinstall   delete stack from reference and then deploy. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
metrics-downsampling-job deploy           Compiles, builds and deploys a stack for a tag. vars: tag, stack, queryid (e.g. tag=latest stack=test queryid=q37f89d)
metrics-downsampling-job redeploy         Compiles, builds and re-deploys a stack for a tag. vars: tag, stack, queryid (e.g. tag=latest stack=test queryid=q37f89d)
metrics-downsampling-job help             this helps
```
