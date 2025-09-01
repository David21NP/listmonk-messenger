package messenger

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpointsmsvoicev2"
	"github.com/francoispqt/onelog"
)

type endUserMessagingCfg struct {
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Region      string `json:"region"`
	MessageType string `json:"message_type"`
	PoolID      string `json:"pool_id"`
	Log         bool   `json:"log"`
}

type endUserMessagingMessenger struct {
	cfg    endUserMessagingCfg
	client *pinpointsmsvoicev2.PinpointSMSVoiceV2

	logger *onelog.Logger
}

func (p endUserMessagingMessenger) Name() string {
	return "pinpointsmsvoicev2"
}

// Push sends the sms through endUserMessaging API.
func (p endUserMessagingMessenger) Push(msg Message) error {
	phone, ok := msg.Subscriber.Attribs["phone"].(string)
	if !ok {
		return fmt.Errorf("could not find subscriber phone")
	}

	body := string(msg.Body)
	payload := &pinpointsmsvoicev2.SendTextMessageInput{
		DestinationPhoneNumber: &phone,
		MessageBody:            &body,
		OriginationIdentity:    &p.cfg.PoolID, // Phone number or pool ID
		MessageType:            &p.cfg.MessageType,
	}

	out, err := p.client.SendTextMessage(payload)
	if err != nil {
		return err
	}

	if p.cfg.Log {
		p.logger.InfoWith("successfully sent sms").String("phone", phone).String("result", fmt.Sprintf("Message ID: %s", *out.MessageId)).Write()
	}

	return nil
}

func (p endUserMessagingMessenger) Flush() error {
	return nil
}

func (p endUserMessagingMessenger) Close() error {
	return nil
}

// NewEndUserMessaging creates new instance of endUserMessaging
func NewEndUserMessaging(cfg []byte, l *onelog.Logger) (Messenger, error) {
	var c endUserMessagingCfg
	if err := json.Unmarshal(cfg, &c); err != nil {
		return nil, err
	}

	config := &aws.Config{
		MaxRetries: aws.Int(3),
	}
	if c.AccessKey != "" && c.SecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	}
	if c.Region != "" {
		config.Region = &c.Region
	}

	sess := session.Must(session.NewSession(config))
	err := checkCredentials(sess)
	if err != nil {
		return nil, err
	}
	svc := pinpointsmsvoicev2.New(sess)

	return endUserMessagingMessenger{
		client: svc,
		cfg:    c,
		logger: l,
	}, nil
}
