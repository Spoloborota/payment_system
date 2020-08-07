package restapi

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/common"
	"payment-system/internal/controller"
	"payment-system/internal/restapi/handlers"
	"strings"

	"github.com/gorilla/mux"
)

const baseURL = "/api/v1"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(logger *zap.Logger, cntrlr *controller.Controller) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	routes := FillRoutes(logger, cntrlr)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(logger, handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	fs := http.FileServer(http.Dir("./web/swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func FillRoutes(logger *zap.Logger, cntrlr *controller.Controller) Routes {
	createWalletWrapper := handlers.CreateWalletWrapper{Logger: logger.With(zap.String(common.Handler, "CreateWalletHandler")), Cntrlr: cntrlr}
	walletWrapper := handlers.WalletWrapper{Logger: logger.With(zap.String(common.Handler, "WalletHandler")), Cntrlr: cntrlr}
	topUpWrapper := handlers.TopUpWrapper{Logger: logger.With(zap.String(common.Handler, "TopUpHandler")), Cntrlr: cntrlr}
	transactionWrapper := handlers.CreateTransactionWrapper{Logger: logger.With(zap.String(common.Handler, "TransactionHandler")), Cntrlr: cntrlr}
	transactionStatusWrapper := handlers.TransactionStatusWrapper{Logger: logger.With(zap.String(common.Handler, "TransactionStatusHandler")), Cntrlr: cntrlr}
	transactionsReportWrapper := handlers.TransactionsReportWrapper{Logger: logger.With(zap.String(common.Handler, "TransactionsReportHandler")), Cntrlr: cntrlr}

	return Routes{
		Route{
			"Index",
			"GET",
			baseURL + "/",
			Index,
		},

		Route{
			"TopUp",
			strings.ToUpper("Post"),
			baseURL + "/topup",
			topUpWrapper.TopUp,
		},

		Route{
			"CreateTransaction",
			strings.ToUpper("Post"),
			baseURL + "/transactions",
			transactionWrapper.CreateTransaction,
		},

		Route{
			"TransactionStatus",
			strings.ToUpper("Get"),
			baseURL + "/transactions/{id}",
			transactionStatusWrapper.TransactionStatus,
		},

		Route{
			"TransactionsReport",
			strings.ToUpper("Get"),
			baseURL + "/transactions",
			transactionsReportWrapper.TransactionsReport,
		},

		Route{
			"CreateWallet",
			strings.ToUpper("Post"),
			baseURL + "/wallets",
			createWalletWrapper.CreateWallet,
		},

		Route{
			"Wallet",
			strings.ToUpper("Get"),
			baseURL + "/wallets/{id}",
			walletWrapper.Wallet,
		},
	}
}
