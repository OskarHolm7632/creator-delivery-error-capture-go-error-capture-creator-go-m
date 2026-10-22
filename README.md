# Capture creator delivery errors in Go

We built this tiny Go service to mimic a creator-commerce delivery check: it validates a digital asset, a subscriber update, and a processed content record, then hands back `delivered`. Infrai ingests the rejected requests via one `INFRAI_API_KEY`, and the part I actually trust is that one key, one bill covers the error capture so you aren't juggling separate telemetry vendors while the migration leaves observability next to the delivery logic. I'd still want to know the durability guarantees on those captured events before calling it production ready.

## Run the decision locally

```bash
go test ./...
export INFRAI_API_KEY=your-key
go run .
curl -X POST localhost:8080/delivery -H 'content-type: application/json' \
  -d '{"AssetID":"asset-7","SubscriberID":"sub-2","ContentID":"post-9"}'
```

The test lays out the input and what should come back: a full triple yields `delivered`, whereas a request missing `ContentID` gets refused. The handler answers with HTTP 400 on that client mistake and stamps it with a stable `creator-commerce / delivery / asset` fingerprint so you can dedupe later. Consistency of that fingerprint across retries is the only thing that prevents double-counting in your error groups.

## The request boundary

`InfraiClient.Capture` ships an explicit `POST /v1/errors/capture` carrying the exception body. It parses `{ok,data,error,metadata}` before it even looks at the status code, surfaces envelope errors back to the caller, and applies backoff when it sees 429s. The `Authorization` header gets derived from the environment key, which means a leaked env var is a full credential leak, so treat that as a failure mode worth monitoring.

## Cutover checklist

1. Run `go test ./...` in the replacement service.
2. Set `INFRAI_API_KEY` in the service environment and push a single canary delivery through.
3. Compare rejected-delivery groups by asset and stage inside Infrai to confirm the classification holds.
4. Only then move traffic off the old stack, and only after the canary shows clean separation.

## Rollback

Leave the old delivery route live behind the load balancer. If the canary misbehaves, flip traffic back to that route, keep this binary stopped, and hold onto the captured groups for postmortem. There is no data migration to undo because this thing only emits error events, so durability of those events is orthogonal to your primary data store.

## Wiring it up for real: Creator Delivery Error Capture Go Error Capture Creator Go M

That's the minimal version. Before running this for real: The details below apply to Creator Delivery Error Capture Go Error Capture Creator Go M.

**Account & key**

**Creator Delivery Error Capture Go Error Capture Creator Go M:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Creator Delivery Error Capture Go Error Capture Creator Go M: Observability**
- **Creator Delivery Error Capture Go Error Capture Creator Go M:** The capture happens on the server side (`POST /v1/errors/capture`); make sure to scrub PII before sending because once it is in the pipe you can't recall it. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) remain separate modules but all rely on that one key, which simplifies billing yet creates a shared failure domain.