// finance_depth.go — gRPC handlers for Wave 4 Finance depth RPCs
// (BL-FIN-020..041).

package grpc_api

import (
	"context"
	"errors"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/logging"

	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// BL-FIN-020 ScheduleBilling
// ---------------------------------------------------------------------------

func (s *Server) ScheduleBilling(ctx context.Context, req *pb.ScheduleBillingRequest) (*pb.ScheduleBillingResponse, error) {
	const op = "grpc_api.Server.ScheduleBilling"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureID()).Msg("")

	result, err := s.svc.ScheduleBilling(ctx, &service.ScheduleBillingParams{
		DepartureID: req.GetDepartureID(),
		DueDate:     req.GetDueDate(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ScheduleBillingResponse{
		InvoicesCreated: result.InvoicesCreated,
		TotalAmount:     result.TotalAmount,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-021 Bank integration
// ---------------------------------------------------------------------------

func (s *Server) RecordBankTransaction(ctx context.Context, req *pb.RecordBankTransactionRequest) (*pb.RecordBankTransactionResponse, error) {
	const op = "grpc_api.Server.RecordBankTransaction"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("account_id", req.GetAccountID()).Msg("")

	result, err := s.svc.RecordBankTransaction(ctx, &service.RecordBankTransactionParams{
		AccountID:   req.GetAccountID(),
		RefNo:       req.GetRefNo(),
		Amount:      req.GetAmount(),
		TxDate:      req.GetTxDate(),
		Description: req.GetDescription(),
		Direction:   req.GetDirection(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordBankTransactionResponse{TransactionID: result.TransactionID}, nil
}

func (s *Server) GetBankReconciliation(ctx context.Context, req *pb.GetBankReconciliationRequest) (*pb.GetBankReconciliationResponse, error) {
	const op = "grpc_api.Server.GetBankReconciliation"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("account_id", req.GetAccountID()).Msg("")

	result, err := s.svc.GetBankReconciliation(ctx, &service.GetBankReconciliationParams{
		AccountID: req.GetAccountID(),
		StartDate: req.GetStartDate(),
		EndDate:   req.GetEndDate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	rows := make([]*pb.BankTxRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.BankTxRow{
			TxID:       r.TxID,
			RefNo:      r.RefNo,
			Amount:     r.Amount,
			TxDate:     r.TxDate,
			Direction:  r.Direction,
			Reconciled: r.Reconciled,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetBankReconciliationResponse{
		AccountID:      result.AccountID,
		OpeningBalance: result.OpeningBalance,
		ClosingBalance: result.ClosingBalance,
		Rows:           rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-022 AR Subledger
// ---------------------------------------------------------------------------

func (s *Server) GetARSubledger(ctx context.Context, req *pb.GetARSubledgerRequest) (*pb.GetARSubledgerResponse, error) {
	const op = "grpc_api.Server.GetARSubledger"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", req.GetBookingID()).Msg("")

	result, err := s.svc.GetARSubledger(ctx, &service.GetARSubledgerParams{
		BookingID: req.GetBookingID(),
		PilgrimID: req.GetPilgrimID(),
		StartDate: req.GetStartDate(),
		EndDate:   req.GetEndDate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	rows := make([]*pb.ARSubledgerRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.ARSubledgerRow{
			EntryID:     r.EntryID,
			Date:        r.Date,
			Description: r.Description,
			Debit:       r.Debit,
			Credit:      r.Credit,
			Balance:     r.Balance,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetARSubledgerResponse{
		BookingID: result.BookingID,
		PilgrimID: result.PilgrimID,
		Rows:      rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-023 Digital receipts
// ---------------------------------------------------------------------------

func (s *Server) IssueDigitalReceipt(ctx context.Context, req *pb.IssueDigitalReceiptRequest) (*pb.IssueDigitalReceiptResponse, error) {
	const op = "grpc_api.Server.IssueDigitalReceipt"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", req.GetBookingID()).Msg("")

	result, err := s.svc.IssueDigitalReceipt(ctx, &service.IssueDigitalReceiptParams{
		BookingID: req.GetBookingID(),
		PaymentID: req.GetPaymentID(),
		Amount:    req.GetAmount(),
		Notes:     req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.IssueDigitalReceiptResponse{
		ReceiptID:     result.ReceiptID,
		ReceiptNumber: result.ReceiptNumber,
		IssuedAt:      result.IssuedAt,
	}, nil
}

func (s *Server) GetDigitalReceipt(ctx context.Context, req *pb.GetDigitalReceiptRequest) (*pb.GetDigitalReceiptResponse, error) {
	const op = "grpc_api.Server.GetDigitalReceipt"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("receipt_id", req.GetReceiptID()).Msg("")

	result, err := s.svc.GetDigitalReceipt(ctx, &service.GetDigitalReceiptParams{
		ReceiptID: req.GetReceiptID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrNotFound) {
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetDigitalReceiptResponse{
		Receipt: &pb.DigitalReceiptProto{
			ReceiptID:     result.ReceiptID,
			ReceiptNumber: result.ReceiptNumber,
			BookingID:     result.BookingID,
			PaymentID:     result.PaymentID,
			Amount:        result.Amount,
			IssuedAt:      result.IssuedAt,
			Notes:         result.Notes,
		},
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-024 Manual payment
// ---------------------------------------------------------------------------

func (s *Server) RecordManualPayment(ctx context.Context, req *pb.RecordManualPaymentRequest) (*pb.RecordManualPaymentResponse, error) {
	const op = "grpc_api.Server.RecordManualPayment"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", req.GetBookingID()).Msg("")

	result, err := s.svc.RecordManualPayment(ctx, &service.RecordManualPaymentParams{
		BookingID:   req.GetBookingID(),
		Amount:      req.GetAmount(),
		PaymentDate: req.GetPaymentDate(),
		Method:      req.GetMethod(),
		EvidenceURL: req.GetEvidenceURL(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordManualPaymentResponse{
		EntryID:   result.EntryID,
		JournalID: result.JournalID,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-025 Vendor master
// ---------------------------------------------------------------------------

func (s *Server) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (*pb.CreateVendorResponse, error) {
	const op = "grpc_api.Server.CreateVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("name", req.GetName()).Msg("")

	result, err := s.svc.CreateVendor(ctx, &service.CreateVendorParams{
		Name:         req.GetName(),
		Category:     req.GetCategory(),
		BankAccount:  req.GetBankAccount(),
		TaxID:        req.GetTaxID(),
		ContactEmail: req.GetContactEmail(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.CreateVendorResponse{VendorID: result.VendorID}, nil
}

func (s *Server) UpdateVendor(ctx context.Context, req *pb.UpdateVendorRequest) (*pb.UpdateVendorResponse, error) {
	const op = "grpc_api.Server.UpdateVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", req.GetVendorID()).Msg("")

	result, err := s.svc.UpdateVendor(ctx, &service.UpdateVendorParams{
		VendorID:     req.GetVendorID(),
		Name:         req.GetName(),
		Category:     req.GetCategory(),
		BankAccount:  req.GetBankAccount(),
		ContactEmail: req.GetContactEmail(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrNotFound) {
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.UpdateVendorResponse{VendorID: result.VendorID}, nil
}

func (s *Server) ListVendors(ctx context.Context, req *pb.ListVendorsRequest) (*pb.ListVendorsResponse, error) {
	const op = "grpc_api.Server.ListVendors"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListVendors(ctx, &service.ListVendorsParams{
		Category: req.GetCategory(),
		PageSize: req.GetPageSize(),
		Cursor:   req.GetCursor(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	vendors := make([]*pb.VendorRow, 0, len(result.Vendors))
	for _, v := range result.Vendors {
		vendors = append(vendors, &pb.VendorRow{
			VendorID:     v.VendorID,
			Name:         v.Name,
			Category:     v.Category,
			BankAccount:  v.BankAccount,
			TaxID:        v.TaxID,
			ContactEmail: v.ContactEmail,
			CreatedAt:    v.CreatedAt,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ListVendorsResponse{Vendors: vendors, NextCursor: result.NextCursor}, nil
}

func (s *Server) DeleteVendor(ctx context.Context, req *pb.DeleteVendorRequest) (*pb.DeleteVendorResponse, error) {
	const op = "grpc_api.Server.DeleteVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", req.GetVendorID()).Msg("")

	result, err := s.svc.DeleteVendor(ctx, &service.DeleteVendorParams{VendorID: req.GetVendorID()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.DeleteVendorResponse{Deleted: result.Deleted}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-026 AP Subledger
// ---------------------------------------------------------------------------

func (s *Server) GetAPSubledger(ctx context.Context, req *pb.GetAPSubledgerRequest) (*pb.GetAPSubledgerResponse, error) {
	const op = "grpc_api.Server.GetAPSubledger"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", req.GetVendorID()).Msg("")

	result, err := s.svc.GetAPSubledger(ctx, &service.GetAPSubledgerParams{
		VendorID:  req.GetVendorID(),
		StartDate: req.GetStartDate(),
		EndDate:   req.GetEndDate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	rows := make([]*pb.APSubledgerRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.APSubledgerRow{
			EntryID:     r.EntryID,
			Date:        r.Date,
			Description: r.Description,
			Debit:       r.Debit,
			Credit:      r.Credit,
			Balance:     r.Balance,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetAPSubledgerResponse{VendorID: result.VendorID, Rows: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-027 Payment authorization
// ---------------------------------------------------------------------------

func (s *Server) ListPendingAuthorizations(ctx context.Context, req *pb.ListPendingAuthorizationsRequest) (*pb.ListPendingAuthorizationsResponse, error) {
	const op = "grpc_api.Server.ListPendingAuthorizations"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListPendingAuthorizations(ctx, &service.ListPendingAuthorizationsParams{Level: req.GetLevel()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	items := make([]*pb.AuthorizationRow, 0, len(result.Items))
	for _, it := range result.Items {
		items = append(items, &pb.AuthorizationRow{
			AuthID:      it.AuthID,
			BatchID:     it.BatchID,
			Amount:      it.Amount,
			RequestedBy: it.RequestedBy,
			Level:       it.Level,
			CreatedAt:   it.CreatedAt,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ListPendingAuthorizationsResponse{Items: items}, nil
}

func (s *Server) DecidePaymentAuthorization(ctx context.Context, req *pb.DecidePaymentAuthorizationRequest) (*pb.DecidePaymentAuthorizationResponse, error) {
	const op = "grpc_api.Server.DecidePaymentAuthorization"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("auth_id", req.GetAuthID()).Msg("")

	result, err := s.svc.DecidePaymentAuthorization(ctx, &service.DecidePaymentAuthorizationParams{
		AuthID:   req.GetAuthID(),
		Decision: req.GetDecision(),
		Notes:    req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.DecidePaymentAuthorizationResponse{AuthID: result.AuthID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-028 Petty cash
// ---------------------------------------------------------------------------

func (s *Server) RecordPettyCash(ctx context.Context, req *pb.RecordPettyCashRequest) (*pb.RecordPettyCashResponse, error) {
	const op = "grpc_api.Server.RecordPettyCash"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordPettyCash(ctx, &service.RecordPettyCashParams{
		Amount:      req.GetAmount(),
		Direction:   req.GetDirection(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		Date:        req.GetDate(),
		EvidenceURL: req.GetEvidenceURL(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordPettyCashResponse{EntryID: result.EntryID, RunningBalance: result.RunningBalance}, nil
}

func (s *Server) ClosePettyCashPeriod(ctx context.Context, req *pb.ClosePettyCashPeriodRequest) (*pb.ClosePettyCashPeriodResponse, error) {
	const op = "grpc_api.Server.ClosePettyCashPeriod"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ClosePettyCashPeriod(ctx, &service.ClosePettyCashPeriodParams{
		PeriodEnd: req.GetPeriodEnd(),
		Notes:     req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ClosePettyCashPeriodResponse{
		ClosingEntryID: result.ClosingEntryID,
		ClosingBalance: result.ClosingBalance,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-029 Project costs
// ---------------------------------------------------------------------------

func (s *Server) GetProjectCosts(ctx context.Context, req *pb.GetProjectCostsRequest) (*pb.GetProjectCostsResponse, error) {
	const op = "grpc_api.Server.GetProjectCosts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureID()).Msg("")

	result, err := s.svc.GetProjectCosts(ctx, &service.GetProjectCostsParams{DepartureID: req.GetDepartureID()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	lines := make([]*pb.CostLineItem, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, &pb.CostLineItem{
			Category:    l.Category,
			Description: l.Description,
			Amount:      l.Amount,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetProjectCostsResponse{
		DepartureID: result.DepartureID,
		TotalCost:   result.TotalCost,
		Lines:       lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-030 Departure P&L
// ---------------------------------------------------------------------------

func (s *Server) GetDeparturePL(ctx context.Context, req *pb.GetDeparturePLRequest) (*pb.GetDeparturePLResponse, error) {
	const op = "grpc_api.Server.GetDeparturePL"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureID()).Msg("")

	result, err := s.svc.GetDeparturePL(ctx, &service.GetDeparturePLParams{DepartureID: req.GetDepartureID()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetDeparturePLResponse{
		DepartureID:   result.DepartureID,
		Revenue:       result.Revenue,
		Costs:         result.Costs,
		GrossProfit:   result.GrossProfit,
		BudgetRevenue: result.BudgetRevenue,
		BudgetCosts:   result.BudgetCosts,
		Variance:      result.Variance,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-031 Budget vs actual
// ---------------------------------------------------------------------------

func (s *Server) GetBudgetVsActual(ctx context.Context, req *pb.GetBudgetVsActualRequest) (*pb.GetBudgetVsActualResponse, error) {
	const op = "grpc_api.Server.GetBudgetVsActual"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetBudgetVsActual(ctx, &service.GetBudgetVsActualParams{
		StartDate:   req.GetStartDate(),
		EndDate:     req.GetEndDate(),
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	lines := make([]*pb.BudgetLine, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, &pb.BudgetLine{
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Budgeted:    l.Budgeted,
			Actual:      l.Actual,
			Variance:    l.Variance,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetBudgetVsActualResponse{
		StartDate: result.StartDate,
		EndDate:   result.EndDate,
		Lines:     lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-032 Auto journal
// ---------------------------------------------------------------------------

func (s *Server) TriggerAutoJournal(ctx context.Context, req *pb.TriggerAutoJournalRequest) (*pb.TriggerAutoJournalResponse, error) {
	const op = "grpc_api.Server.TriggerAutoJournal"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("event_kind", req.GetEventKind()).Msg("")

	result, err := s.svc.TriggerAutoJournal(ctx, &service.TriggerAutoJournalParams{
		EventKind: req.GetEventKind(),
		SourceID:  req.GetSourceID(),
		Amount:    req.GetAmount(),
		Notes:     req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.TriggerAutoJournalResponse{JournalID: result.JournalID, Skipped: result.Skipped}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-033 Revenue recognition policy
// ---------------------------------------------------------------------------

func (s *Server) GetRevenueRecognitionPolicy(ctx context.Context, _ *pb.GetRevenueRecognitionPolicyRequest) (*pb.GetRevenueRecognitionPolicyResponse, error) {
	const op = "grpc_api.Server.GetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetRevenueRecognitionPolicy(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetRevenueRecognitionPolicyResponse{
		Policy: &pb.RevenueRecognitionPolicy{
			TriggerStatus:      result.TriggerStatus,
			DeferralAccount:    result.DeferralAccount,
			RecognitionAccount: result.RecognitionAccount,
		},
	}, nil
}

func (s *Server) SetRevenueRecognitionPolicy(ctx context.Context, req *pb.SetRevenueRecognitionPolicyRequest) (*pb.SetRevenueRecognitionPolicyResponse, error) {
	const op = "grpc_api.Server.SetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	err := s.svc.SetRevenueRecognitionPolicy(ctx, &service.SetRevenueRecognitionPolicyParams{
		TriggerStatus:      req.GetTriggerStatus(),
		DeferralAccount:    req.GetDeferralAccount(),
		RecognitionAccount: req.GetRecognitionAccount(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.SetRevenueRecognitionPolicyResponse{Updated: true}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-034 Exchange rates
// ---------------------------------------------------------------------------

func (s *Server) SetExchangeRate(ctx context.Context, req *pb.SetExchangeRateRequest) (*pb.SetExchangeRateResponse, error) {
	const op = "grpc_api.Server.SetExchangeRate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SetExchangeRate(ctx, &service.SetExchangeRateParams{
		FromCurrency:  req.GetFromCurrency(),
		ToCurrency:    req.GetToCurrency(),
		Rate:          req.GetRate(),
		EffectiveDate: req.GetEffectiveDate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.SetExchangeRateResponse{RateID: result.RateID}, nil
}

func (s *Server) GetExchangeRate(ctx context.Context, req *pb.GetExchangeRateRequest) (*pb.GetExchangeRateResponse, error) {
	const op = "grpc_api.Server.GetExchangeRate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetExchangeRate(ctx, &service.GetExchangeRateParamsService{
		FromCurrency: req.GetFromCurrency(),
		ToCurrency:   req.GetToCurrency(),
		AsOf:         req.GetAsOf(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrNotFound) {
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetExchangeRateResponse{
		RateID:        result.RateID,
		Rate:          result.Rate,
		EffectiveDate: result.EffectiveDate,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-035 Fixed assets
// ---------------------------------------------------------------------------

func (s *Server) CreateFixedAsset(ctx context.Context, req *pb.CreateFixedAssetRequest) (*pb.CreateFixedAssetResponse, error) {
	const op = "grpc_api.Server.CreateFixedAsset"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("name", req.GetName()).Msg("")

	result, err := s.svc.CreateFixedAsset(ctx, &service.CreateFixedAssetParams{
		Name:             req.GetName(),
		Category:         req.GetCategory(),
		PurchaseDate:     req.GetPurchaseDate(),
		PurchaseCost:     req.GetPurchaseCost(),
		UsefulLifeMonths: req.GetUsefulLifeMonths(),
		ResidualValue:    req.GetResidualValue(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.CreateFixedAssetResponse{AssetID: result.AssetID}, nil
}

func (s *Server) ListFixedAssets(ctx context.Context, req *pb.ListFixedAssetsRequest) (*pb.ListFixedAssetsResponse, error) {
	const op = "grpc_api.Server.ListFixedAssets"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListFixedAssets(ctx, req.GetCategory())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	assets := make([]*pb.FixedAssetRow, 0, len(result.Assets))
	for _, a := range result.Assets {
		assets = append(assets, &pb.FixedAssetRow{
			AssetID:                 a.AssetID,
			Name:                    a.Name,
			Category:                a.Category,
			PurchaseDate:            a.PurchaseDate,
			PurchaseCost:            a.PurchaseCost,
			AccumulatedDepreciation: a.AccumulatedDepreciation,
			BookValue:               a.BookValue,
			UsefulLifeMonths:        a.UsefulLifeMonths,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ListFixedAssetsResponse{Assets: assets}, nil
}

func (s *Server) RunDepreciation(ctx context.Context, req *pb.RunDepreciationRequest) (*pb.RunDepreciationResponse, error) {
	const op = "grpc_api.Server.RunDepreciation"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RunDepreciation(ctx, &service.RunDepreciationParams{AsOf: req.GetAsOf()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RunDepreciationResponse{
		AssetsProcessed:   result.AssetsProcessed,
		TotalDepreciation: result.TotalDepreciation,
		JournalID:         result.JournalID,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-036 Tax
// ---------------------------------------------------------------------------

func (s *Server) CalculateTax(ctx context.Context, req *pb.CalculateTaxRequest) (*pb.CalculateTaxResponse, error) {
	const op = "grpc_api.Server.CalculateTax"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("tax_type", req.GetTaxType()).Msg("")

	result, err := s.svc.CalculateTax(ctx, &service.CalculateTaxParams{
		BaseAmount: req.GetBaseAmount(),
		TaxType:    req.GetTaxType(),
		Rate:       req.GetRate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.CalculateTaxResponse{
		BaseAmount: result.BaseAmount,
		TaxAmount:  result.TaxAmount,
		TaxType:    result.TaxType,
		Rate:       result.Rate,
	}, nil
}

func (s *Server) GetTaxReport(ctx context.Context, req *pb.GetTaxReportRequest) (*pb.GetTaxReportResponse, error) {
	const op = "grpc_api.Server.GetTaxReport"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetTaxReport(ctx, &service.GetTaxReportParams{
		StartDate: req.GetStartDate(),
		EndDate:   req.GetEndDate(),
		TaxType:   req.GetTaxType(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	rows := make([]*pb.TaxReportRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.TaxReportRow{
			Date:        r.Date,
			Description: r.Description,
			BaseAmount:  r.BaseAmount,
			TaxAmount:   r.TaxAmount,
			TaxType:     r.TaxType,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetTaxReportResponse{
		StartDate: result.StartDate,
		EndDate:   result.EndDate,
		TotalTax:  result.TotalTax,
		Rows:      rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-037 Agent commission
// ---------------------------------------------------------------------------

func (s *Server) CreateCommissionPayout(ctx context.Context, req *pb.CreateCommissionPayoutRequest) (*pb.CreateCommissionPayoutResponse, error) {
	const op = "grpc_api.Server.CreateCommissionPayout"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("agent_id", req.GetAgentID()).Msg("")

	result, err := s.svc.CreateCommissionPayout(ctx, &service.CreateCommissionPayoutParams{
		AgentID:     req.GetAgentID(),
		DepartureID: req.GetDepartureID(),
		Amount:      req.GetAmount(),
		BasisAmount: req.GetBasisAmount(),
		RatePercent: req.GetRatePercent(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.CreateCommissionPayoutResponse{PayoutID: result.PayoutID, Status: result.Status}, nil
}

func (s *Server) DecideCommissionPayout(ctx context.Context, req *pb.DecideCommissionPayoutRequest) (*pb.DecideCommissionPayoutResponse, error) {
	const op = "grpc_api.Server.DecideCommissionPayout"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("payout_id", req.GetPayoutID()).Msg("")

	result, err := s.svc.DecideCommissionPayout(ctx, &service.DecideCommissionPayoutParams{
		PayoutID: req.GetPayoutID(),
		Decision: req.GetDecision(),
		Notes:    req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		if errors.Is(err, service.ErrValidation) {
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
		}
		if errors.Is(err, service.ErrNotFound) {
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		}
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.DecideCommissionPayoutResponse{PayoutID: result.PayoutID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-038 Realtime summary
// ---------------------------------------------------------------------------

func (s *Server) GetRealtimeFinancialSummary(ctx context.Context, _ *pb.GetRealtimeFinancialSummaryRequest) (*pb.GetRealtimeFinancialSummaryResponse, error) {
	const op = "grpc_api.Server.GetRealtimeFinancialSummary"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetRealtimeFinancialSummary(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	accounts := make([]*pb.RealtimeSummaryAccount, 0, len(result.Accounts))
	for _, a := range result.Accounts {
		accounts = append(accounts, &pb.RealtimeSummaryAccount{
			AccountCode: a.AccountCode,
			AccountName: a.AccountName,
			Balance:     a.Balance,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetRealtimeFinancialSummaryResponse{
		AsOf:         result.AsOf,
		TotalRevenue: result.TotalRevenue,
		TotalExpense: result.TotalExpense,
		NetIncome:    result.NetIncome,
		CashBalance:  result.CashBalance,
		ARBalance:    result.ARBalance,
		APBalance:    result.APBalance,
		Accounts:     accounts,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-039 Cash flow dashboard
// ---------------------------------------------------------------------------

func (s *Server) GetCashFlowDashboard(ctx context.Context, req *pb.GetCashFlowDashboardRequest) (*pb.GetCashFlowDashboardResponse, error) {
	const op = "grpc_api.Server.GetCashFlowDashboard"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetCashFlowDashboard(ctx, &service.GetCashFlowDashboardParams{
		StartDate: req.GetStartDate(),
		EndDate:   req.GetEndDate(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	lines := make([]*pb.CashFlowLine, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, &pb.CashFlowLine{
			Date:           l.Date,
			Description:    l.Description,
			Amount:         l.Amount,
			RunningBalance: l.RunningBalance,
			Category:       l.Category,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetCashFlowDashboardResponse{
		StartDate:      result.StartDate,
		EndDate:        result.EndDate,
		OpeningBalance: result.OpeningBalance,
		ClosingBalance: result.ClosingBalance,
		OperatingNet:   result.OperatingNet,
		InvestingNet:   result.InvestingNet,
		FinancingNet:   result.FinancingNet,
		Lines:          lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-040 Aging alerts
// ---------------------------------------------------------------------------

func (s *Server) GetAgingAlerts(ctx context.Context, req *pb.GetAgingAlertsRequest) (*pb.GetAgingAlertsResponse, error) {
	const op = "grpc_api.Server.GetAgingAlerts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("ledger_type", req.GetLedgerType()).Msg("")

	result, err := s.svc.GetAgingAlerts(ctx, &service.GetAgingAlertsParams{
		LedgerType:    req.GetLedgerType(),
		ThresholdDays: req.GetThresholdDays(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	alerts := make([]*pb.AgingAlert, 0, len(result.Alerts))
	for _, a := range result.Alerts {
		alerts = append(alerts, &pb.AgingAlert{
			EntityID:    a.EntityID,
			EntityName:  a.EntityName,
			Amount:      a.Amount,
			DaysOverdue: a.DaysOverdue,
			LedgerType:  a.LedgerType,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetAgingAlertsResponse{Alerts: alerts, TotalOverdue: result.TotalOverdue}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-041 Audit trail
// ---------------------------------------------------------------------------

func (s *Server) SearchFinanceAuditLog(ctx context.Context, req *pb.SearchFinanceAuditLogRequest) (*pb.SearchFinanceAuditLogResponse, error) {
	const op = "grpc_api.Server.SearchFinanceAuditLog"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SearchFinanceAuditLog(ctx, &service.SearchFinanceAuditLogParams{
		UserID:     req.GetUserID(),
		Action:     req.GetAction(),
		EntityType: req.GetEntityType(),
		EntityID:   req.GetEntityID(),
		StartDate:  req.GetStartDate(),
		EndDate:    req.GetEndDate(),
		PageSize:   req.GetPageSize(),
		Cursor:     req.GetCursor(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	rows := make([]*pb.AuditLogRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.AuditLogRow{
			LogID:      r.LogID,
			UserID:     r.UserID,
			Action:     r.Action,
			EntityType: r.EntityType,
			EntityID:   r.EntityID,
			Diff:       r.Diff,
			CreatedAt:  r.CreatedAt,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.SearchFinanceAuditLogResponse{Rows: rows, NextCursor: result.NextCursor}, nil
}
