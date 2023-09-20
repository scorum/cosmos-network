package keeper_test

import (
	"testing"

	codec "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/google/uuid"
	"github.com/mikluke/co-pilot/handle"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func TestQueryServer_PlaneByName(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()

	id := uuid.New().String()

	s := keeper.NewMsgServerImpl(set.keeper)
	q := types.QueryServer(set.keeper)

	act, err := q.PlaneByName(ctx.Context, &types.QueryPlaneByNameRequest{Name: "name"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
	require.Nil(t, act)

	_, err = s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
		Id:         id,
		Supervisor: set.supervisor.String(),
		Owner:      addr.String(),
		Meta: &types.PlaneMeta{
			Name:       "name",
			Color:      "white",
			Experience: 100,
		},
	})
	require.NoError(t, err)

	act, err = q.PlaneByName(ctx.Context, &types.QueryPlaneByNameRequest{Name: "name"})
	require.NoError(t, err)
	require.Equal(t, types.NftClassID, act.Nft.ClassId)
	require.Equal(t, id, act.Nft.Id)
	require.Equal(t, handle.Must(codec.NewAnyWithValue(&types.PlaneMeta{
		Name:       "name",
		Color:      "white",
		Experience: 100,
	})).Value, act.Nft.Data.Value)
}
