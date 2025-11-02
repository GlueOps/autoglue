# DtoCreateServerRequest

## Properties

| Name                 | Type   |
| -------------------- | ------ |
| `hostname`           | string |
| `private_ip_address` | string |
| `public_ip_address`  | string |
| `role`               | string |
| `ssh_key_id`         | string |
| `ssh_user`           | string |
| `status`             | string |

## Example

```typescript
import type { DtoCreateServerRequest } from "@glueops/autoglue-sdk";

// TODO: Update the object below with actual values
const example = {
  hostname: null,
  private_ip_address: null,
  public_ip_address: null,
  role: master | worker | bastion,
  ssh_key_id: null,
  ssh_user: null,
  status: pending | provisioning | ready | failed,
} satisfies DtoCreateServerRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateServerRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
