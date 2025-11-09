
# DtoClusterResponse


## Properties

Name | Type
------------ | -------------
`bastion_server` | [DtoServerResponse](DtoServerResponse.md)
`captain_domain` | string
`certificate_key` | string
`cluster_load_balancer` | string
`control_load_balancer` | string
`created_at` | string
`id` | string
`name` | string
`node_pools` | [Array&lt;DtoNodePoolResponse&gt;](DtoNodePoolResponse.md)
`provider` | string
`random_token` | string
`region` | string
`status` | string
`updated_at` | string

## Example

```typescript
import type { DtoClusterResponse } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "bastion_server": null,
  "captain_domain": null,
  "certificate_key": null,
  "cluster_load_balancer": null,
  "control_load_balancer": null,
  "created_at": null,
  "id": null,
  "name": null,
  "node_pools": null,
  "provider": null,
  "random_token": null,
  "region": null,
  "status": null,
  "updated_at": null,
} satisfies DtoClusterResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoClusterResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


