# Developer Guide: for Logging with Google Cloud Logging

This document outlines key best practices for developers to effectively utilize Google Cloud Logging, ensuring logs are actionable, easy to query, and integrated seamlessly with the Google Cloud ecosystem.

## Adopt Structured Logging

Instead of logging plain text, use a JSON format to output log entries. This enables Google Cloud Logging to parse and index specific fields, making logs significantly easier to search, filter, and analyze.

| Key Practices | Details| GCL Integration|
| --------------- | --------------- | --------------- |
| output JSON | Configure your logging library to emit logs as single-line JSON object to `stdout` or `stderr` | GCL automatically parse valid JSON output to populate the `jsonPayload` field. |
| Use GCL specific fields | When logging to `stdout` / `stderr` from services like Cloud Run, GKE, or Cloud Functions, include special JSON keys that GCL will automatically map to top level `LogEntry` fields. | `severity`: Maps to log level (e.g., `INFO`, `ERROR`). `message`: Maps to the primary text message. `logging.googleapis.com/trace`: Essential for request correlation. |
| Log via Client Libraries | For more granular control over fields like `logName` or `resource` labels, use the official Google Cloud Logging client library for your language. | Allow direct control over the entire `LogEntry` structure. |

Structured Log Example

```json
{
  "insertId": "42",
  "jsonPayload": {
    "message": "There was an error in the application",
    "times": "2020-10-12T07:20:50.52Z"
  },
  "httpRequest": {
    "requestMethod": "GET"
  },
  "resource": {
    "type": "k8s_container",
    "labels": {
      "container_name": "hello-app",
      "pod_name": "helloworld-gke-6cfd6f4599-9wff8",
      "project_id": "stackdriver-sandbox-92334288",
      "namespace_name": "default",
      "location": "us-west4",
      "cluster_name": "helloworld-gke"
    }
  },
  "timestamp": "2020-11-07T15:57:35.945508391Z",
  "severity": "ERROR",
  "labels": {
    "user_label_2": "value_2",
    "user_label_1": "value_1"
  },
  "logName": "projects/stackdriver-sandbox-92334288/logs/stdout",
  "operation": {
    "id": "get_data",
    "producer": "github.com/MyProject/MyApplication",
    "first": true
  },
  "trace": "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824![](chrome-extension://efbjojhplkelaegfbieplglfidafgoka/icons/vt-logo.svg)",
  "sourceLocation": {
    "file": "get_data.py",
    "line": "142",
    "function": "getData"
  },
  "receiveTimestamp": "2020-11-07T15:57:42.411414059Z",
  "spanId": "000000000000004a"
}
```

## Context and correlation

Effective logging requires the ability to trace an event through your entire system, especially in a micro services architecture.

| Key Practices | Details | Impact in GCL |
| --------------- | --------------- | --------------- |
| Requiest ID / correlation ID | Include a unique identifier for every inbound request across all log entries related to that request. Pass this ID between services. | Essential for queriying all logs related to a single user interation or request flow. |
| Trace Integration | Leverage Cloud Trace integration by adding the trace context to your log entry. | GCL automatically links log entries to the corresponding trace span, allowing developers to see the logs alongside the request latency timeline. Use the key `logging.googleapis.com/trace` |
| Resouce & Service Context | Always include service, name, module, version, and environment (e.g., `dev`, `prod`) | GCL often auto-populates the `resouce` field for managed services (e.g., Cloud Run, GKE). Explicitly logging the service name helps in filtering and log-based metrics. |
| Caconical Log | Caconical log line is a single log, comprehensive log entry that is `created at the end of each requests` | redaction of noise and ensure each log correlate with status of each request for example `success` or `failed` |

## Log Levels and Severity

Use log levels consistently to manage noise, control costs, and quickly identify critical issues. GCL maps standard log levels to the `severity` filed.

| Level | Purpose | Recommended Usage |
| --------------- | --------------- | --------------- |
| `DEBUG` | Extreamly detailed inforamtion, often too voluminous for production | Development and deep troubleshooting in non-production. |
| `INFO` | General application flow, major application state change, successfull operations. | Low-volume, high-value status updates in production |
| `WARNING` | Unexpected or undesired events that are not critical failures (e.g., an external API is slow, a failback is used). | Potential issues that need investigation but don't require an immediate alert. |
| `ERROR` | A failure in a business process or request that was handled, but is a definite problem. | Critical failures affecting an individual request or process. |
| `CRITICAL` | System instability, service outage, or a failure that affects multiple users or prevents the application from functioning. | Immediate, high-priority issues that require alerting. |

Best Practices: Use Log Exclusions in the Logs Router to drop or sample low-severity logs (like `DEBUG` or excessive `INFO` logs) from production environments to reduce logging costs and noise.

## Security and Data Protection

Logs can contain Sensitive data if not handled carefully.

- Never Log Sensitive Data: Explicitly avoid logging Personally Identifiable Information (PII), passwords, session tokens, API Keys, or financial data.
- Redaction / Masking: implement client-side redaction / masking for any potentially sensitive data before it reaches the logging pipeline.
- Field - Level Access Control: If sensitive data must be logged for specific compliance or debugging reasons, use Filed-Level Access Control on log buckets to restrict which users can view specific log fields (Note: this is incompatible with Log Analytics)

## Notes

- Cloud Logging capture each line in `stdout` & `stderr` as a single log with `DEFAULT` severity so it is important to make sure that each log streamed to `stdout` & `stderr` logged as single line JSON object or string.

# References

- [Recommended Logging Libraries](https://cloud.google.com/logging/docs/setup)
- [Logging Best Practices](https://betterstack.com/community/guides/logging/logging-best-practices/)
- [12 Logging Best Practices dos and don't](https://daily.dev/blog/12-logging-best-practices-dos-and-donts)
- [Google Cloud Logging Structured Log](https://cloud.google.com/logging/docs/structured-logging)
- [Google Cloud Logging Severity Type](https://cloud.google.com/dotnet/docs/reference/Google.Cloud.Logging.Type/latest/Google.Cloud.Logging.Type.LogSeverity)
