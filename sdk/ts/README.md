# @glueops/autoglue-sdk-go@0.1.0

A TypeScript SDK client for the localhost API.

## Usage

First, install the SDK from npm.

```bash
npm install @glueops/autoglue-sdk-go --save
```

Next, try it out.

```ts
import {
  Configuration,
  AnnotationsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateAnnotationRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AnnotationsApi(config);

  const body = {
    // DtoCreateAnnotationRequest | Annotation payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateAnnotationRequest;

  try {
    const data = await api.createAnnotation(body);
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

All URIs are relative to _/api/v1_

| Class            | Method                                                                          | HTTP request                                           | Description                                           |
| ---------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------ | ----------------------------------------------------- |
| _AnnotationsApi_ | [**createAnnotation**](docs/AnnotationsApi.md#createannotation)                 | **POST** /annotations                                  | Create annotation (org scoped)                        |
| _AnnotationsApi_ | [**deleteAnnotation**](docs/AnnotationsApi.md#deleteannotation)                 | **DELETE** /annotations/{id}                           | Delete annotation (org scoped)                        |
| _AnnotationsApi_ | [**getAnnotation**](docs/AnnotationsApi.md#getannotation)                       | **GET** /annotations/{id}                              | Get annotation by ID (org scoped)                     |
| _AnnotationsApi_ | [**listAnnotations**](docs/AnnotationsApi.md#listannotations)                   | **GET** /annotations                                   | List annotations (org scoped)                         |
| _AnnotationsApi_ | [**updateAnnotation**](docs/AnnotationsApi.md#updateannotation)                 | **PATCH** /annotations/{id}                            | Update annotation (org scoped)                        |
| _ArcherAdminApi_ | [**adminCancelArcherJob**](docs/ArcherAdminApi.md#admincancelarcherjob)         | **POST** /admin/archer/jobs/{id}/cancel                | Cancel an Archer job (admin)                          |
| _ArcherAdminApi_ | [**adminEnqueueArcherJob**](docs/ArcherAdminApi.md#adminenqueuearcherjob)       | **POST** /admin/archer/jobs                            | Enqueue a new Archer job (admin)                      |
| _ArcherAdminApi_ | [**adminListArcherJobs**](docs/ArcherAdminApi.md#adminlistarcherjobs)           | **GET** /admin/archer/jobs                             | List Archer jobs (admin)                              |
| _ArcherAdminApi_ | [**adminListArcherQueues**](docs/ArcherAdminApi.md#adminlistarcherqueues)       | **GET** /admin/archer/queues                           | List Archer queues (admin)                            |
| _ArcherAdminApi_ | [**adminRetryArcherJob**](docs/ArcherAdminApi.md#adminretryarcherjob)           | **POST** /admin/archer/jobs/{id}/retry                 | Retry a failed/canceled Archer job (admin)            |
| _AuthApi_        | [**authCallback**](docs/AuthApi.md#authcallback)                                | **GET** /auth/{provider}/callback                      | Handle social login callback                          |
| _AuthApi_        | [**authStart**](docs/AuthApi.md#authstart)                                      | **POST** /auth/{provider}/start                        | Begin social login                                    |
| _AuthApi_        | [**getJWKS**](docs/AuthApi.md#getjwks)                                          | **GET** /.well-known/jwks.json                         | Get JWKS                                              |
| _AuthApi_        | [**logout**](docs/AuthApi.md#logout)                                            | **POST** /auth/logout                                  | Revoke refresh token family (logout everywhere)       |
| _AuthApi_        | [**refresh**](docs/AuthApi.md#refresh)                                          | **POST** /auth/refresh                                 | Rotate refresh token                                  |
| _HealthApi_      | [**healthCheckOperationId**](docs/HealthApi.md#healthcheckoperationid)          | **GET** /healthz                                       | Basic health check                                    |
| _LabelsApi_      | [**createLabel**](docs/LabelsApi.md#createlabel)                                | **POST** /labels                                       | Create label (org scoped)                             |
| _LabelsApi_      | [**deleteLabel**](docs/LabelsApi.md#deletelabel)                                | **DELETE** /labels/{id}                                | Delete label (org scoped)                             |
| _LabelsApi_      | [**getLabel**](docs/LabelsApi.md#getlabel)                                      | **GET** /labels/{id}                                   | Get label by ID (org scoped)                          |
| _LabelsApi_      | [**listLabels**](docs/LabelsApi.md#listlabels)                                  | **GET** /labels                                        | List node labels (org scoped)                         |
| _LabelsApi_      | [**updateLabel**](docs/LabelsApi.md#updatelabel)                                | **PATCH** /labels/{id}                                 | Update label (org scoped)                             |
| _MeApi_          | [**getMe**](docs/MeApi.md#getme)                                                | **GET** /me                                            | Get current user profile                              |
| _MeApi_          | [**updateMe**](docs/MeApi.md#updateme)                                          | **PATCH** /me                                          | Update current user profile                           |
| _MeAPIKeysApi_   | [**createUserAPIKey**](docs/MeAPIKeysApi.md#createuserapikey)                   | **POST** /me/api-keys                                  | Create a new user API key                             |
| _MeAPIKeysApi_   | [**deleteUserAPIKey**](docs/MeAPIKeysApi.md#deleteuserapikey)                   | **DELETE** /me/api-keys/{id}                           | Delete a user API key                                 |
| _MeAPIKeysApi_   | [**listUserAPIKeys**](docs/MeAPIKeysApi.md#listuserapikeys)                     | **GET** /me/api-keys                                   | List my API keys                                      |
| _NodePoolsApi_   | [**attachNodePoolAnnotations**](docs/NodePoolsApi.md#attachnodepoolannotations) | **POST** /node-pools/{id}/annotations                  | Attach annotation to a node pool (org scoped)         |
| _NodePoolsApi_   | [**attachNodePoolLabels**](docs/NodePoolsApi.md#attachnodepoollabels)           | **POST** /node-pools/{id}/labels                       | Attach labels to a node pool (org scoped)             |
| _NodePoolsApi_   | [**attachNodePoolServers**](docs/NodePoolsApi.md#attachnodepoolservers)         | **POST** /node-pools/{id}/servers                      | Attach servers to a node pool (org scoped)            |
| _NodePoolsApi_   | [**attachNodePoolTaints**](docs/NodePoolsApi.md#attachnodepooltaints)           | **POST** /node-pools/{id}/taints                       | Attach taints to a node pool (org scoped)             |
| _NodePoolsApi_   | [**createNodePool**](docs/NodePoolsApi.md#createnodepool)                       | **POST** /node-pools                                   | Create node pool (org scoped)                         |
| _NodePoolsApi_   | [**deleteNodePool**](docs/NodePoolsApi.md#deletenodepool)                       | **DELETE** /node-pools/{id}                            | Delete node pool (org scoped)                         |
| _NodePoolsApi_   | [**detachNodePoolAnnotation**](docs/NodePoolsApi.md#detachnodepoolannotation)   | **DELETE** /node-pools/{id}/annotations/{annotationId} | Detach one annotation from a node pool (org scoped)   |
| _NodePoolsApi_   | [**detachNodePoolLabel**](docs/NodePoolsApi.md#detachnodepoollabel)             | **DELETE** /node-pools/{id}/labels/{labelId}           | Detach one label from a node pool (org scoped)        |
| _NodePoolsApi_   | [**detachNodePoolServer**](docs/NodePoolsApi.md#detachnodepoolserver)           | **DELETE** /node-pools/{id}/servers/{serverId}         | Detach one server from a node pool (org scoped)       |
| _NodePoolsApi_   | [**detachNodePoolTaint**](docs/NodePoolsApi.md#detachnodepooltaint)             | **DELETE** /node-pools/{id}/taints/{taintId}           | Detach one taint from a node pool (org scoped)        |
| _NodePoolsApi_   | [**getNodePool**](docs/NodePoolsApi.md#getnodepool)                             | **GET** /node-pools/{id}                               | Get node pool by ID (org scoped)                      |
| _NodePoolsApi_   | [**listNodePoolAnnotations**](docs/NodePoolsApi.md#listnodepoolannotations)     | **GET** /node-pools/{id}/annotations                   | List annotations attached to a node pool (org scoped) |
| _NodePoolsApi_   | [**listNodePoolLabels**](docs/NodePoolsApi.md#listnodepoollabels)               | **GET** /node-pools/{id}/labels                        | List labels attached to a node pool (org scoped)      |
| _NodePoolsApi_   | [**listNodePoolServers**](docs/NodePoolsApi.md#listnodepoolservers)             | **GET** /node-pools/{id}/servers                       | List servers attached to a node pool (org scoped)     |
| _NodePoolsApi_   | [**listNodePoolTaints**](docs/NodePoolsApi.md#listnodepooltaints)               | **GET** /node-pools/{id}/taints                        | List taints attached to a node pool (org scoped)      |
| _NodePoolsApi_   | [**listNodePools**](docs/NodePoolsApi.md#listnodepools)                         | **GET** /node-pools                                    | List node pools (org scoped)                          |
| _NodePoolsApi_   | [**updateNodePool**](docs/NodePoolsApi.md#updatenodepool)                       | **PATCH** /node-pools/{id}                             | Update node pool (org scoped)                         |
| _OrgsApi_        | [**addOrUpdateMember**](docs/OrgsApi.md#addorupdatemember)                      | **POST** /orgs/{id}/members                            | Add or update a member (owner/admin)                  |
| _OrgsApi_        | [**createOrg**](docs/OrgsApi.md#createorg)                                      | **POST** /orgs                                         | Create organization                                   |
| _OrgsApi_        | [**createOrgKey**](docs/OrgsApi.md#createorgkey)                                | **POST** /orgs/{id}/api-keys                           | Create org key/secret pair (owner/admin)              |
| _OrgsApi_        | [**deleteOrg**](docs/OrgsApi.md#deleteorg)                                      | **DELETE** /orgs/{id}                                  | Delete organization (owner)                           |
| _OrgsApi_        | [**deleteOrgKey**](docs/OrgsApi.md#deleteorgkey)                                | **DELETE** /orgs/{id}/api-keys/{key_id}                | Delete org key (owner/admin)                          |
| _OrgsApi_        | [**getOrg**](docs/OrgsApi.md#getorg)                                            | **GET** /orgs/{id}                                     | Get organization                                      |
| _OrgsApi_        | [**listMembers**](docs/OrgsApi.md#listmembers)                                  | **GET** /orgs/{id}/members                             | List members in org                                   |
| _OrgsApi_        | [**listMyOrgs**](docs/OrgsApi.md#listmyorgs)                                    | **GET** /orgs                                          | List organizations I belong to                        |
| _OrgsApi_        | [**listOrgKeys**](docs/OrgsApi.md#listorgkeys)                                  | **GET** /orgs/{id}/api-keys                            | List org-scoped API keys (no secrets)                 |
| _OrgsApi_        | [**removeMember**](docs/OrgsApi.md#removemember)                                | **DELETE** /orgs/{id}/members/{user_id}                | Remove a member (owner/admin)                         |
| _OrgsApi_        | [**updateOrg**](docs/OrgsApi.md#updateorg)                                      | **PATCH** /orgs/{id}                                   | Update organization (owner/admin)                     |
| _ServersApi_     | [**createServer**](docs/ServersApi.md#createserver)                             | **POST** /servers                                      | Create server (org scoped)                            |
| _ServersApi_     | [**deleteServer**](docs/ServersApi.md#deleteserver)                             | **DELETE** /servers/{id}                               | Delete server (org scoped)                            |
| _ServersApi_     | [**getServer**](docs/ServersApi.md#getserver)                                   | **GET** /servers/{id}                                  | Get server by ID (org scoped)                         |
| _ServersApi_     | [**listServers**](docs/ServersApi.md#listservers)                               | **GET** /servers                                       | List servers (org scoped)                             |
| _ServersApi_     | [**updateServer**](docs/ServersApi.md#updateserver)                             | **PATCH** /servers/{id}                                | Update server (org scoped)                            |
| _SshApi_         | [**createSSHKey**](docs/SshApi.md#createsshkey)                                 | **POST** /ssh                                          | Create ssh keypair (org scoped)                       |
| _SshApi_         | [**deleteSSHKey**](docs/SshApi.md#deletesshkey)                                 | **DELETE** /ssh/{id}                                   | Delete ssh keypair (org scoped)                       |
| _SshApi_         | [**downloadSSHKey**](docs/SshApi.md#downloadsshkey)                             | **GET** /ssh/{id}/download                             | Download ssh key files by ID (org scoped)             |
| _SshApi_         | [**getSSHKey**](docs/SshApi.md#getsshkey)                                       | **GET** /ssh/{id}                                      | Get ssh key by ID (org scoped)                        |
| _SshApi_         | [**listPublicSshKeys**](docs/SshApi.md#listpublicsshkeys)                       | **GET** /ssh                                           | List ssh keys (org scoped)                            |
| _TaintsApi_      | [**createTaint**](docs/TaintsApi.md#createtaint)                                | **POST** /taints                                       | Create node taint (org scoped)                        |
| _TaintsApi_      | [**deleteTaint**](docs/TaintsApi.md#deletetaint)                                | **DELETE** /taints/{id}                                | Delete taint (org scoped)                             |
| _TaintsApi_      | [**getTaint**](docs/TaintsApi.md#gettaint)                                      | **GET** /taints/{id}                                   | Get node taint by ID (org scoped)                     |
| _TaintsApi_      | [**listTaints**](docs/TaintsApi.md#listtaints)                                  | **GET** /taints                                        | List node pool taints (org scoped)                    |
| _TaintsApi_      | [**updateTaint**](docs/TaintsApi.md#updatetaint)                                | **PATCH** /taints/{id}                                 | Update node taint (org scoped)                        |

### Models

- [DtoAnnotationResponse](docs/DtoAnnotationResponse.md)
- [DtoAttachAnnotationsRequest](docs/DtoAttachAnnotationsRequest.md)
- [DtoAttachLabelsRequest](docs/DtoAttachLabelsRequest.md)
- [DtoAttachServersRequest](docs/DtoAttachServersRequest.md)
- [DtoAttachTaintsRequest](docs/DtoAttachTaintsRequest.md)
- [DtoAuthStartResponse](docs/DtoAuthStartResponse.md)
- [DtoCreateAnnotationRequest](docs/DtoCreateAnnotationRequest.md)
- [DtoCreateLabelRequest](docs/DtoCreateLabelRequest.md)
- [DtoCreateNodePoolRequest](docs/DtoCreateNodePoolRequest.md)
- [DtoCreateSSHRequest](docs/DtoCreateSSHRequest.md)
- [DtoCreateServerRequest](docs/DtoCreateServerRequest.md)
- [DtoCreateTaintRequest](docs/DtoCreateTaintRequest.md)
- [DtoJWK](docs/DtoJWK.md)
- [DtoJWKS](docs/DtoJWKS.md)
- [DtoJob](docs/DtoJob.md)
- [DtoJobStatus](docs/DtoJobStatus.md)
- [DtoLabelResponse](docs/DtoLabelResponse.md)
- [DtoLogoutRequest](docs/DtoLogoutRequest.md)
- [DtoNodePoolResponse](docs/DtoNodePoolResponse.md)
- [DtoPageJob](docs/DtoPageJob.md)
- [DtoQueueInfo](docs/DtoQueueInfo.md)
- [DtoRefreshRequest](docs/DtoRefreshRequest.md)
- [DtoServerResponse](docs/DtoServerResponse.md)
- [DtoSshResponse](docs/DtoSshResponse.md)
- [DtoSshRevealResponse](docs/DtoSshRevealResponse.md)
- [DtoTaintResponse](docs/DtoTaintResponse.md)
- [DtoTokenPair](docs/DtoTokenPair.md)
- [DtoUpdateAnnotationRequest](docs/DtoUpdateAnnotationRequest.md)
- [DtoUpdateLabelRequest](docs/DtoUpdateLabelRequest.md)
- [DtoUpdateNodePoolRequest](docs/DtoUpdateNodePoolRequest.md)
- [DtoUpdateServerRequest](docs/DtoUpdateServerRequest.md)
- [DtoUpdateTaintRequest](docs/DtoUpdateTaintRequest.md)
- [HandlersCreateUserKeyRequest](docs/HandlersCreateUserKeyRequest.md)
- [HandlersHealthStatus](docs/HandlersHealthStatus.md)
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
