# Kafka

Брокер находится на порту 9092.


Kafka-UI на порту 8090.

## Топики
- messages-to-process
- processed-messages

### messages-to-process
Формат данных:

```json
{
  "id": 0,
  "content": "string",
  "processed": false
}
```


### processed-messages
Формат данных:

```json
{
  "id": 0
}
```