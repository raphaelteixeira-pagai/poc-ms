package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/raphaelteixeira-pagai/poc-ms/db/mock"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/entities"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/services"
	"github.com/raphaelteixeira-pagai/poc-ms/pkg/utils/random"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	app   *gin.Engine
	store *mockdb.MockIWalletRepository
	srv   services.IWalletService
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	app = gin.Default()

	code := m.Run()

	os.Exit(code)
}

func setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store = mockdb.NewMockIWalletRepository(ctrl)
	srv = services.NewWalletService(store)

	router := app.Group("/")
	RegisterWalletRoutes(srv, router)
}

func TestWallethandlers_CreateWallet(t *testing.T) {
	setup(t)
	payload := entities.Wallet{
		ID:      int64(rand.Intn(100)),
		Balance: 0,
		Owner:   random.Owner(),
	}
	empty := entities.Wallet{
		Balance: 0,
		Owner:   "",
	}
	negativeBalance := entities.Wallet{
		Owner:   random.Owner(),
		Balance: -10,
	}

	testCases := []struct {
		name          string
		wallet        entities.Wallet
		buildStubd    func(store *mockdb.MockIWalletRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			wallet: payload,
			buildStubd: func(store *mockdb.MockIWalletRepository) {
				store.EXPECT().
					Get(gomock.Any(), gomock.Eq(payload.Owner)).
					Times(1).
					Return(entities.Wallet{}, dbr.ErrNotFound)

				store.EXPECT().
					Create(gomock.Any(), gomock.Eq(payload)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:   "Duplicated",
			wallet: payload,
			buildStubd: func(store *mockdb.MockIWalletRepository) {
				store.EXPECT().
					Get(gomock.Any(), gomock.Eq(payload.Owner)).
					Times(1).
					Return(payload, nil)

				store.EXPECT().
					Create(gomock.Any(), gomock.Eq(payload)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Uncompleted",
			wallet: empty,
			buildStubd: func(store *mockdb.MockIWalletRepository) {
				store.EXPECT().
					Get(gomock.Any(), gomock.Eq(payload.Owner)).
					Times(1).
					Return(empty, dbr.ErrNotFound)

				store.EXPECT().
					Create(gomock.Any(), gomock.Eq(empty)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Negative Balance",
			wallet: negativeBalance,
			buildStubd: func(store *mockdb.MockIWalletRepository) {
				store.EXPECT().
					Get(gomock.Any(), gomock.Eq(negativeBalance.Owner)).
					Times(1).
					Return(empty, dbr.ErrNotFound)

				store.EXPECT().
					Create(gomock.Any(), gomock.Eq(negativeBalance)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubd(store)

			data, _ := json.Marshal(tc.wallet)
			url := "/wallet/"

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			w := httptest.NewRecorder()

			app.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}
