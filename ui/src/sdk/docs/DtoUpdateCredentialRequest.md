
# DtoUpdateCredentialRequest


## Properties

Name | Type
------------ | -------------
`account_id` | string
`name` | string
`region` | string
`scope` | object
`scope_kind` | string
`scope_version` | number
`secret` | object

## Example

```typescript
import type { DtoUpdateCredentialRequest } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "account_id": null,
  "name": null,
  "region": null,
  "scope": null,
  "scope_kind": null,
  "scope_version": null,
  "secret": null,
} satisfies DtoUpdateCredentialRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoUpdateCredentialRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


