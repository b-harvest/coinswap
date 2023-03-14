package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bharvest/coinswap/modules/coinswap/types"
)

func (s *TestSuite) TestGRPCParams() {
	resp, err := s.queryClient.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(s.keeper.GetParams(s.ctx), resp.Params)
}
