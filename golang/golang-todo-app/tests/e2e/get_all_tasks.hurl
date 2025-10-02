GET {{base_url}}/api/v1/tasks
HTTP 200

[Asserts]
status == 200
jsonpath "$.message" == "OK"
jsonpath "$.data" count != 0

