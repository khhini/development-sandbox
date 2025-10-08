GET {{base_url}}/api/v1/healthz
HTTP 200

[Asserts]
status == 200
jsonpath "$.status" == "healthy"
