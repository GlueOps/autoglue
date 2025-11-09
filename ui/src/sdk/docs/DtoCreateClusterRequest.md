
# DtoCreateClusterRequest


## Properties

Name | Type
------------ | -------------
`captain_domain` | string
`cluster_load_balancer` | string
`control_load_balancer` | string
`name` | string
`provider` | string
`region` | string
`status` | string

## Example

```typescript
import type { DtoCreateClusterRequest } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "captain_domain": null,
  "cluster_load_balancer": null,
  "control_load_balancer": null,
  "name": null,
  "provider": null,
  "region": null,
  "status": null,
} satisfies DtoCreateClusterRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateClusterRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


