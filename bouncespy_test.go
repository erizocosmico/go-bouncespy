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

var msg5 = "--001a114118acf80c3205277d9377\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nDelivery to the following recipient failed permanently:\r\n\r\n     foo@foo.se\r\n\r\nTechnical details of permanent failure:\r\nGoogle tried to deliver your message, but it was rejected by the\r\nserver for the recipient domain foo.se by aspmx.l.google.com.\r\n[2a00:1450:400c:c0a::1a].\r\n\r\nThe error that the other server returned was:\r\n550-5.1.1 The email account that you tried to reach does not exist. Please try\r\n550-5.1.1 double-checking the recipient's email address for typos or\r\n550-5.1.1 unnecessary spaces. Learn more at\r\n550 5.1.1  https://support.google.com/mail/answer/6596\r\nus2si56046121wjc.170 - gsmtp\r\n\r\n\r\n----- Original message -----\r\n\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\r\n        d=tyba.com; s=google;\r\n        h=message-id:mime-version:date:from:to:subject:content-type\r\n         :content-transfer-encoding;\r\n        bh=MwWkcFnO/EDKbeKgUqiNEQ4tzRilRsASOD9zlXracn4=;\r\n        b=AZRjeeGBJ6idmtTHVOoSuBrUg7DULBjoIuNriXuqDCM9HH4gk1NeD/yGKMCe4thywE\r\n         yvURCf2Mt4L2J8sh6ciGTRE9+1MK7vKclyWuu/oIsQE0+PFP4y4NL1tvBcUhsrgbSdjN\r\n         EzPBqnzmoQY9Ciu/0dEPkfr4apnXPJghxxeLs=\r\nX-Google-DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\r\n        d=1e100.net; s=20130820;\r\n        h=x-gm-message-state:message-id:mime-version:date:from:to:subject\r\n         :content-type:content-transfer-encoding;\r\n        bh=MwWkcFnO/EDKbeKgUqiNEQ4tzRilRsASOD9zlXracn4=;\r\n        b=f8IbTkx4+M4rGnno6iwjNFJYM0LhMWESWxRrm6MmKAm8FIZoaSGNm25ooe4O8t0+n5\r\n         SJceFIPX2zl2bgPT5uxFVk/wbqQOZ4N7M4OquLfARCAL/8url2eFTZybJe95qljeEO6D\r\n         WdUagutbchqjrmLHYSqURdkOD9TtMgM4y1fUACUTlQnlwQkZ1yOa+Wu1AXz4Fwe9VF62\r\n         hobBqg5hsrmVe17XiSESa3vu0oESU7/T1N1nPpuG6h3dmHsdGri7+aPTT1dfYNaIo02k\r\n         04u9qIwosLr/5o8LkCJCRcAOOEvmLNYo+AwsLeXENtxCRj4BB3tCjmDz08odJxW/KcjY\r\n         u4Rw==\r\nX-Gm-Message-State:\r\nALoCoQl96JHjNkheqFScvUG6W4G6hZ2yxKymdH2eb/v8UzXNSwpGYXBN+fBFTyTEKDlxKsOnIFGR8or5vTCibxZFcI2Mub7TGQ==\r\nX-Received: by 10.28.49.65 with SMTP id x62mr27673005wmx.49.1450784677418;\r\n        Tue, 22 Dec 2015 03:44:37 -0800 (PST)\r\nReturn-Path: <foo@foo.foo>\r\nReceived: [1.1.1.1])\r\n        by smtp.gmail.com with ESMTPSA id\r\no132sm19541744wmb.7.2015.12.22.03.44.35\r\n        for <christian@cpf.se>\r\n        (version=TLSv1/SSLv3 cipher=OTHER);\r\n        Tue, 22 Dec 2015 03:44:36 -0800 (PST)\r\nMessage-ID: <567937a4.8a5a1c0a.d2aae.ffffc555@mx.google.com>\r\nMime-Version: 1.0\r\nDate: Tue, 22 Dec 2015 12:45:14 +0100\r\n\r\n----- End of message -----\r\n\r\n"

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
	}

	for _, cs := range cases {
		c.Assert(FindBounceReason([]byte(cs.msg)), Equals, cs.r)
	}
}
