package bouncespy

import (
	"net/mail"
	"testing"

	ch "gopkg.in/check.v1"
)

type BounceSuite struct{}

func Test(t *testing.T) { ch.TestingT(t) }

var _ = ch.Suite(&BounceSuite{})

var Equals = ch.Equals

func (s *BounceSuite) TestSpamScore(c *ch.C) {
	h1 := make(mail.Header)
	h1["X-Spam-Score"] = []string{"-4.0"}

	h2 := make(mail.Header)

	c.Assert(SpamScore(h1), Equals, -4.0)
	c.Assert(SpamScore(h2), Equals, 0.0)
}

func (s *BounceSuite) TestBounceReasonCompare(c *ch.C) {
	cases := []struct {
		r1, r2 BounceReason
		result int
	}{
		{NotFound, NotFound, BothNotFound},
		{ServiceNotAvailable, CryptoFailure, LessSpecific},
		{NotFound, CryptoFailure, LessSpecific},
		{MailboxUnavailable, MailboxUnavailable, Equal},
		{CryptoFailure, CryptoFailure, Equal},
		{AddressDoesntExist, MailboxUnavailable, MoreSpecific},
		{AddressDoesntExist, NotFound, MoreSpecific},
	}

	for _, cs := range cases {
		c.Assert(cs.r1.Compare(cs.r2), Equals, cs.result)
	}
}

func (s *BounceSuite) TestAnalyzeLine(c *ch.C) {
	cases := []struct {
		ln string
		r  BounceReason
	}{
		{"421 a ksk sogjsdhvkfg dk", ServiceNotAvailable},
		{"421- sdfiu a ksk sogjsdhvkfg dk", ServiceNotAvailable},
		{"421- 1.2.3 a ksk sogjsdhvkfg dk", ServiceNotAvailable},
		{"421- 5.0.0 a ksk sogjsdhvkfg dk", AddressDoesntExist},
		{"421 5.0.0 a ksk sogjsdhvkfg dk", AddressDoesntExist},
		{"421 (a ksk sogjsdhvkfg dk)", ServiceNotAvailable},
		{"5.0.0 (a ksk sogjsdhvkfg dk)", AddressDoesntExist},
		{"5.0.0- a ksk sogjsdhvkfg dk", AddressDoesntExist},
		{"a ksk sogjsdhvkfg dk", NotFound},
	}

	for _, cs := range cases {
		c.Assert(analyzeLine(cs.ln), Equals, cs.r)
	}
}

var msg1 = `Delivery to the following recipient failed permanently:

     foo@foo.foo

Technical details of permanent failure: 
Google tried to deliver your message, but it was rejected by the server for the recipient domain foo.foo by mx.foo.foo. [1.1.1.1].

The error that the other server returned was:
550 Account discontinued, cancelled by user


----- Original message -----

Foo

----- End of transmission -----`
var msg2 = `The following message to <foo@foo.foo> was undeliverable.
The reason for the problem:
5.1.0 - Unknown address error 550-'No such user (foo) -ERR foo@foo.foo not found'

Reporting-MTA: dns; mx1.foo.foo

Final-Recipient: rfc822;foo@foo.foo
Action: failed
Status: 5.0.0 (permanent failure)
Remote-MTA: dns; [1.1.1.1]
Diagnostic-Code: smtp; 5.1.0 - Unknown address error 550-'No such user (foo) -ERR foo@foo.foo not found' (delivery attempts: 0)`
var msg3 = `The following message to <foo@foo.foo> was undeliverable.
The reason for the problem:
5.1.0 - Unknown address error 550-'No such user (foo) -ERR foo@foo.foo not found'
`
var msg4 = ``

var msg5 = "--001a114118acf80c3205277d9377\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nDelivery to the following recipient failed permanently:\r\n\r\n     foo@foo.se\r\n\r\nTechnical details of permanent failure:\r\nGoogle tried to deliver your message, but it was rejected by the\r\nserver for the recipient domain foo.se by aspmx.l.google.com.\r\n[2a00:1450:400c:c0a::1a].\r\n\r\nThe error that the other server returned was:\r\n550-5.1.1 The email account that you tried to reach does not exist. Please try\r\n550-5.1.1 double-checking the recipient's email address for typos or\r\n550-5.1.1 unnecessary spaces. Learn more at\r\n550 5.1.1  https://support.google.com/mail/answer/6596\r\nus2si56046121wjc.170 - gsmtp\r\n\r\n\r\n----- Original message -----\r\n\r\nFoo"

var msg6 = "Delivery to the following recipient failed permanently:\r\n\r\n    foo@foo.foo\r\n\r\nTechnical details of permanent failure:=20\r\nDNS Error: Address resolution of foo.foo failed: Domain name not found\r\n\r\n----- Original message -----\r\n\r\nFoo"

var msg7 = "This is an automatically generated Delivery Status Notification\r\n\r\nTHIS IS A WARNING MESSAGE ONLY.\r\n\r\nYOU DO NOT NEED TO RESEND YOUR MESSAGE.\r\n\r\nDelivery to the following recipient has been delayed:\r\n\r\n     foo@foo.foo\r\n\r\nMessage will be retried for 2 more day(s)\r\n\r\nTechnical details of temporary failure:=20\r\nDNS Error: MX lookup of foo.foo returned error DNS server returned =\r\ngeneral failure\r\n\r\n----- Original message -----\r\n\r\nFoo"

var msg8 = "Delivery to the following recipient failed permanently:\r\n\r\n     foo@foo.foo\r\n\r\nTechnical details of permanent failure:=20\r\nThe recipient server did not accept our requests to connect. Learn more at =\r\nhttps://support.google.com/mail/answer/7720=20\r\n[(0) foo@foo. [100.64.174.100]:25: socket error=\r\n]\r\n\r\n----- Original message -----\r\n\r\nFoo"

func (s *BounceSuite) TestFindBounceReason(c *ch.C) {
	cases := []struct {
		msg string
		r   BounceReason
	}{
		{msg1, MailboxUnavailable},
		{msg2, AddressDoesntExist},
		{msg3, OtherAddressError},
		{msg4, NotFound},
		{msg5, BadDestinationMailboxAddress},
		{msg6, UndefinedCode},
		{msg7, ServiceNotAvailable},
		{msg8, UndefinedCode},
	}

	for _, cs := range cases {
		c.Assert(FindBounceReason([]byte(cs.msg)), Equals, cs.r)
	}
}
