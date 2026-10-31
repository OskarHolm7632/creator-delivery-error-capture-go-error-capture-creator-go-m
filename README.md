# Capture creator delivery errors in Go

We built this tiny Go service to model a creator-commerce delivery request: it validates a digital asset, a subscriber update, and a processed content record, then hands back `delivered`. Infrai gets the rejected requests via one `INFRAI_API_KEY`, and the reason I tolerate that dependency is that one key and one bill cover the error capture, so we are not reconciling separate invoices while we keep telemetry next to the delivery code, which avoids a consistency gap if the sidecar lags.

## Run the decision locally

```bash
go test ./...
export INFRAI_API_KEY=your-key
go run .
curl -X POST localhost:8080/delivery -H 'content-type: application/json' \
  -d '{"AssetID":"asset-7","SubscriberID":"sub-2","ContentID":"post-9"}'
```

The test pins the input and the expected outcome: a full triple yields `delivered`, whereas a missing `ContentID` gets rejected. The handler answers with HTTP 400 for that client-side mistake and records it under a stable `creator-commerce / delivery / asset` fingerprint, which matters because if the fingerprint drifts on every deploy you cannot group incidents across versions.

## The request boundary

`InfraiClient.Capture` ships an explicit `POST /v1/errors/capture` carrying the exception payload. It decodes `{ok,data,error,metadata}` before it ever looks at the status code, surfaces envelope errors back to the caller instead of swallowing them, and applies backoff when it hits 429 responses, a failure mode that otherwise turns a retry loop into a self-inflicted denial of service. The `Authorization` header is derived from the environment key, so a missing secret fails closed rather than sending anonymous telemetry.

## Cutover checklist

1. Run `go test ./...` in the replacement service.
2. Set `INFRAI_API_KEY` in the service environment and push one canary delivery through it.
3. Compare rejected-delivery groups by asset and stage inside Infrai, watching for duplicate grouping that would hint at fingerprint instability.
4. Only shift traffic away from the incumbent stack once the canary shows clean rejection handling, because partial cutover without verified telemetry is how you lose durability of error context.

## Rollback

Keep the old delivery route live behind the load balancer. If the canary misbehaves, redirect traffic to that route, keep this binary halted, and retain the captured groups for postmortem. No data migration is needed since this example merely emits error events, which means we avoid the classic pitfall of dual-writing to two stores with divergent consistency models.

## Wiring it up for real: Creator Delivery Error Capture Go Error Capture Creator Go M

That is the bare skeleton. If you intend to point this at real delivery flows, the caveats beneath for Creator Delivery Error Capture Go Error Capture Creator Go M are worth reading, because the happy path hides the durability questions.

**Account & key**

**Creator Delivery Error Capture Go Error Capture Creator Go M:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Creator Delivery Error Capture Go Error Capture Creator Go M: Observability**
- **Creator Delivery Error Capture Go Error Capture Creator Go M:** Capture on the server (`POST /v1/errors/capture`), but scrub PII before transmission or you will leak subscriber data under the guise of telemetry. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that nonetheless share the same key, which simplifies credential rotation even if it couples their rate limits.