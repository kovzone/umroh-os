// payment_reissue_messages.go — hand-written proto message types for ReissuePaymentLink RPC.
//
// BL-PAY-020: CS closing — Customer Service can re-issue a VA link for an
// existing booking whose VA may have expired or never been claimed.
//
// Unlike IssueVirtualAccount (called by booking-svc's submit saga), this RPC:
//   - Only works on bookings that already have an invoice (returns ErrNotFound otherwise)
//   - Returns the existing active VA if one is alive (idempotent no-op)
//   - Creates a fresh VA on the same invoice if the existing one is expired/consumed

package pb

// ReissuePaymentLinkRequest is sent by the CS console (via gateway-svc) to
// re-issue or retrieve the active VA link for an existing booking.
type ReissuePaymentLinkRequest struct {
	// BookingId is the booking.bookings.id for which to (re-)issue a VA.
	BookingId string
	// BankCode is the preferred bank for the VA (empty = gateway default).
	BankCode string
	// GatewayPref selects the gateway ("midtrans"|"xendit"|"mock"); empty = primary.
	GatewayPref string
	// ActorUserId is the IAM user triggering this operation (for audit log).
	ActorUserId string
}

func (x *ReissuePaymentLinkRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *ReissuePaymentLinkRequest) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *ReissuePaymentLinkRequest) GetGatewayPref() string {
	if x == nil {
		return ""
	}
	return x.GatewayPref
}
func (x *ReissuePaymentLinkRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

// ReissuePaymentLinkResponse carries the (re-)issued VA details.
type ReissuePaymentLinkResponse struct {
	InvoiceId     string
	BookingId     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     string // RFC3339
	Gateway       string
	// IsNew is true when a fresh VA was created; false when an existing active
	// VA was returned unchanged.
	IsNew bool
}

func (x *ReissuePaymentLinkResponse) GetInvoiceId() string {
	if x == nil {
		return ""
	}
	return x.InvoiceId
}
func (x *ReissuePaymentLinkResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *ReissuePaymentLinkResponse) GetAccountNumber() string {
	if x == nil {
		return ""
	}
	return x.AccountNumber
}
func (x *ReissuePaymentLinkResponse) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *ReissuePaymentLinkResponse) GetAmountTotal() float64 {
	if x == nil {
		return 0
	}
	return x.AmountTotal
}
func (x *ReissuePaymentLinkResponse) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}
func (x *ReissuePaymentLinkResponse) GetGateway() string {
	if x == nil {
		return ""
	}
	return x.Gateway
}
func (x *ReissuePaymentLinkResponse) GetIsNew() bool {
	if x == nil {
		return false
	}
	return x.IsNew
}
