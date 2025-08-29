package sms

import (
	"context"
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type KavenegarProvider struct {
	client *kavenegar.Kavenegar
}

func NewKavenegarProvider(apiKey string) Provider {
	return &KavenegarProvider{
		client: kavenegar.New(apiKey),
	}
}

func (p *KavenegarProvider) SendOTP(
	ctx context.Context,
	phone, code, template string,
) error {
	params := &kavenegar.VerifyLookupParam{}

	if _, err := p.client.Verify.Lookup(
		phone,
		template,
		code,
		params,
	); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			return fmt.Errorf(
				"kavenegar API error: %s",
				err.Error(),
			)
		case *kavenegar.HTTPError:
			return fmt.Errorf(
				"kavenegar HTTP error: %s",
				err.Error(),
			)
		default:
			return fmt.Errorf(
				"kavenegar error: %s",
				err.Error(),
			)
		}
	}

	return nil
}
