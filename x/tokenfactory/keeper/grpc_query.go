package keeper

import (
	"github.com/wfblockchain/noblechain/v5/x/tokenfactory/types"
)

var _ types.QueryServer = Keeper{}
