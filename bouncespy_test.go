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

func (s *BounceSuite) TestFindBounceReason(c *ch.C) {
	cases := []struct {
		msg string
		r   BounceReason
	}{
		{msg1, MailboxUnavailable},
		{msg2, AddressDoesntExist},
		{msg3, OtherAddressError},
		{msg4, NotFound},
	}

	for _, cs := range cases {
		c.Assert(FindBounceReason([]byte(cs.msg)), Equals, cs.r)
	}
}
