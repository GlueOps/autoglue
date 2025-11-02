# DtoJWK

## Properties

| Name  | Type   |
| ----- | ------ |
| `alg` | string |
| `e`   | string |
| `kid` | string |
| `kty` | string |
| `n`   | string |
| `use` | string |
| `x`   | string |

## Example

```typescript
import type { DtoJWK } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "alg": RS256,
  "e": AQAB,
  "kid": 7c6f1d0a-7a98-4e6a-9dbf-6b1af4b9f345,
  "kty": RSA,
  "n": null,
  "use": sig,
  "x": null,
} satisfies DtoJWK

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoJWK
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
