# @glueops/autoglue-sdk@0.1.0

A TypeScript SDK client for the localhost API.

## Usage

First, install the SDK from npm.

```bash
npm install @glueops/autoglue-sdk --save
```

Next, try it out.

```ts
import { Configuration, AuthApi } from "@glueops/autoglue-sdk";
import type { AuthCallbackRequest } from "@glueops/autoglue-sdk";

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // string | google|github
    provider: provider_example,
  } satisfies AuthCallbackRequest;

  try {
    const data = await api.authCallback(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

## Documentation

### API Endpoints

All URIs are relative to _http://localhost:8080/api/v1_

| Class          | Method                                                        | HTTP request                            | Description                                     |
| -------------- | ------------------------------------------------------------- | --------------------------------------- | ----------------------------------------------- |
| _AuthApi_      | [**authCallback**](docs/AuthApi.md#authcallback)              | **GET** /auth/{provider}/callback       | Handle social login callback                    |
| _AuthApi_      | [**authStart**](docs/AuthApi.md#authstart)                    | **POST** /auth/{provider}/start         | Begin social login                              |
| _AuthApi_      | [**getJWKS**](docs/AuthApi.md#getjwks)                        | **GET** /.well-known/jwks.json          | Get JWKS                                        |
| _AuthApi_      | [**logout**](docs/AuthApi.md#logout)                          | **POST** /auth/logout                   | Revoke refresh token family (logout everywhere) |
| _AuthApi_      | [**refresh**](docs/AuthApi.md#refresh)                        | **POST** /auth/refresh                  | Rotate refresh token                            |
| _LabelsApi_    | [**createLabel**](docs/LabelsApi.md#createlabel)              | **POST** /labels                        | Create label (org scoped)                       |
| _LabelsApi_    | [**deleteLabel**](docs/LabelsApi.md#deletelabel)              | **DELETE** /labels/{id}                 | Delete label (org scoped)                       |
| _LabelsApi_    | [**getLabel**](docs/LabelsApi.md#getlabel)                    | **GET** /labels/{id}                    | Get label by ID (org scoped)                    |
| _LabelsApi_    | [**listLabels**](docs/LabelsApi.md#listlabels)                | **GET** /labels                         | List node labels (org scoped)                   |
| _LabelsApi_    | [**updateLabel**](docs/LabelsApi.md#updatelabel)              | **PATCH** /labels/{id}                  | Update label (org scoped)                       |
| _MeApi_        | [**getMe**](docs/MeApi.md#getme)                              | **GET** /me                             | Get current user profile                        |
| _MeApi_        | [**updateMe**](docs/MeApi.md#updateme)                        | **PATCH** /me                           | Update current user profile                     |
| _MeAPIKeysApi_ | [**createUserAPIKey**](docs/MeAPIKeysApi.md#createuserapikey) | **POST** /me/api-keys                   | Create a new user API key                       |
| _MeAPIKeysApi_ | [**deleteUserAPIKey**](docs/MeAPIKeysApi.md#deleteuserapikey) | **DELETE** /me/api-keys/{id}            | Delete a user API key                           |
| _MeAPIKeysApi_ | [**listUserAPIKeys**](docs/MeAPIKeysApi.md#listuserapikeys)   | **GET** /me/api-keys                    | List my API keys                                |
| _OrgsApi_      | [**addOrUpdateMember**](docs/OrgsApi.md#addorupdatemember)    | **POST** /orgs/{id}/members             | Add or update a member (owner/admin)            |
| _OrgsApi_      | [**createOrg**](docs/OrgsApi.md#createorg)                    | **POST** /orgs                          | Create organization                             |
| _OrgsApi_      | [**createOrgKey**](docs/OrgsApi.md#createorgkey)              | **POST** /orgs/{id}/api-keys            | Create org key/secret pair (owner/admin)        |
| _OrgsApi_      | [**deleteOrg**](docs/OrgsApi.md#deleteorg)                    | **DELETE** /orgs/{id}                   | Delete organization (owner)                     |
| _OrgsApi_      | [**deleteOrgKey**](docs/OrgsApi.md#deleteorgkey)              | **DELETE** /orgs/{id}/api-keys/{key_id} | Delete org key (owner/admin)                    |
| _OrgsApi_      | [**getOrg**](docs/OrgsApi.md#getorg)                          | **GET** /orgs/{id}                      | Get organization                                |
| _OrgsApi_      | [**listMembers**](docs/OrgsApi.md#listmembers)                | **GET** /orgs/{id}/members              | List members in org                             |
| _OrgsApi_      | [**listMyOrgs**](docs/OrgsApi.md#listmyorgs)                  | **GET** /orgs                           | List organizations I belong to                  |
| _OrgsApi_      | [**listOrgKeys**](docs/OrgsApi.md#listorgkeys)                | **GET** /orgs/{id}/api-keys             | List org-scoped API keys (no secrets)           |
| _OrgsApi_      | [**removeMember**](docs/OrgsApi.md#removemember)              | **DELETE** /orgs/{id}/members/{user_id} | Remove a member (owner/admin)                   |
| _OrgsApi_      | [**updateOrg**](docs/OrgsApi.md#updateorg)                    | **PATCH** /orgs/{id}                    | Update organization (owner/admin)               |
| _ServersApi_   | [**createServer**](docs/ServersApi.md#createserver)           | **POST** /servers                       | Create server (org scoped)                      |
| _ServersApi_   | [**deleteServer**](docs/ServersApi.md#deleteserver)           | **DELETE** /servers/{id}                | Delete server (org scoped)                      |
| _ServersApi_   | [**getServer**](docs/ServersApi.md#getserver)                 | **GET** /servers/{id}                   | Get server by ID (org scoped)                   |
| _ServersApi_   | [**listServers**](docs/ServersApi.md#listservers)             | **GET** /servers                        | List servers (org scoped)                       |
| _ServersApi_   | [**updateServer**](docs/ServersApi.md#updateserver)           | **PATCH** /servers/{id}                 | Update server (org scoped)                      |
| _SshApi_       | [**createSSHKey**](docs/SshApi.md#createsshkey)               | **POST** /ssh                           | Create ssh keypair (org scoped)                 |
| _SshApi_       | [**deleteSSHKey**](docs/SshApi.md#deletesshkey)               | **DELETE** /ssh/{id}                    | Delete ssh keypair (org scoped)                 |
| _SshApi_       | [**downloadSSHKey**](docs/SshApi.md#downloadsshkey)           | **GET** /ssh/{id}/download              | Download ssh key files by ID (org scoped)       |
| _SshApi_       | [**getSSHKey**](docs/SshApi.md#getsshkey)                     | **GET** /ssh/{id}                       | Get ssh key by ID (org scoped)                  |
| _SshApi_       | [**listPublicSshKeys**](docs/SshApi.md#listpublicsshkeys)     | **GET** /ssh                            | List ssh keys (org scoped)                      |
| _TaintsApi_    | [**createTaint**](docs/TaintsApi.md#createtaint)              | **POST** /taints                        | Create node taint (org scoped)                  |
| _TaintsApi_    | [**deleteTaint**](docs/TaintsApi.md#deletetaint)              | **DELETE** /taints/{id}                 | Delete taint (org scoped)                       |
| _TaintsApi_    | [**getTaint**](docs/TaintsApi.md#gettaint)                    | **GET** /taints/{id}                    | Get node taint by ID (org scoped)               |
| _TaintsApi_    | [**listTaints**](docs/TaintsApi.md#listtaints)                | **GET** /taints                         | List node pool taints (org scoped)              |
| _TaintsApi_    | [**updateTaint**](docs/TaintsApi.md#updatetaint)              | **PATCH** /taints/{id}                  | Update node taint (org scoped)                  |

### Models

- [DtoAuthStartResponse](docs/DtoAuthStartResponse.md)
- [DtoCreateLabelRequest](docs/DtoCreateLabelRequest.md)
- [DtoCreateSSHRequest](docs/DtoCreateSSHRequest.md)
- [DtoCreateServerRequest](docs/DtoCreateServerRequest.md)
- [DtoCreateTaintRequest](docs/DtoCreateTaintRequest.md)
- [DtoJWK](docs/DtoJWK.md)
- [DtoJWKS](docs/DtoJWKS.md)
- [DtoLabelResponse](docs/DtoLabelResponse.md)
- [DtoLogoutRequest](docs/DtoLogoutRequest.md)
- [DtoRefreshRequest](docs/DtoRefreshRequest.md)
- [DtoServerResponse](docs/DtoServerResponse.md)
- [DtoSshResponse](docs/DtoSshResponse.md)
- [DtoSshRevealResponse](docs/DtoSshRevealResponse.md)
- [DtoTaintResponse](docs/DtoTaintResponse.md)
- [DtoTokenPair](docs/DtoTokenPair.md)
- [DtoUpdateLabelRequest](docs/DtoUpdateLabelRequest.md)
- [DtoUpdateServerRequest](docs/DtoUpdateServerRequest.md)
- [DtoUpdateTaintRequest](docs/DtoUpdateTaintRequest.md)
- [HandlersCreateUserKeyRequest](docs/HandlersCreateUserKeyRequest.md)
- [HandlersMeResponse](docs/HandlersMeResponse.md)
- [HandlersMemberOut](docs/HandlersMemberOut.md)
- [HandlersMemberUpsertReq](docs/HandlersMemberUpsertReq.md)
- [HandlersOrgCreateReq](docs/HandlersOrgCreateReq.md)
- [HandlersOrgKeyCreateReq](docs/HandlersOrgKeyCreateReq.md)
- [HandlersOrgKeyCreateResp](docs/HandlersOrgKeyCreateResp.md)
- [HandlersOrgUpdateReq](docs/HandlersOrgUpdateReq.md)
- [HandlersUpdateMeRequest](docs/HandlersUpdateMeRequest.md)
- [HandlersUserAPIKeyOut](docs/HandlersUserAPIKeyOut.md)
- [ModelsAPIKey](docs/ModelsAPIKey.md)
- [ModelsOrganization](docs/ModelsOrganization.md)
- [ModelsUser](docs/ModelsUser.md)
- [ModelsUserEmail](docs/ModelsUserEmail.md)
- [UtilsErrorResponse](docs/UtilsErrorResponse.md)

### Authorization

Authentication schemes defined for the API:
<a id="ApiKeyAuth"></a>

#### ApiKeyAuth

- **Type**: API key
- **API key parameter name**: `X-API-KEY`
- **Location**: HTTP header
  <a id="BearerAuth"></a>

#### BearerAuth

- **Type**: API key
- **API key parameter name**: `Authorization`
- **Location**: HTTP header
  <a id="OrgKeyAuth"></a>

#### OrgKeyAuth

- **Type**: API key
- **API key parameter name**: `X-ORG-KEY`
- **Location**: HTTP header
  <a id="OrgSecretAuth"></a>

#### OrgSecretAuth

- **Type**: API key
- **API key parameter name**: `X-ORG-SECRET`
- **Location**: HTTP header

## About

This TypeScript SDK client supports the [Fetch API](https://fetch.spec.whatwg.org/)
and is automatically generated by the
[OpenAPI Generator](https://openapi-generator.tech) project:

- API version: `1.0`
- Package version: `0.1.0`
- Generator version: `7.17.0`
- Build package: `org.openapitools.codegen.languages.TypeScriptFetchClientCodegen`

The generated npm module supports the following:

- Environments
  - Node.js
  - Webpack
  - Browserify
- Language levels
  - ES5 - you must have a Promises/A+ library installed
  - ES6
- Module systems
  - CommonJS
  - ES6 module system

## Development

### Building

To build the TypeScript source code, you need to have Node.js and npm installed.
After cloning the repository, navigate to the project directory and run:

```bash
npm install
npm run build
```

### Publishing

Once you've built the package, you can publish it to npm:

```bash
npm publish
```

## License

[]()
