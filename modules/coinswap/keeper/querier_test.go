package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bharvest/coinswap/modules/coinswap/keeper"
	"github.com/bharvest/coinswap/modules/coinswap/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	legacyAmino := suite.app.LegacyAmino()

	k, ctx := suite.app.CoinswapKeeper, suite.ctx

	params := k.GetParams(ctx)
	params.WhitelistedDenoms = []string{denomBTC}
	k.SetParams(ctx, params)

	querier := keeper.NewQuerier(suite.app.CoinswapKeeper, legacyAmino)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	btcAmt, _ := sdk.NewIntFromString("100")
	standardAmt, _ := sdk.NewIntFromString("1000000000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender1.String())
	lptCoin, _ := k.AddLiquidity(ctx, msg)

	// test queryLiquidity

	bz, errRes := legacyAmino.MarshalJSON(types.QueryPoolParams{LptDenom: lptCoin.Denom})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryPool}, req)
	suite.NoError(err)

	var response types.QueryLiquidityPoolResponse
	errRes = suite.app.LegacyAmino().UnmarshalJSON(res, &response)
	suite.NoError(errRes)
	standard := sdk.NewCoin(denomStandard, standardAmt)
	token := sdk.NewCoin(denomBTC, btcAmt)
	liquidity := sdk.NewCoin(lptCoin.Denom, standardAmt)
	suite.Equal(standard, response.Pool.Standard)
	suite.Equal(token, response.Pool.Token)
	suite.Equal(liquidity, response.Pool.Lpt)
	suite.Equal(k.GetParams(ctx).Fee.String(), response.Pool.Fee)
}
