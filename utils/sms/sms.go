package sms

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kavenegar/kavenegar-go"
	//rest_error "gitlab.com/irmitcod/go-technician/src/rest_errors"
)

var (
	SMS smsClientInterface = &smsClient{}
)

type smsClient struct {
	c *resty.Client
}

func (s *smsClient) setClient(r *resty.Request) {
	panic("implement me")
}

type smsClientInterface interface {
	SendSms(phone string, sms string) error
	setClient(*resty.Request)
}

func init() {
	//c := resty.New()
	//smsClient.setClient(c.R())
}

func (s *smsClient) SendSms(phone string, sms string) error {

	api := kavenegar.New("33694E584334586A5839734B667878323661384254773D3D")

	receptor := "09333033375"
	template := "verify"
	token := "1234"
	params := &kavenegar.VerifyLookupParam{
	}
	if res, err := api.Verify.Lookup(receptor, template, token, params); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
		default:
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("MessageID 	= ", res.MessageID)
		fmt.Println("Status    	= ", res.Status)
		//...
	}



	//sender := "30006703503503"
	//receptor := []string{"09333033375"}
	//message := "Hello Go!"
	//if res, err := api.Message.Send(sender, receptor, message, nil); err != nil {
	//	switch err := err.(type) {
	//	case *kavenegar.APIError:
	//		fmt.Println(err.Error())
	//	case *kavenegar.HTTPError:
	//		fmt.Println(err.Error())
	//	default:
	//		fmt.Println(err.Error())
	//	}
	//} else {
	//	for _, r := range res {
	//		fmt.Println("MessageID 	= ", r.MessageID)
	//		fmt.Println("Status    	= ", r.Status)
	//		//...
	//	}
	//}


	// Create a Resty Client
	//myURL := "https://RestfulSms.com/api/VerificationCode"
	//
	//myHeaders := map[string]string{
	//	"Content-Type":          "application/json",
	//	"x-sms-ir-secure-token": "9dc4890011c914e0aec86309",
	//}
	//dict := map[string]string{
	//	"Code":         sms,
	//	"MobileNumber": phone,
	//}
	//
	//c := resty.New()
	//
	//res, err := c.R().SetBody(dict).SetHeaders(myHeaders).Post(myURL)
	//if err != nil {
	//	return rest_error.NewInternalServerError(err.Error(), errors.New("send sms faild"))
	//
	//}
	//
	//if res.StatusCode() == 200 {
	//	return nil
	//}
	//
	return nil
}
