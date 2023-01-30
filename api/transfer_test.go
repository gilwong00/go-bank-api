package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "go-bank-api/pkg/db/mock"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	account1 := generateRandomAccount()
	account2 := generateRandomAccount()
	account3 := generateRandomAccount()

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.EUR

	testCases := []struct {
		name string
		body gin.H

		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account2.ID,
				"amount":        amount,
				"currency":      util.USD,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferFundsParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// {
		// 	name: "UnauthorizedUser",
		// 	body: gin.H{
		// 		"fromAccountId": account1.ID,
		// 		"toAccountId":   account2.ID,
		// 		"amount":        amount,
		// 		"currency":      util.USD,
		// 	},

		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
		// 		store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "NoAuthorization",
		// 	body: gin.H{
		// 		"fromAccountId": account1.ID,
		// 		"toAccountId":   account2.ID,
		// 		"amount":        amount,
		// 		"currency":      util.USD,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Any()).Times(0)
		// 		store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "FromAccountNotFound",
		// 	body: gin.H{
		// 		"fromAccountId": account1.ID,
		// 		"toAccountId":   account2.ID,
		// 		"amount":        amount,
		// 		"currency":      util.USD,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
		// 		store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "ToAccountNotFound",
		// 	body: gin.H{
		// 		"fromAccountId": account1.ID,
		// 		"toAccountId":   account2.ID,
		// 		"amount":        amount,
		// 		"currency":      util.USD,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
		// 		store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
		// 		store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"fromAccountId": account3.ID,
				"toAccountId":   account2.ID,
				"amount":        amount,
				"currency":      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account3.ID,
				"amount":        amount,
				"currency":      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account2.ID,
				"amount":        amount,
				"currency":      "XYZ",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account2.ID,
				"amount":        -amount,
				"currency":      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account2.ID,
				"amount":        amount,
				"currency":      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferFundTxError",
			body: gin.H{
				"fromAccountId": account1.ID,
				"toAccountId":   account2.ID,
				"amount":        amount,
				"currency":      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferFundsTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferFundsResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/transfer"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
