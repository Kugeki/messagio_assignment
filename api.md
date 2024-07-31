


# Messagio Assigment
Test task to Messagio.
  

## Informations

### Version

0.1

### Contact

  

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  messages

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /messages/stats | [get messages stats](#get-messages-stats) | Get messages stats |
| POST | /messages | [post messages](#post-messages) | Create a message |
  


## Paths

### <span id="get-messages-stats"></span> Get messages stats (*GetMessagesStats*)

```
GET /messages/stats
```

get messages stats

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-messages-stats-200) | OK | OK | ✓ | [schema](#get-messages-stats-200-schema) |
| [429](#get-messages-stats-429) | Too Many Requests | Too Many Requests | ✓ | [schema](#get-messages-stats-429-schema) |
| [500](#get-messages-stats-500) | Internal Server Error | Internal Server Error | ✓ | [schema](#get-messages-stats-500-schema) |

#### Responses


##### <span id="get-messages-stats-200"></span> 200 - OK
Status: OK

###### <span id="get-messages-stats-200-schema"></span> Schema
   
  

[DtoGetStatsResp](#dto-get-stats-resp)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="get-messages-stats-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-messages-stats-429-schema"></span> Schema
   
  

[DtoHTTPError](#dto-http-error)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="get-messages-stats-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-messages-stats-500-schema"></span> Schema

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

### <span id="post-messages"></span> Create a message (*PostMessages*)

```
POST /messages
```

create a message

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| message | `body` | [DtoCreateMessageReq](#dto-create-message-req) | `models.DtoCreateMessageReq` | | ✓ | | Create message |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-messages-201) | Created | Created | ✓ | [schema](#post-messages-201-schema) |
| [400](#post-messages-400) | Bad Request | Bad Request | ✓ | [schema](#post-messages-400-schema) |
| [409](#post-messages-409) | Conflict | Conflict | ✓ | [schema](#post-messages-409-schema) |
| [422](#post-messages-422) | Unprocessable Entity | Unprocessable Entity | ✓ | [schema](#post-messages-422-schema) |
| [429](#post-messages-429) | Too Many Requests | Too Many Requests | ✓ | [schema](#post-messages-429-schema) |
| [500](#post-messages-500) | Internal Server Error | Internal Server Error | ✓ | [schema](#post-messages-500-schema) |

#### Responses


##### <span id="post-messages-201"></span> 201 - Created
Status: Created

###### <span id="post-messages-201-schema"></span> Schema
   
  

[DtoCreateMessageResp](#dto-create-message-resp)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="post-messages-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-messages-400-schema"></span> Schema
   
  

[DtoHTTPError](#dto-http-error)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="post-messages-409"></span> 409 - Conflict
Status: Conflict

###### <span id="post-messages-409-schema"></span> Schema
   
  

[DtoHTTPError](#dto-http-error)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="post-messages-422"></span> 422 - Unprocessable Entity
Status: Unprocessable Entity

###### <span id="post-messages-422-schema"></span> Schema
   
  

[DtoHTTPError](#dto-http-error)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="post-messages-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="post-messages-429-schema"></span> Schema
   
  

[DtoHTTPError](#dto-http-error)

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

##### <span id="post-messages-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-messages-500-schema"></span> Schema

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| X-RateLimit-Limit | string | `string` |  |  | Request limit per minute |
| X-RateLimit-Remaining | string | `string` |  |  | The number of requests left for the time window |
| X-RateLimit-Reset | string | `string` |  |  | The remaining window before the rate limit resets in UTC epoch seconds |

## Models

### <span id="dto-create-message-req"></span> dto.CreateMessageReq


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | |  |  |
| processed | boolean| `bool` |  | |  |  |



### <span id="dto-create-message-resp"></span> dto.CreateMessageResp


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | |  |  |
| id | integer| `int64` |  | |  |  |
| processed | boolean| `bool` |  | |  |  |



### <span id="dto-get-stats-resp"></span> dto.GetStatsResp


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| all | integer| `int64` |  | |  |  |
| processed | integer| `int64` |  | |  |  |



### <span id="dto-http-error"></span> dto.HTTPError


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| error | string| `string` |  | |  |  |


###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###
###