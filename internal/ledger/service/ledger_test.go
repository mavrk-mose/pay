package service_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/ledger/models"
	"github.com/mavrk-mose/pay/internal/ledger/mocks"
	"github.com/mavrk-mose/pay/internal/payment"

	"github.com/stretchr/testify/assert"
)

func TestPaymentService_MakePayment(t *testing.T) {
	t.Parallel() 

	testCases := []struct {
		name        string
		txn         models.Transaction
		mockSetup   func(mock *mocks.LedgerService, ctx *gin.Context, txn models.Transaction)
		expectedErr error
	}{
		{
			name: "Successful transaction",
			txn: models.Transaction{
				ID:       "txn-001",
				Amount:   100.0,
				Currency: "USD",
				Status:   models.StatusPending,
			},
			mockSetup: func(mock *mocks.LedgerService, ctx *gin.Context, txn models.Transaction) {
				mock.On("RecordTransaction", ctx, txn).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Ledger service error",
			txn: models.Transaction{
				ID:       "txn-002",
				Amount:   250.0,
				Currency: "KES",
				Status:   models.StatusPending,
			},
			mockSetup: func(mock *mocks.LedgerService, ctx *gin.Context, txn models.Transaction) {
				mock.On("RecordTransaction", ctx, txn).Return(errors.New("ledger service failure"))
			},
			expectedErr: errors.New("ledger service failure"),
		},
		{
			name: "Duplicate transaction error",
			txn: models.Transaction{
				ID:       "txn-003",
				Amount:   50.0,
				Currency: "EUR",
				Status:   models.StatusPending,
			},
			mockSetup: func(mock *mocks.LedgerService, ctx *gin.Context, txn models.Transaction) {
				mock.On("RecordTransaction", ctx, txn).Return(errors.New("duplicate transaction"))
			},
			expectedErr: errors.New("duplicate transaction"),
		},
	}

	for _, tc := range testCases {
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() 

			ctx := &gin.Context{}
			mockLedger := mocks.NewLedgerService(t)

			tc.mockSetup(mockLedger, ctx, tc.txn)

			paymentService := payment.NewPaymentService(mockLedger)
			err := paymentService.MakePayment(ctx, tc.txn)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockLedger.AssertExpectations(t)
		})
	}
}
