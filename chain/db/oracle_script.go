package db

import (
	"errors"
	"strconv"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createOracleScript(
	id int64,
	name, description string,
	owner sdk.AccAddress,
	blockTime time.Time,
	schema string,
	sourceCodeURL string,
) OracleScript {

	return OracleScript{
		ID:            id,
		Name:          &name,
		Description:   &description,
		Owner:         owner.String(),
		LastUpdated:   blockTime.UnixNano() / int64(time.Millisecond),
		Schema:        &schema,
		SourceCodeURL: &sourceCodeURL,
	}
}

func (b *BandDB) AddOracleScript(
	id int64,
	name, description string,
	owner sdk.AccAddress,
	code []byte,
	blockTime time.Time,
	blockHeight int64,
	txHash []byte,
	schema string,
	sourceCodeURL string,
) error {

	oracleScript := createOracleScript(
		id,
		name,
		description,
		owner,
		blockTime,
		schema,
		sourceCodeURL,
	)
	err := b.tx.Create(&oracleScript).Error
	if err != nil {
		return err
	}

	return b.tx.Create(&OracleScriptRevision{
		OracleScriptID: id,
		Name:           name,
		Timestamp:      blockTime.UnixNano() / int64(time.Millisecond),
		BlockHeight:    blockHeight,
		TxHash:         txHash,
	}).Error
}

func (b *BandDB) handleMsgCreateOracleScript(
	txHash []byte,
	msg oracle.MsgCreateOracleScript,
	events map[string][]string,
) error {
	ids := events[otypes.EventTypeCreateOracleScript+"."+otypes.AttributeKeyID]
	if len(ids) != 1 {
		return errors.New("handleMsgCreateOracleScript: cannot find oracle script id")
	}
	id, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return err
	}
	return b.AddOracleScript(
		id, msg.Name, msg.Description, msg.Owner, msg.Code,
		b.ctx.BlockTime(), b.ctx.BlockHeight(), txHash, msg.Schema, msg.SourceCodeURL,
	)
}

func (b *BandDB) handleMsgEditOracleScript(
	txHash []byte,
	msg oracle.MsgEditOracleScript,
) error {
	var oracleScript OracleScript
	err := b.tx.First(&oracleScript, int64(msg.OracleScriptID)).Error
	if err != nil {
		return err
	}

	newName := modify(*oracleScript.Name, msg.Name)
	newDescription := modify(*oracleScript.Description, msg.Description)
	newSchema := modify(*oracleScript.Schema, msg.Schema)
	newSourceCodeURL := modify(*oracleScript.SourceCodeURL, msg.SourceCodeURL)

	oracleScript = createOracleScript(
		int64(msg.OracleScriptID), newName, newDescription, msg.Owner,
		b.ctx.BlockTime(), newSchema, newSourceCodeURL,
	)

	err = b.tx.Save(&oracleScript).Error
	if err != nil {
		return err
	}

	return b.tx.Create(&OracleScriptRevision{
		OracleScriptID: int64(msg.OracleScriptID),
		Name:           newName,
		Timestamp:      b.ctx.BlockTime().UnixNano() / int64(time.Millisecond),
		BlockHeight:    b.ctx.BlockHeight(),
		TxHash:         txHash,
	}).Error
}
