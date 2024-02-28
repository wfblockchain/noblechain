package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/wfblockchain/noblechain/v5/x/tariff/types"
)

type (
	Keeper struct {
		paramstore       paramtypes.Subspace
		authKeeper       types.AccountKeeper
		bankKeeper       types.BankKeeper
		feeCollectorName string // name of the FeeCollector ModuleAccount
		ics4Wrapper      porttypes.ICS4Wrapper
	}
)

// NewKeeper constructs a new fee collector keeper.
func NewKeeper(
	ps paramtypes.Subspace,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	feeCollectorName string,
	ics4Wrapper porttypes.ICS4Wrapper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		paramstore:       ps,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		ics4Wrapper:      ics4Wrapper,
	}
}

// WriteAcknowledgement implements the ICS4Wrapper interface.
func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
	ack exported.Acknowledgement,
) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (k Keeper) GetAppVersion(
	ctx sdk.Context,
	portID string,
	channelID string,
) (string, bool) {
	return k.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
