
# DtoCredentialOut


## Properties

Name | Type
------------ | -------------
`account_id` | string
`created_at` | string
`id` | string
`kind` | string
`name` | string
`provider` | string
`region` | string
`schema_version` | number
`scope` | object
`scope_kind` | string
`scope_version` | number
`updated_at` | string

## Example

```typescript
import type { DtoCredentialOut } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "account_id": null,
  "created_at": null,
  "id": null,
  "kind": null,
  "name": null,
  "provider": null,
  "region": null,
  "schema_version": null,
  "scope": null,
  "scope_kind": null,
  "scope_version": null,
  "updated_at": null,
} satisfies DtoCredentialOut

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCredentialOut
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


