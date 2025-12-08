# DtoClusterResponse

## Properties

| Name                       | Type                                                       |
| -------------------------- | ---------------------------------------------------------- |
| `apps_load_balancer`       | [DtoLoadBalancerResponse](DtoLoadBalancerResponse.md)      |
| `bastion_server`           | [DtoServerResponse](DtoServerResponse.md)                  |
| `captain_domain`           | [DtoDomainResponse](DtoDomainResponse.md)                  |
| `certificate_key`          | string                                                     |
| `cluster_provider`         | string                                                     |
| `control_plane_fqdn`       | string                                                     |
| `control_plane_record_set` | [DtoRecordSetResponse](DtoRecordSetResponse.md)            |
| `created_at`               | string                                                     |
| `docker_image`             | string                                                     |
| `docker_tag`               | string                                                     |
| `glueops_load_balancer`    | [DtoLoadBalancerResponse](DtoLoadBalancerResponse.md)      |
| `id`                       | string                                                     |
| `last_error`               | string                                                     |
| `name`                     | string                                                     |
| `node_pools`               | [Array&lt;DtoNodePoolResponse&gt;](DtoNodePoolResponse.md) |
| `random_token`             | string                                                     |
| `region`                   | string                                                     |
| `status`                   | string                                                     |
| `updated_at`               | string                                                     |

## Example

```typescript
import type { DtoClusterResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  apps_load_balancer: null,
  bastion_server: null,
  captain_domain: null,
  certificate_key: null,
  cluster_provider: null,
  control_plane_fqdn: null,
  control_plane_record_set: null,
  created_at: null,
  docker_image: null,
  docker_tag: null,
  glueops_load_balancer: null,
  id: null,
  last_error: null,
  name: null,
  node_pools: null,
  random_token: null,
  region: null,
  status: null,
  updated_at: null,
} satisfies DtoClusterResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoClusterResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
