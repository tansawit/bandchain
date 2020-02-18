package zoracle

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keep "github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func TestNewExecutionEnvironment(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	_, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.NotNil(t, err)

	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100,
	))

	_, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
}

func TestGetCurrentRequestID(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100,
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1), env.GetCurrentRequestID())
}

func TestGetRequestedValidatorCount(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100,
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(2), env.GetRequestedValidatorCount())
}

func TestGetReceivedValidatorCount(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100,
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(0), env.GetReceivedValidatorCount())

	keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("val1")))

	env, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1), env.GetReceivedValidatorCount())

}

func TestGetPrepareBlockTime(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 20, 1581589790, 100,
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1581589790), env.GetPrepareBlockTime())
}

func TestGetAggregateBlockTime(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100,
	))

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(0), env.GetAggregateBlockTime())

	// Add received validator
	err = keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("val1")))
	require.Nil(t, err)

	// After report is greater or equal SufficientValidatorCount, it will resolve in current block time.
	env, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1581589790), env.GetAggregateBlockTime())
}

func TestGetValidatorPubKey(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}
	validatorAddress1 := keep.SetupTestValidator(
		ctx,
		keeper,
		pubStr[0],
		10,
	)
	validatorAddress2 := keep.SetupTestValidator(
		ctx,
		keeper,
		pubStr[1],
		10,
	)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{validatorAddress1, validatorAddress2}, 1, 0, 0, 100,
	))

	env, errSDK := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, errSDK)

	addr1, err := env.GetValidatorAddress(0)
	require.Nil(t, err)
	require.Equal(t, validatorAddress1, sdk.ValAddress(addr1))

	addr2, err := env.GetValidatorAddress(1)
	require.Nil(t, err)
	require.Equal(t, validatorAddress2, sdk.ValAddress(addr2))

	_, err = env.GetValidatorAddress(2)
	require.NotNil(t, err)

	_, err = env.GetValidatorAddress(-1)
	require.NotNil(t, err)
}

func TestRequestExternalData(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetMaxCalldataSize(ctx, 20)
	keeper.SetMaxDataSourceCountPerRequest(ctx, 1)
	// Set Request
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100,
	))

	// Set Datasource
	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	envErr := env.RequestExternalData(1, 42, []byte("prepare32"))
	require.Nil(t, envErr)

	rawRequest, err := keeper.GetRawDataRequest(ctx, 1, 42)
	require.Nil(t, err)
	require.Equal(t, types.NewRawDataRequest(1, []byte("prepare32")), rawRequest)
}

func TestRequestExternalDataExceedMaxDataSourceCountPerRequest(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetMaxCalldataSize(ctx, 20)

	// Set Request
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100,
	))

	// Set Datasource
	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)

	// Set MaxDataSourceCountPerRequest to 5
	keeper.SetMaxDataSourceCountPerRequest(ctx, 5)
	envErr := env.RequestExternalData(1, 41, []byte("prepare32"))
	require.Nil(t, envErr)
	envErr = env.RequestExternalData(1, 42, []byte("prepare32"))
	require.Nil(t, envErr)
	envErr = env.RequestExternalData(1, 43, []byte("prepare32"))
	require.Nil(t, envErr)
	envErr = env.RequestExternalData(1, 44, []byte("prepare32"))
	require.Nil(t, envErr)
	envErr = env.RequestExternalData(1, 45, []byte("prepare32"))
	require.Nil(t, envErr)

	// Should get an error if trying to request more external data
	envErr = env.RequestExternalData(1, 46, []byte("prepare32"))
	require.NotNil(t, envErr)
}

func TestGetExternalData(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100,
	))

	keeper.SetRawDataReport(ctx, 1, 42, sdk.ValAddress([]byte("val1")), []byte("data42"))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)

	// Get report from reported validator
	report, envErr := env.GetExternalData(42, 0)
	require.Nil(t, envErr)
	require.Equal(t, []byte("data42"), report)

	// Get report from missing validator
	_, envErr = env.GetExternalData(42, 1)
	require.EqualError(t, envErr, "ERROR:\nCodespace: zoracle\nCode: 109\nMessage: \"report not found\"\n")

	// Get report from invalid validator index
	_, envErr = env.GetExternalData(42, 2)
	require.NotNil(t, envErr, "validator out of range")

	// Get report from invalid validator index
	_, envErr = env.GetExternalData(42, -2)
	require.NotNil(t, envErr, "validator out of range")
}