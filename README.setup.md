## How to setup mainflux agent v2

### Setup Mainflux API_TOKEN

```bash
API_TOKEN="<Your_Secret_Token>"
```

### Create 2 channels mannually

```bash
Control channel: 6631e8e5-df48-44e4-bcc0-6e6191cb1b68
Data channel: 8fefdd5e-ad30-4397-90cf-7b3d94089cad
```



### Create bootstrap config

Got thing_id: /things/configs/931c0052-1d6b-44bf-931c-0b805d5177a7

Prepare your external id and key:

```bash
external_id=00-AA-BB-CD-12-18
external_key=00e0dcfa-6b46-11e9-a923-1681be663d3e
```

Then execute curl to create thing_id by bootstrap server:

```bash
curl -v \
  -H "Authorization: $API_TOKEN" \
  -H "Content-Type: application/json" \
  http://190.190.190.81:8202/things/configs \
  -d '{
    "external_id": "00-AA-BB-CD-12-18",
    "external_key": "00e0dcfa-6b46-11e9-a923-1681be663d3e",
    "channels": [
      "6631e8e5-df48-44e4-bcc0-6e6191cb1b68",
      "8fefdd5e-ad30-4397-90cf-7b3d94089cad"
    ]
  }'
```


### You could find all configs (Optional)

```bash
# Get configs
curl -H "Authorization: $API_TOKEN" \
  -H "Content-Type: application/json" \
  http://190.190.190.81:8202/things/configs
```

### You could check boostrap by external id (Optional)

```bash
curl -s -S -i \
  -H "Authorization: 00e0dcfa-6b46-11e9-a923-1681be663d3e" \
  http://190.190.190.81:8202/things/bootstrap/00-AC-87-B4-86-18
```


### **Enable thing by thing_id (Important!)**

```bash
curl -s -S -i -X PUT -H "Authorization: $API_TOKEN" \
  -H "Content-Type: application/json" \
  http://190.190.190.81:8202/things/state/931c0052-1d6b-44bf-931c-0b805d5177a7 -d '{"state": 1}'
```

### Run mainflux agent

```bash
MF_AGENT_BOOTSTRAP_ID=00-AA-BB-CD-12-18 \
MF_AGENT_BOOTSTRAP_KEY=00e0dcfa-6b46-11e9-a923-1681be663d3e \
MF_AGENT_BOOTSTRAP_URL=http://190.190.190.81:8202/things/bootstrap \
MF_AGENT_MQTT_USERNAME=931c0052-1d6b-44bf-931c-0b805d5177a7 \
MF_AGENT_MQTT_PASSWORD=6e0f0fe5-5a2f-43ef-be5e-a424ba2c7f0c \
MF_AGENT_NATS_URL=nats://190.190.190.81:4222 \
build/mainflux-agent
```


### MQTT publish and subcribe

```bash
# Subscribe to control channel
mosquitto_sub \
  -u 931c0052-1d6b-44bf-931c-0b805d5177a7 -P 6e0f0fe5-5a2f-43ef-be5e-a424ba2c7f0c \
  -t channels/6631e8e5-df48-44e4-bcc0-6e6191cb1b68/messages/res \
  -h 190.190.190.81 -p 1883

# Publish to control channel
mosquitto_pub \
  -u 931c0052-1d6b-44bf-931c-0b805d5177a7 -P 6e0f0fe5-5a2f-43ef-be5e-a424ba2c7f0c \
  -t channels/6631e8e5-df48-44e4-bcc0-6e6191cb1b68/messages/req \
  -i kyya-9527 \
  -h 190.190.190.81 -p 1883 -m '[{"bn":"1:", "n":"config", "vs":"view"}]'  
```