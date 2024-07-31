<!--
generated at https://swagger-markdown-ui.netlify.app/
-->

# Messagio Assigment
Test task to Messagio.

## Version: 0.1

### /messages

#### POST
##### Summary:

Create a message

##### Description:

create a message

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| message | body | Create message | Yes | [dto.CreateMessageReq](#dto.CreateMessageReq) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [dto.CreateMessageResp](#dto.CreateMessageResp) |
| 400 | Bad Request | [dto.HTTPError](#dto.HTTPError) |
| 409 | Conflict | [dto.HTTPError](#dto.HTTPError) |
| 422 | Unprocessable Entity | [dto.HTTPError](#dto.HTTPError) |
| 429 | Too Many Requests | [dto.HTTPError](#dto.HTTPError) |
| 500 | Internal Server Error |  |

### /messages/stats

#### GET
##### Summary:

Get messages stats

##### Description:

get messages stats

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [dto.GetStatsResp](#dto.GetStatsResp) |
| 429 | Too Many Requests | [dto.HTTPError](#dto.HTTPError) |
| 500 | Internal Server Error |  |

### Models


#### dto.CreateMessageReq

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| processed | boolean |  | No |

#### dto.CreateMessageResp

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| id | integer |  | No |
| processed | boolean |  | No |

#### dto.GetStatsResp

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| all | integer |  | No |
| processed | integer |  | No |

#### dto.HTTPError

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |