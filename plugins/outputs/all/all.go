package all

import (
	_ "github.com/masami10/telegraf/plugins/outputs/amon"
	_ "github.com/masami10/telegraf/plugins/outputs/amqp"
	_ "github.com/masami10/telegraf/plugins/outputs/cloudwatch"
	_ "github.com/masami10/telegraf/plugins/outputs/datadog"
	_ "github.com/masami10/telegraf/plugins/outputs/discard"
	_ "github.com/masami10/telegraf/plugins/outputs/elasticsearch"
	_ "github.com/masami10/telegraf/plugins/outputs/file"
	_ "github.com/masami10/telegraf/plugins/outputs/graphite"
	_ "github.com/masami10/telegraf/plugins/outputs/graylog"
	_ "github.com/masami10/telegraf/plugins/outputs/influxdb"
	_ "github.com/masami10/telegraf/plugins/outputs/instrumental"
	_ "github.com/masami10/telegraf/plugins/outputs/kafka"
	_ "github.com/masami10/telegraf/plugins/outputs/kinesis"
	_ "github.com/masami10/telegraf/plugins/outputs/librato"
	_ "github.com/masami10/telegraf/plugins/outputs/mqtt"
	_ "github.com/masami10/telegraf/plugins/outputs/mqtt_flatten"
	_ "github.com/masami10/telegraf/plugins/outputs/nats"
	_ "github.com/masami10/telegraf/plugins/outputs/nsq"
	_ "github.com/masami10/telegraf/plugins/outputs/opentsdb"
	_ "github.com/masami10/telegraf/plugins/outputs/prometheus_client"
	_ "github.com/masami10/telegraf/plugins/outputs/riemann"
	_ "github.com/masami10/telegraf/plugins/outputs/riemann_legacy"
	_ "github.com/masami10/telegraf/plugins/outputs/socket_writer"
)
