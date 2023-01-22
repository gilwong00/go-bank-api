package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockDB "go-bank-api/pkg/db/mock"
	"go-bank-api/pkg/util"
	"go-bank-api/sqlc"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountById(t *testing.T) {
	account := generateRandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockDB.NewMockStore(ctrl)
	// build stubs

	// expect the GetAccountById method to be called with any context but this account id
	store.EXPECT().
		GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func generateRandomAccount() sqlc.Account {
	return sqlc.Account{
		ID:       util.GetRandomInt(1, 1000),
		Owner:    util.GetRandomOwner(),
		Balance:  util.GetRandomBalance(),
		Currency: util.GetCurrencyType(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account sqlc.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount sqlc.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []sqlc.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []sqlc.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
