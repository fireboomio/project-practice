package customize

import (
	"custom-go/generated"
	"custom-go/pkg/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
	"time"
)

const (
	manyMembershipQueryPath        = generated.Membership__GetManyMembership
	iosRenewalQueryPath            = generated.Payment__GetNearlyIOSRenewal
	timeDay                        = time.Hour * 24
	appleSubscriptionRenewalSwitch = "apple_subscription_renewal_switch"
)

type (
	manyMembershipI  = generated.Membership__GetManyMembershipInput
	manyMembershipRD = generated.Membership__GetManyMembershipResponseData
	iosRenewalI      = generated.Payment__GetNearlyIOSRenewalInternalInput
	iosRenewalRD     = generated.Payment__GetNearlyIOSRenewalResponseData
)

func init() {
	base.AddRegisteredHook(startPaymentAppleSubscription)
}

func startPaymentAppleSubscription(logger echo.Logger) {
	if !strings.EqualFold(os.Getenv(appleSubscriptionRenewalSwitch), "1") {
		return
	}

	internalClient := plugins.DefaultInternalClient
	renewalFunction := func() {
		now := time.Now()
		manyMembershipResp, _ := plugins.ExecuteInternalRequestQueries[manyMembershipI, manyMembershipRD](internalClient, manyMembershipQueryPath, manyMembershipI{})
		grc := &base.GraphqlRequestContext{Logger: logger, InternalClient: internalClient}
		for _, membershipItem := range manyMembershipResp.Data {
			paymentDateGte := now.Add(timeDay * time.Duration(membershipItem.Lifespan) * -1)
			paymentDateLte := paymentDateGte.Add(timeDay * 1)
			iosRenewalInput := iosRenewalI{PaymentDateGte: paymentDateGte.Format(utils.ISO8601Layout), PaymentDateLte: paymentDateLte.Format(utils.ISO8601Layout)}
			iosRenewalResp, _ := plugins.ExecuteInternalRequestQueries[iosRenewalI, iosRenewalRD](internalClient, iosRenewalQueryPath, iosRenewalInput)
			for _, renewalItem := range iosRenewalResp.Data {
				verifyInput := &AppleVerifyInput{Receipt: renewalItem.Sn, AccountId: renewalItem.AccountId, AppleProductId: membershipItem.ProductId}
				_, _ = IosVerify(grc, verifyInput)
			}
		}
	}
	renewalFunction()
	for range time.Tick(time.Hour) {
		renewalFunction()
	}
}
