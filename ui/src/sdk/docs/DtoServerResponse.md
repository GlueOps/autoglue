
# DtoServerResponse


## Properties

Name | Type
------------ | -------------
`created_at` | string
`hostname` | string
`id` | string
`organization_id` | string
`private_ip_address` | string
`public_ip_address` | string
`role` | string
`ssh_key_id` | string
`ssh_user` | string
`status` | string
`updated_at` | string

## Example

```typescript
import type { DtoServerResponse } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "created_at": null,
  "hostname": null,
  "id": null,
  "organization_id": null,
  "private_ip_address": null,
  "public_ip_address": null,
  "role": master|worker|bastion,
  "ssh_key_id": null,
  "ssh_user": null,
  "status": pending|provisioning|ready|failed,
  "updated_at": null,
} satisfies DtoServerResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoServerResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


