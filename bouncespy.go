package bouncespy

import (
	"net/mail"
	"strconv"
	"strings"
)

// BounceType defines the type of the bounce, that is, hard or soft
type BounceType int

const (
	Soft BounceType = 0
	Hard BounceType = 1
)

// BounceReason is a status code that tells why the message was bounced according to
// https://tools.ietf.org/html/rfc3463#section-3 and https://tools.ietf.org/html/rfc821#section-4.2.2
type BounceReason string

const (
	ServiceNotAvailable               BounceReason = "421"
	MailActionNotTaken                BounceReason = "450"
	ActionAbortedErrorProcessing      BounceReason = "451"
	ActionAbortedInsufficientStorage  BounceReason = "452"
	CmdSyntaxError                    BounceReason = "500"
	ArgumentsSyntaxError              BounceReason = "501"
	CmdNotImplemented                 BounceReason = "502"
	BadCmdSequence                    BounceReason = "503"
	CmdParamNotImplemented            BounceReason = "504"
	MailboxUnavailable                BounceReason = "550"
	RecipientNotLocal                 BounceReason = "551"
	ActionAbortedExceededStorageAlloc BounceReason = "552"
	MailboxNameInvalid                BounceReason = "553"
	TransactionFailed                 BounceReason = "554"

	AddressDoesntExist                       BounceReason = "5.0.0"
	OtherAddressError                        BounceReason = "5.1.0"
	BadDestinationMailboxAddress             BounceReason = "5.1.1"
	BadDestinationSystemAddress              BounceReason = "5.1.2"
	BadDestinationMailboxAddressSyntax       BounceReason = "5.1.3"
	DestinationMailboxAmbiguous              BounceReason = "5.1.4"
	DestinationMailboxAddressInvalid         BounceReason = "5.1.5"
	MailboxMoved                             BounceReason = "5.1.6"
	BadSenderMailboxAddressSyntax            BounceReason = "5.1.7"
	BadSenderSystemAddress                   BounceReason = "5.1.8"
	UndefinedMailboxError                    BounceReason = "5.2.0"
	MailboxDisabled                          BounceReason = "5.2.1"
	MailboxFull                              BounceReason = "5.2.2"
	MessageLenExceedsLimit                   BounceReason = "5.2.3"
	MailingListExpansionProblem              BounceReason = "5.2.4"
	UndefinedMailSystemStatus                BounceReason = "5.3.0"
	MailSystemFull                           BounceReason = "5.3.1"
	SystemNotAcceptingNetworkMessages        BounceReason = "5.3.2"
	SystemNotCapableOfFeatures               BounceReason = "5.3.3"
	MessageTooBigForSystem                   BounceReason = "5.3.4"
	UndefinedNetworkStatus                   BounceReason = "5.4.0"
	NoAnswerFromHost                         BounceReason = "5.4.1"
	BadConnection                            BounceReason = "5.4.2"
	RoutingServerFailure                     BounceReason = "5.4.3"
	UnableToRoute                            BounceReason = "5.4.4"
	NetworkCongestion                        BounceReason = "5.4.5"
	RoutingLoopDetected                      BounceReason = "5.4.6"
	DeliveryTimeExpired                      BounceReason = "5.4.7"
	UndefinedProtocolStatus                  BounceReason = "5.5.0"
	InvalidCommand                           BounceReason = "5.5.1"
	SyntaxError                              BounceReason = "5.5.2"
	TooManyRecipients                        BounceReason = "5.5.3"
	InvalidCommandArguments                  BounceReason = "5.5.4"
	WrongProtocolVersion                     BounceReason = "5.5.5"
	UndefinedMediaError                      BounceReason = "5.6.0"
	MediaNotSupported                        BounceReason = "5.6.1"
	ConversionRequiredAndProhibited          BounceReason = "5.6.2"
	ConversionRequiredButNotSupported        BounceReason = "5.6.3"
	ConversionWithLossPerformed              BounceReason = "5.6.4"
	ConversionFailed                         BounceReason = "5.6.5"
	UndefinedSecurityStatus                  BounceReason = "5.7.0"
	MessageRefused                           BounceReason = "5.7.1"
	MailingListExpansionProhibited           BounceReason = "5.7.2"
	SecurityConversionRequiredButNotPossible BounceReason = "5.7.3"
	SecurityFeaturesNotSupported             BounceReason = "5.7.4"
	CryptoFailure                            BounceReason = "5.7.5"
	CryptoAlgorithmNotSupported              BounceReason = "5.7.6"
	MessageIntegrityFailure                  BounceReason = "5.7.7"
	UndefinedCode                            BounceReason = "9.1.1"

	// NotFound means we did not found the reason in the email
	NotFound BounceReason = ""
)

// StatusMap is a map indexed by bounce reason that returns an object with
// its bounce type and whether it's an specific error or not (an enhanced)
var StatusMap = map[BounceReason]struct {
	Type     BounceType
	Specific bool
}{
	ServiceNotAvailable:               {Soft, false},
	MailActionNotTaken:                {Soft, false},
	ActionAbortedErrorProcessing:      {Soft, false},
	ActionAbortedInsufficientStorage:  {Soft, false},
	CmdSyntaxError:                    {Hard, false},
	ArgumentsSyntaxError:              {Hard, false},
	CmdNotImplemented:                 {Hard, false},
	BadCmdSequence:                    {Hard, false},
	CmdParamNotImplemented:            {Hard, false},
	MailboxUnavailable:                {Hard, false},
	RecipientNotLocal:                 {Hard, false},
	ActionAbortedExceededStorageAlloc: {Hard, false},
	MailboxNameInvalid:                {Hard, false},
	TransactionFailed:                 {Hard, false},

	AddressDoesntExist:                       {Hard, true},
	OtherAddressError:                        {Hard, true},
	BadDestinationMailboxAddress:             {Hard, true},
	BadDestinationSystemAddress:              {Hard, true},
	BadDestinationMailboxAddressSyntax:       {Hard, true},
	DestinationMailboxAmbiguous:              {Hard, true},
	DestinationMailboxAddressInvalid:         {Hard, true},
	MailboxMoved:                             {Hard, true},
	BadSenderMailboxAddressSyntax:            {Hard, true},
	BadSenderSystemAddress:                   {Hard, true},
	UndefinedMailboxError:                    {Soft, true},
	MailboxDisabled:                          {Soft, true},
	MailboxFull:                              {Soft, true},
	MessageLenExceedsLimit:                   {Hard, true},
	MailingListExpansionProblem:              {Hard, true},
	UndefinedMailSystemStatus:                {Hard, true},
	MailSystemFull:                           {Soft, true},
	SystemNotAcceptingNetworkMessages:        {Hard, true},
	SystemNotCapableOfFeatures:               {Hard, true},
	MessageTooBigForSystem:                   {Hard, true},
	UndefinedNetworkStatus:                   {Hard, true},
	NoAnswerFromHost:                         {Hard, true},
	BadConnection:                            {Hard, true},
	RoutingServerFailure:                     {Hard, true},
	UnableToRoute:                            {Hard, true},
	NetworkCongestion:                        {Soft, true},
	RoutingLoopDetected:                      {Hard, true},
	DeliveryTimeExpired:                      {Hard, true},
	UndefinedProtocolStatus:                  {Hard, true},
	InvalidCommand:                           {Hard, true},
	SyntaxError:                              {Hard, true},
	TooManyRecipients:                        {Soft, true},
	InvalidCommandArguments:                  {Hard, true},
	WrongProtocolVersion:                     {Hard, true},
	UndefinedMediaError:                      {Hard, true},
	MediaNotSupported:                        {Hard, true},
	ConversionRequiredAndProhibited:          {Hard, true},
	ConversionRequiredButNotSupported:        {Hard, true},
	ConversionWithLossPerformed:              {Hard, true},
	ConversionFailed:                         {Hard, true},
	UndefinedSecurityStatus:                  {Hard, true},
	MessageRefused:                           {Hard, true},
	MailingListExpansionProhibited:           {Hard, true},
	SecurityConversionRequiredButNotPossible: {Hard, true},
	SecurityFeaturesNotSupported:             {Hard, true},
	CryptoFailure:                            {Hard, true},
	CryptoAlgorithmNotSupported:              {Hard, true},
	MessageIntegrityFailure:                  {Hard, true},
	UndefinedCode:                            {Hard, true},
}

const (
	LessSpecific = -1
	MoreSpecific = 1
	Equal        = 0
	BothNotFound = -2
)

// Compare returns a comparison code depending on the relation between the two
// reasons to compare.
//
// Let A be the compared reason and B the reason to compare it to
// - If A and B are NotFound, the result is BothNotFound
// - If A and B are both enhanced reasons, the result is Equal
// - If A is enhanced but B is not, the result is MoreSpecific
// - If As is not enhanced but B is, the result is LessSpecific
func (r BounceReason) Compare(o BounceReason) int {
	if r == NotFound && o == NotFound {
		return BothNotFound
	}

	infoSelf := StatusMap[r]
	infoOther := StatusMap[o]
	if infoSelf.Specific == infoOther.Specific {
		return Equal
	} else if infoSelf.Specific {
		return MoreSpecific
	} else {
		return LessSpecific
	}
}

// Result is the returned value of the analysis. It contains the bounce type, the reason,
// and the spam score if it was present.
type Result struct {
	Type      BounceType
	Reason    BounceReason
	SpamScore float64
}

// Analyze returns a Result given the headers and body of an email message
func Analyze(headers mail.Header, body []byte) Result {
	reason := FindBounceReason(body)
	return Result{
		SpamScore: SpamScore(headers),
		Reason:    reason,
		Type:      StatusMap[reason].Type,
	}
}

// SpamScore finds the spam score given the email headers
func SpamScore(headers mail.Header) float64 {
	score := headers.Get("X-Spam-Score")
	if score == "" {
		return .0
	}

	scoreNum, err := strconv.ParseFloat(score, 64)
	if err != nil {
		return .0
	}

	return scoreNum
}

const (
	errorOtherServerReturned = "the error that the other server returned was:"
	reasonOfTheProblem       = "the reason of the problem:"
	reasonForTheProblem      = "the reason for the problem:"
)

// FindBounceReason returns the bounce reason found in the body of the email if it was found
func FindBounceReason(body []byte) BounceReason {
	lns := strings.Split(strings.ToLower(string(body)), "\n")
	numLines := len(lns)

	// we need to reverse the lines because some servers send a bounce email
	// with a more specific error code in the end of the message and a less
	// specific error at the beginning
	var lines = make([]string, numLines)
	for i, ln := range lns {
		lines[numLines-i-1] = ln
	}

	for i, line := range lines {
		if strings.HasPrefix(line, "status:") {
			if reason := analyzeLine(line[6:]); reason != NotFound {
				return reason
			}
		}

		if (strings.HasSuffix(line, reasonOfTheProblem) ||
			strings.HasSuffix(line, reasonForTheProblem) ||
			strings.HasSuffix(line, errorOtherServerReturned)) && i-1 >= 0 {
			if reason := analyzeLine(lines[i-1]); reason != NotFound {
				return reason
			}
		}
	}

	return NotFound
}

func analyzeLine(line string) BounceReason {
	var firstStatus, secondStatus BounceReason
	parts := strings.Split(removeUnnecessaryChars(line), " ")

	if len(parts) > 1 {
		secondStatus = parseStatus(parts[1])
	}

	if len(parts) > 0 {
		firstStatus = parseStatus(parts[0])
	}

	switch firstStatus.Compare(secondStatus) {
	case LessSpecific:
		return secondStatus
	case MoreSpecific, Equal:
		return firstStatus
	default:
		return NotFound
	}
}

func parseStatus(status string) BounceReason {
	status = strings.TrimSpace(status)
	if _, ok := StatusMap[BounceReason(status)]; ok {
		return BounceReason(status)
	}
	return NotFound
}

func removeUnnecessaryChars(line string) string {
	return strings.Replace(strings.Replace(strings.Replace(line, "-", " ", -1), "  ", " ", -1), "  ", " ", -1)
}
