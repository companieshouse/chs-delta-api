# Logging in the chs-delta-api
This application uses the following standard format to log messages using the chs.go logger.

## Notable logs
The following sub sections will cover all notable logs that are output by this service.


#### Validation fails
if a request fails validation you will get the following log:

```go
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"error":{},"message":"ch_errors_object_array"},"event":"error","namespace":"chs-delta-api"}
```

The following message will be logged:
```go
"message":"Request validated. Errors found."
```

followed by a message containing the errors found.
```go
"message":"{Error:value is not one of the allowed values, ErrorValues:map[secure_director:test], Location:officers.0.secure_director, LocationType:json-path, Type:ch:validation},"
```


The following variables are of interest to Kibana:

`"context":"context_id"` - This is the only variable which can be used to track a request through the full end to end journey.

With an example of this being:
```go
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.042891994+01:00","data":{"message":"Request validated. Errors found."},"event":"info","namespace":"chs-delta-api"}
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.045197901+01:00","data":{"error":{},"message":"{Error:value is not one of the allowed values, ErrorValues:map[secure_director:test], Location:officers.0.secure_director, LocationType:json-path, Type:ch:validation},"},"event":"error","namespace":"chs-delta-api"}
```
---

#### Message sent to Kafka
```go
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text","offset":offset_int,"partition":partition_int,"topic":"topic_choice"},"event":"info","namespace":"chs-delta-api"}
```
The following message will be logged:
```go
"message":"Sent message"
```

The following variables are of interest to Kibana:

`"context":"context_id"` - This is the only variable which can be used to track a request through the full end to end journey.

`"offset":offset_int, "partition":partition_int` - These values are only returned if a message has successfully been sent to Kafka.

With an example of this being:
```go
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.825429905+01:00","data":{"message":"Sent message","offset":1,"partition":6,"topic":"officers-delta"},"event":"info","namespace":"chs-delta-api"}
```
---

#### Message fails to be sent to Kafka
If a message fails to be sent to Kafka you will get the following logs:

```go
{"context":"context_id","created":"date_time_stamp","data":{"error":kafka_error_int,"message":"kafka_message","topic":"topic_choice"},"event":"error","namespace":"chs-delta-api"}
```

The message that will be logged changes depending on the error Kafka returns. This can vary quite a bit, but here are 2 examples:

```go
"message":"kafka server: Replication-factor is invalid."
```

```go
"message":"kafka server: In the middle of a leadership election, there is currently no leader for this partition and hence it is unavailable for writes."
```

The following variables are of interest to Kibana:

`"context":"context_id"` - This is the only variable which can be used to track a request through the full end to end journey.

`"error":kafka_error_int` - This is the error code returned by Kafka which will allow you to trace the cause of the error.

`"topic":"topic_choice"` - This is the topic which was attempted to send the message too.

With an example of this being:
```go
{"context":"aaron-temp-id","created":"2021-09-14T13:01:48.67308Z","data":{"error":5,"message":"kafka server: In the middle of a leadership election, there is currently no leader for this partition and hence it is unavailable for writes.","topic":"officers-delta"},
"event":"error","namespace":"chs-delta-api"}
```

---

#### Status codes logged

`"status":http_status` - Below is a full list of HTTP status codes returned by this service, to indicate whether a request was successful or not:
- `200` (Successful request)
- `400` (Validation error)
- `500` (Internal Server Error)
- `401` (Unauthorized)

This does not cover other status codes such as timeouts and responses from other services such as ERIC.

Example of a Successful request:
```go
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.82546636+01:00","data":{"duration":time,"end":"2021-09-13T12:19:29.825462326+01:00","method":"POST","path":"/delta/officers","start":"2021-09-13T12:19:29.820460095+01:00","status":200},"event":"request","namespace":"chs-delta-api"}
```
 
Example of a Validation Error request:
```go
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.045648992+01:00","data":{"duration":3111079,"end":"2021-09-14T11:03:02.045644507+01:00","method":"POST","path":"/delta/officers","start":"2021-09-14T11:03:02.042533743+01:00","status":400},"event":"request","namespace":"chs-delta-api"}
```

---

#### ContextId and its use in tracking requests from CHIPS to CHS
The context_id variable as mentioned in other sections of this document is the variable which should be used to track a filing from CHIPS to CHS.
It is retrieved / created in the chs-delta-api at the start of handling a request. First the chs-delta-api will attempt to pull the contxt_id from 
the request headers; one of the provided headers from ERIC is the x-request-id, this is the context_id. If the x-request-id header is not present 
then the chs-delta-api will create a context_id.

It is passed with the request received from CHIPS onto a Kafka topic as part of the avro schema (inside the avro schema you will find the contextId field).
A go struct implementation of the avro schema can be found in the `models/chsDelta.go` file.