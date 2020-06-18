# STAN-connector

The STAN connector connects OpenFaaS functions to STAN subjects.

## Configuration

This configuration can be set in the YAML files for Kubernetes or Swarm.

| Env var                  | Description                                                                        | Default               |
| ------------------------ | ---------------------------------------------------------------------------------- | --------------------- |
| `broker_host`            | The DNS entry for NATS                                                             | `bus`                 |
| `broker_max_reconnect`   | An integer of the amount of reconnection attempts when the NATS connection is lost | `120`                 |
| `broker_reconnect_delay` | Delay between retrying to connect to NATS                                          | `2s`                  |
| `broker_durable`         |                                                                                    | ``                    |
| `broker_ack_wait`        |                                                                                    | `30s`                 |
| `broker_max_in_flight`   |                                                                                    | `1`                   |
| `broker_unsubscribe`     |                                                                                    | `false`               |
| `topics`                 | Topics to which the connector will bind                                            | ``                    |
| `gateway_url`            | The URL for the API gateway                                                        | `http://gateway:8080` |
| `rebuild_interval`       | Go duration - interval for rebuilding function to topic map                        | ``                    |
| `print_response`         | This will output information about the response of calling a function in the logs, | `true`                |
                             including the HTTP status, topic that triggered invocation, the function name, 
                             and the length of the response body in bytes

