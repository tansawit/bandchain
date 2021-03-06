package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	idTag       = "idTag"
	dataHashTag = "dataHashTag"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/params", storeName), getParamsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/counts", storeName), getCountsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data/{%s}", storeName, dataHashTag), getDataByHashHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data_sources/{%s}", storeName, idTag), getDataSourceByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/oracle_scripts/{%s}", storeName, idTag), getOracleScriptByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/requests/{%s}", storeName, idTag), getRequestByIDHandler(cliCtx, storeName)).Methods("GET")
}
