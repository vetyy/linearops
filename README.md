# Linearops

Linearops allows to use Linear [Webhooks](https://developers.linear.app/docs/graphql/webhooks)
to trigger OpsGenie [Alerts](https://docs.opsgenie.com/docs/alerts-and-alert-fields).

## Overview

Linearops listens on `<bind.URL.Prefix>/webhook/<lienar.Webhook.ID>` for
example `/linearops/webhook/419cfbaf-b98d-481b-9ad7-661eb432bdbf`.
Linear currently doesn't support any form of authentication, so you can use random webhook ID to get unique URL.

It will parse the webhook data and based on configured Linear workflow `unstarted states` send request to OpsGenie to
trigger an alert.

It will pass title, description and url for the ticket.

After receiving webhook of configure `completed state` it will close the alert.

## Configuration

Linearops can be configured using yaml configuration file or environment variable. When both are provided, environment
variables take precedence.

| yaml | env variable | description | default | required |
| --- | --- | --- | --- | --- |
| log.Level | LINEAROPS_LOG_LEVEL | Logging level | info | no |
| bind.HTTP | LINEAROPS_BIND_HTTP | Address to bind HTTP gateway | :8080 | no |
| bind.URL.Prefix | LINEAROPS_BIND_URL_PREFIX | URL path prefixed for every route except healtcheck. | / | no |
| opsGenie.API.Key | LINEAROPS_OPSGENIE_API_KEY | OpsGenie API Integration key. |  | yes |
| opsGenie.Responders | LINEAROPS_OPSGENIE_RESPONDERS | OpsGenie name of the team to receive alerts. |  | yes |
| linear.Webhook.ID | LINEAROPS_LINEAR_WEBHOOK_ID | Random string that will be used as last part of the webhook path. |  | yes |
| linear.UserAgent | LINEAROPS_LINEAR_USERAGENT | User-agent that will be used to send requests to OpsGenie. | Linear | no |
| linear.Unstarted.States | LINEAROPS_LINEAR_UNSTARTED_STATES | Linear workflow unstarted states that should be used to trigger the alert.. | Reported | no |
| linear.Completed.States | LINEAROPS_LINEAR_COMPLETED_STATES | Linear workflow completed state that should be used to close the alert. | Resolved,Postmortem,Rejected | no |

## API

Healthcheck - `/health` returns:

```
200 OK.
```

Linear Webhook - `<bind.URL.Prefix>/webhook/<lienar.Webhook.ID>` for
example `/linearops/webhook/419cfbaf-b98d-481b-9ad7-661eb432bdbf` returns:

```
200 OK.
```

Who Is OnCall - `/on-call/<schedule_id>` for example `/on-call/ops_team_schedule` returns:

```json
{
  "on_call": [
    "email@domain.com"
  ],
  "next_on_call": [
    "email2@domain.com"
  ]
}
```
