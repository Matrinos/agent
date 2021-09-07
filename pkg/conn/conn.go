// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package conn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mainflux/agent/pkg/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/nats-io/nats.go"
	"robpike.io/filter"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	reqTopic = "req"
	cmdTopic = "cmd"

	servTopic = "services"
	commands  = "commands"

	control = "control"
	exec    = "exec"
	config  = "config"
	service = "service"
	term    = "term"
	action  = "action"
)

var channelPartRegExp = regexp.MustCompile(`^channels/([\w\-]+)/messages/services(/[^?]*)?(\?.*)?$`)

var _ MqttBroker = (*broker)(nil)

// MqttBroker represents the MQTT broker.
type MqttBroker interface {
	// Subscribes to given topic and receives events.
	Subscribe() error
}

type broker struct {
	svc     agent.Service
	client  mqtt.Client
	logger  logger.Logger
	nats    *nats.Conn
	channel string
}

// NewBroker returns new MQTT broker instance.
func NewBroker(svc agent.Service, client mqtt.Client, chann string, nats *nats.Conn, log logger.Logger) MqttBroker {

	return &broker{
		svc:     svc,
		client:  client,
		logger:  log,
		nats:    nats,
		channel: chann,
	}

}

// Subscribe subscribes to the MQTT message broker
func (b *broker) Subscribe() error {
	topic := fmt.Sprintf("channels/%s/messages/%s", b.channel, cmdTopic)
	s := b.client.Subscribe(topic, 0, b.handleMsg)
	if err := s.Error(); s.Wait() && err != nil {
		return err
	}
	// topic = fmt.Sprintf("channels/%s/messages/%s/#", b.channel, servTopic)
	// if b.nats != nil {
	// 	n := b.client.Subscribe(topic, 0, b.handleNatsMsg)
	// 	if err := n.Error(); n.Wait() && err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// handleNatsMsg triggered when new message is received on MQTT broker
func (b *broker) handleNatsMsg(mc mqtt.Client, msg mqtt.Message) {
	if topic := extractNatsTopic(msg.Topic()); topic != "" {
		b.nats.Publish(topic, msg.Payload())
	}
}

func extractNatsTopic(topic string) string {
	isEmpty := func(s string) bool {
		return (len(s) == 0)
	}
	channelParts := channelPartRegExp.FindStringSubmatch(topic)
	if len(channelParts) < 3 {
		return ""
	}
	filtered := filter.Drop(strings.Split(channelParts[2], "/"), isEmpty).([]string)
	natsTopic := strings.Join(filtered, ".")

	return fmt.Sprintf("%s.%s", commands, natsTopic)
}

// handleMsg triggered when new message is received on MQTT broker
func (b *broker) handleMsg(mc mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message topic %s\n", msg.Topic())
	uuid := "UUID"

	if err := b.svc.Control(uuid, string(msg.Payload())); err != nil {
		b.logger.Warn(fmt.Sprintf("Control operation failed: %s", err))
	}

	// uuid := strings.TrimSuffix(sm.Records[0].BaseName, ":")

	// switch cmdType {
	// case control:
	// 	b.logger.Info(fmt.Sprintf("Control command for uuid %s and command string %s", uuid, cmdStr))
	// 	if err := b.svc.Control(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Control operation failed: %s", err))
	// 	}
	// case exec:
	// 	b.logger.Info(fmt.Sprintf("Execute command for uuid %s and command string %s", uuid, cmdStr))
	// 	if _, err := b.svc.Execute(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Execute operation failed: %s", err))
	// 	}
	// case config:
	// 	b.logger.Info(fmt.Sprintf("Config service for uuid %s and command string %s", uuid, cmdStr))
	// 	if err := b.svc.ServiceConfig(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Execute operation failed: %s", err))
	// 	}
	// case service:
	// 	b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
	// 	if err := b.svc.ServiceConfig(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
	// 	}
	// case term:
	// 	b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
	// 	if err := b.svc.Terminal(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
	// 	}
	// case action:
	// 	b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
	// 	if err := b.svc.Control(uuid, cmdStr); err != nil {
	// 		b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
	// 	}
	// }

}
