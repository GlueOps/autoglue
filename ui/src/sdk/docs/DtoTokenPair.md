
# DtoTokenPair


## Properties

Name | Type
------------ | -------------
`access_token` | string
`expires_in` | number
`refresh_token` | string
`token_type` | string

## Example

```typescript
import type { DtoTokenPair } from '@glueops/autoglue-sdk'

// TODO: Update the object below with actual values
const example = {
  "access_token": eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ij...,
  "expires_in": 3600,
  "refresh_token": m0l9o8rT3t0V8d3eFf....,
  "token_type": Bearer,
} satisfies DtoTokenPair

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoTokenPair
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


