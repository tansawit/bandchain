package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k Keeper) HasReport(ctx sdk.Context, rid types.RID, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.RawDataReportStoreKey(rid, val))
}

// SetDataReport saves the report to the storage without performing validation.
func (k Keeper) SetReport(ctx sdk.Context, rid types.RID, report types.Report) {
	key := types.RawDataReportStoreKey(rid, report.Validator)
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshalBinaryBare(report))
}

// AddReports performs sanity checks and adds a new batch from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k Keeper) AddReport(ctx sdk.Context, rid types.RID, report types.Report) error {
	req, err := k.GetRequest(ctx, rid)
	if err != nil {
		return err
	}

	// TODO: Keep statistics of validator reports, so we can jail inactive validators!

	if !ContainsVal(req.RequestedValidators, report.Validator) {
		return sdkerrors.Wrapf(
			types.ErrValidatorNotRequested, "reqID: %d, val: %s", rid, report.Validator.String())
	}
	if k.HasReport(ctx, rid, report.Validator) {
		return sdkerrors.Wrapf(
			types.ErrValidatorAlreadyReported, "reqID: %d, val: %s", rid, report.Validator.String())
	}
	if int64(len(report.RawDataReports)) != k.GetRawRequestCount(ctx, rid) {
		return types.ErrInvalidDataSourceCount
	}

	for _, rep := range report.RawDataReports {
		// Here we can safely assume that external IDs are unique, as this has already been
		// checked by ValidateBasic performed in baseapp's runTx function.
		if !k.HasRawRequest(ctx, rid, rep.ExternalID) {
			return sdkerrors.Wrapf(
				types.ErrRawRequestNotFound, "reqID: %d, extID: %d", rid, rep.ExternalID)
		}
		if err := k.EnsureLength(ctx, types.KeyMaxRawDataReportSize, len(rep.Data)); err != nil {
			return err
		}
	}

	k.SetReport(ctx, rid, report)
	return nil
}

// GetReportIterator returns the iterator for all reports of the given request ID.
func (k Keeper) GetReportIterator(ctx sdk.Context, rid types.RequestID) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataReportStoreKeyPrefix, rid)
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

// GetReportCount returns the number of reports for the given request ID.
func (k Keeper) GetReportCount(ctx sdk.Context, rid types.RID) (count int64) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

// GetReports returns all reports for the given request ID, or nil if there is none.
func (k Keeper) GetReports(ctx sdk.Context, rid types.RID) (res []types.Report) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)
		res = append(res, report)
	}
	return res
}
