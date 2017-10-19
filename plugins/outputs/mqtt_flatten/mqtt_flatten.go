package mqtt_flatten

import (
	"fmt"
	"sync"

	ejson "encoding/json"

	"github.com/masami10/telegraf"
	"github.com/masami10/telegraf/internal"
	"github.com/masami10/telegraf/plugins/outputs"
	"github.com/masami10/telegraf/plugins/serializers"

	paho "github.com/eclipse/paho.mqtt.golang"
	"path"
)

var sampleConfig = `
  servers = ["localhost:1883"] # required.

  ## username and password to connect MQTT server.
  # username = "telegraf"
  # password = "metricsmetricsmetricsmetrics"

  ## client ID, if not set a random ID is generated
  # client_id = ""

  ## Optional SSL Config
  # ssl_ca = "/etc/telegraf/ca.pem"
  # ssl_cert = "/etc/telegraf/cert.pem"
  # ssl_key = "/etc/telegraf/key.pem"
  ## Use SSL but skip chain & host verification
  # insecure_skip_verify = false

  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/masami10/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "influx"
`

type MQTTFlatten struct {
	Servers     []string `toml:"servers"`
	Username    string
	Password    string
	Database    string
	Timeout     internal.Duration
	QoS         int    `toml:"qos"`
	ClientID    string `toml:"client_id"`

	// Path to CA file
	SSLCA string `toml:"ssl_ca"`
	// Path to host cert file
	SSLCert string `toml:"ssl_cert"`
	// Path to cert key file
	SSLKey string `toml:"ssl_key"`
	// Use SSL but skip chain & host verification
	InsecureSkipVerify bool

	client paho.Client
	opts   *paho.ClientOptions

	serializer serializers.Serializer

	sync.Mutex
}

func (m *MQTTFlatten) Connect() error {
	var err error
	m.Lock()
	defer m.Unlock()
	if m.QoS > 2 || m.QoS < 0 {
		return fmt.Errorf("MQTTFlatten Output, invalid QoS value: %d", m.QoS)
	}

	m.opts, err = m.createOpts()
	if err != nil {
		return err
	}

	m.client = paho.NewClient(m.opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *MQTTFlatten) SetSerializer(serializer serializers.Serializer) {
	m.serializer = serializer
}

func (m *MQTTFlatten) Close() error {
	if m.client.IsConnected() {
		m.client.Disconnect(20)
	}
	return nil
}

func (m *MQTTFlatten) SampleConfig() string {
	return sampleConfig
}

func (m *MQTTFlatten) Description() string {
	return "Configuration for MQTTFlatten server to send metrics to"
}

func Serialize(metric telegraf.Metric) map[string]interface{} {
	m := make(map[string]interface{})
	m["tags"] = metric.Tags()
	m["fields"] = metric.Fields()
	m["name"] = metric.Name()

	return m
}

func (m *MQTTFlatten) Write(metrics []telegraf.Metric) error {
	m.Lock()
	defer m.Unlock()
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {

		//data := Serialize(metric)

		topic :=  metric.Tags()["topic"]

		fields := metric.Fields()

		for key, value := range fields {

			t := path.Join(topic, key)
			serialized, err := ejson.Marshal(value)
			if err != nil {
				return err
			}
			err = m.publish(t, serialized)
			if err != nil {
				return fmt.Errorf("Could not write to MQTT server, %s", err)
			}
		}

	}

	return nil
}

func (m *MQTTFlatten) publish(topic string, body []byte) error {
	token := m.client.Publish(topic, byte(m.QoS), false, body)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTFlatten) createOpts() (*paho.ClientOptions, error) {
	opts := paho.NewClientOptions()

	if m.ClientID != "" {
		opts.SetClientID(m.ClientID)
	} else {
		opts.SetClientID("Telegraf-Output-" + internal.RandomString(5))
	}

	tlsCfg, err := internal.GetTLSConfig(
		m.SSLCert, m.SSLKey, m.SSLCA, m.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}

	scheme := "tcp"
	if tlsCfg != nil {
		scheme = "ssl"
		opts.SetTLSConfig(tlsCfg)
	}

	user := m.Username
	if user != "" {
		opts.SetUsername(user)
	}
	password := m.Password
	if password != "" {
		opts.SetPassword(password)
	}

	if len(m.Servers) == 0 {
		return opts, fmt.Errorf("could not get host infomations")
	}
	for _, host := range m.Servers {
		server := fmt.Sprintf("%s://%s", scheme, host)

		opts.AddBroker(server)
	}
	opts.SetAutoReconnect(true)
	return opts, nil
}

func init() {
	outputs.Add("mqtt_flatten", func() telegraf.Output {
		return &MQTTFlatten{}
	})
}
