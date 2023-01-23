package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "go-bank-api/pkg/db/mock"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
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

	store := mockdb.NewMockStore(ctrl)
	// build stubs

	store.EXPECT().
		GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/api/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func generateRandomAccount() db.Account {
	return db.Account{
		ID:       util.GetRandomInt(1, 1000),
		Owner:    util.GetRandomOwner(),
		Balance:  util.GetRandomBalance(),
		Currency: util.GetCurrencyType(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
