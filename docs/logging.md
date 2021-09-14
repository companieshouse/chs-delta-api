# Logging in the chs-delta-api

This application uses the following standard format to log messages using the chs.go logger.

## Successful journey

When a request is successful you will get the following stack of logs:
```go
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text","message_key":"message_key_value"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text","offset":offset_int,"partition":partition_int,"topic":"officers-delta"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"duration":time,"end":"date_time_stamp","method":"http_method","path":"url_path","start":"date_time_stamp","status":http_status},"event":"request","namespace":"chs-delta-api"}
```

With an example of this being:
```go
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.820620448+01:00","data":{"message":"Validating request using: ","spec_location":"apispec/api-spec.yml"},"event":"info","namespace":"chs-delta-api"}
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.823892338+01:00","data":{"message":"Request validated. No errors were found."},"event":"info","namespace":"chs-delta-api"}
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.825429905+01:00","data":{"message":"Sending message","offset":1,"partition":6,"topic":"officers-delta"},"event":"info","namespace":"chs-delta-api"}
{"context":"zWJcenQdHgzRmoK1hL_ridjiybm1","created":"2021-09-13T12:19:29.82546636+01:00","data":{"duration":time,"end":"2021-09-13T12:19:29.825462326+01:00","method":"POST","path":"/delta/officers","start":"2021-09-13T12:19:29.820460095+01:00","status":200},"event":"request","namespace":"chs-delta-api"}
```

## Validation error journey

When a request returns validation errors you will get the following stack of logs:
```go
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text","message_key":"message_key_value"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"message":"message_text"},"event":"info","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"error":{},"message":"ch_errors_object_array"},"event":"error","namespace":"chs-delta-api"}
{"context":"context_id","created":"date_time_stamp","data":{"duration":time,"end":"date_time_stamp","method":"http_method","path":"url_path","start":"date_time_stamp","status":http_status},"event":"request","namespace":"chs-delta-api"}
```

With an example of this being:
```go
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.042735672+01:00","data":{"message":"Validating request using: ","spec_location":"apispec/api-spec.yml"},"event":"info","namespace":"chs-delta-api"}
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.042891994+01:00","data":{"message":"Request validated. Errors found."},"event":"info","namespace":"chs-delta-api"}
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.045197901+01:00","data":{"error":{},"message":"{Error:value is not one of the allowed values, ErrorValues:map[secure_director:test], Location:officers.0.secure_director, LocationType:json-path, Type:ch:validation},"},"event":"error","namespace":"chs-delta-api"}
{"context":"Dy7rFtAq3G5G9m60MZ1jlgkgeyLD","created":"2021-09-14T11:03:02.045648992+01:00","data":{"duration":3111079,"end":"2021-09-14T11:03:02.045644507+01:00","method":"POST","path":"/delta/officers","start":"2021-09-14T11:03:02.042533743+01:00","status":400},"event":"request","namespace":"chs-delta-api"}
```