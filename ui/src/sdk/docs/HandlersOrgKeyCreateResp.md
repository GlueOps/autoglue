
# HandlersOrgKeyCreateResp


## Properties

Name | Type
------------ | -------------
`created_at` | string
`expires_at` | string
`id` | string
`name` | string
`org_key` | string
`org_secret` | string
`scope` | string

## Example

```typescript
import type { HandlersOrgKeyCreateResp } from '@glueops/autoglue-sdk'

// TODO: Update the object below with actual values
const example = {
  "created_at": null,
  "expires_at": null,
  "id": null,
  "name": null,
  "org_key": null,
  "org_secret": null,
  "scope": null,
} satisfies HandlersOrgKeyCreateResp

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersOrgKeyCreateResp
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


