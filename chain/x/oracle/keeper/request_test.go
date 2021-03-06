package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHasRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a request ID 42 without setting it.
	require.False(t, k.HasRequest(ctx, 42))
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, "", nil, nil))
	require.True(t, k.HasRequest(ctx, 42))
}

func TestDeleteRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, "", nil, nil))
	require.True(t, k.HasRequest(ctx, 42))
	// After we delete it, we should not find it anymore.
	k.DeleteRequest(ctx, 42)
	require.False(t, k.HasRequest(ctx, 42))
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
}

func TestSetterGetterRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// Getting a non-existent request should return error.
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
	// Creates some basic requests.
	req1 := types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, "", nil, nil)
	req2 := types.NewRequest(2, BasicCalldata, nil, 1, 1, 1, "", nil, nil)
	// Sets id 42 with request 1 and id 42 with request 2.
	k.SetRequest(ctx, 42, req1)
	k.SetRequest(ctx, 43, req2)
	// Checks that Get and MustGet perform correctly.
	req1Res, err := k.GetRequest(ctx, 42)
	require.Nil(t, err)
	require.Equal(t, req1, req1Res)
	require.Equal(t, req1, k.MustGetRequest(ctx, 42))
	req2Res, err := k.GetRequest(ctx, 43)
	require.Nil(t, err)
	require.Equal(t, req2, req2Res)
	require.Equal(t, req2, k.MustGetRequest(ctx, 43))
	// Replaces id 42 with another request.
	k.SetRequest(ctx, 42, req2)
	require.NotEqual(t, req1, k.MustGetRequest(ctx, 42))
	require.Equal(t, req2, k.MustGetRequest(ctx, 42))
}

func TestSetterGettterPendingResolveList(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// After we set something, we should get that thing back.
	k.SetPendingResolveList(ctx, []types.RequestID{5, 6, 7, 8})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{5, 6, 7, 8})
	// Let's also try setting it back to empty list.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// Nil should also works.
	k.SetPendingResolveList(ctx, nil)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
}

func TestAddDataSourceBasic(t *testing.T) {
	_, ctx, k := createTestInput()
	// We start by setting an oracle request available at ID 42.
	k.SetOracleScript(ctx, 42, types.NewOracleScript(
		Owner.Address, BasicName, BasicDesc, BasicFilename, BasicSchema, BasicSourceCodeURL,
	))
	// Adding the first request should return ID 1.
	id := k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, 1, "", nil, nil))
	require.Equal(t, id, types.RequestID(1))
	// Adding another request should return ID 2.
	id = k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, 1, "", nil, nil))
	require.Equal(t, id, types.RequestID(2))
}

func TestAddPendingResolveList(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// Everytime we append a new request ID, it should show up.
	k.AddPendingRequest(ctx, 42)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{42})
	k.AddPendingRequest(ctx, 43)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{42, 43})
}

func TestGetRandomValidatorsSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	hash := sha256.Sum256([]byte("Hello"))
	ctx = ctx.WithBlockHeader(abci.Header{
		LastBlockId: abci.BlockID{Hash: hash[:32]},
	})
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	vals, err := k.GetRandomValidators(ctx, 3, int64(1))
	expect := []sdk.ValAddress{Validator1.ValAddress, Validator3.ValAddress, Validator2.ValAddress}

	require.NoError(t, err)
	require.Equal(t, vals, expect)

	hash = sha256.Sum256([]byte("Ni Hao"))
	ctx = ctx.WithBlockHeader(abci.Header{
		LastBlockId: abci.BlockID{Hash: hash[:32]},
	})
	vals, err = k.GetRandomValidators(ctx, 3, int64(2))
	expect = []sdk.ValAddress{Validator1.ValAddress, Validator3.ValAddress, Validator2.ValAddress}
	require.NoError(t, err)
	require.Equal(t, vals, expect)

	vals, err = k.GetRandomValidators(ctx, 1, int64(2))
	expect = []sdk.ValAddress{Validator1.ValAddress}
	require.NoError(t, err)
	require.Equal(t, vals, expect)

}

func TestGetRandomValidatorsTooBigSize(t *testing.T) {
	_, ctx, k := createTestInput()
	_, err := k.GetRandomValidators(ctx, 9999, int64(1))
	require.Error(t, err)
}
