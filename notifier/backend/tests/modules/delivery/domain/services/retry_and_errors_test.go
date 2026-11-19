package services

import (
    "errors"
    "net"
    "net/textproto"
    "testing"

    "notifier-backend/tests"

    "notifier-backend/src/modules/delivery/domain/exceptions"
    "notifier-backend/src/modules/delivery/domain/services"
    "notifier-backend/src/modules/delivery/infrastructure/components/attachmentdownloader"
)

func TestRetryDelaysGrowAndStopAtMax(t *testing.T) {
    retryScheduler := services.RetrySchedulerGetInstance()
    tests.AssertEquals(t, 60, retryScheduler.GetDelaySeconds(1))
    tests.AssertEquals(t, 300, retryScheduler.GetDelaySeconds(2))
    tests.AssertEquals(t, 7200, retryScheduler.GetDelaySeconds(4))
    tests.AssertEquals(t, 7200, retryScheduler.GetDelaySeconds(9))
    tests.AssertEquals(t, false, retryScheduler.IsLastAttempt(4, 5))
    tests.AssertEquals(t, true, retryScheduler.IsLastAttempt(5, 5))
}

func TestSmtpErrorsAre5xxPermanentAndTheRestTemporary(t *testing.T) {
    classifier := services.SmtpErrorClassifierGetInstance()
    tests.AssertEquals(t, true, classifier.GetDeliveryException(&textproto.Error{Code: 535, Msg: "bad credentials"}).IsPermanent())
    tests.AssertEquals(t, false, classifier.GetDeliveryException(&textproto.Error{Code: 421, Msg: "try later"}).IsPermanent())
    tests.AssertEquals(t, false, classifier.GetDeliveryException(errors.New("connection reset")).IsPermanent())
    tests.AssertEquals(t, true, classifier.GetDeliveryException(exceptions.PermanentCustom("x")).IsPermanent())
}

func TestOnlyPublicIpsAreDownloadable(t *testing.T) {
    for _, privateIp := range []string{"127.0.0.1", "10.1.2.3", "172.16.0.1", "192.168.1.1", "169.254.169.254", "100.64.0.1", "::1", "fd00::1", "0.0.0.0"} {
        tests.AssertEquals(t, false, attachmentdownloader.IsPublicIp(net.ParseIP(privateIp)), privateIp)
    }
    tests.AssertEquals(t, true, attachmentdownloader.IsPublicIp(net.ParseIP("8.8.8.8")))
}
