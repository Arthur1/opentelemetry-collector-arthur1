[comment]: <> (Code generated by mdatagen. DO NOT EDIT.)

# runn

## Default Metrics

The following metrics are emitted by default. Each of them can be disabled by applying the following configuration:

```yaml
metrics:
  <metric_name>:
    enabled: false
```

### runn.elapsed_time

elapsed time of step

| Unit | Metric Type | Value Type |
| ---- | ----------- | ---------- |
| s | Gauge | Double |

#### Attributes

| Name | Description | Values |
| ---- | ----------- | ------ |
| runn.runbook.desc | description of runbook | Any Str |
| runn.runbook.step.key | step key in runbook | Any Str |

### runn.status

1 if the operation has finished successfully, otherwise 0.

| Unit | Metric Type | Value Type | Aggregation Temporality | Monotonic |
| ---- | ----------- | ---------- | ----------------------- | --------- |
| 1 | Sum | Int | Cumulative | false |

#### Attributes

| Name | Description | Values |
| ---- | ----------- | ------ |
| runn.runbook.desc | description of runbook | Any Str |
