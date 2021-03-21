# prometheus_relay

Small app to periodically run a [prometheus](https://prometheus.io) query and relay
the results to an HTTP POST endpoint.

I wrote this app to relay some metrics from my own prometheus server to a
simple web app that I use to display the values in various friendly ways.

Example usage:

    prometheus_relay 'http://localhost:9090/prometheus' 'https://mywebserver.com/api/values' 'ambient_wx_temperature' 'location'

This executes the query `ambient_wx_temperature` (get the latest values for the metric named `ambient_wx_temperature`)
using the prometheus HTTP API running at 'http://localhost:9090/prometheus', takes each returned value and posts
a payload to 'https://mywebserver.com/api/values'.  The post payload looks like

    {
        "name": "<value of the 'location' label in the metric>",
        "value": <value of the metric>
    }

I plan to make this more flexible, but so far this format suits my needs.
