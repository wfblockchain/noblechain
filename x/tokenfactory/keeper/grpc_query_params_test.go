package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/wfblockchain/noblechain/v5/testutil/keeper"
	"github.com/wfblockchain/noblechain/v5/x/tokenfactory/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.TokenfactoryKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
