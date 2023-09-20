package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func TestMsgServer_CreatePlane(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()

	s := keeper.NewMsgServerImpl(set.keeper)

	t.Run("not supervisor", func(t *testing.T) {
		id := uuid.New()

		_, err := s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
			Id:         id.String(),
			Supervisor: addr.String(),
			Owner:      addr.String(),
			Meta: &types.PlaneMeta{
				Name:       "name",
				Color:      "white",
				Experience: 100,
			},
		})

		require.Error(t, err)
		require.True(t, errorsmod.IsOf(err, sdkerrors.ErrUnauthorized))
	})

	t.Run("supervisor_self", func(t *testing.T) {
		id := uuid.New()

		_, err := s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
			Id:         id.String(),
			Supervisor: set.supervisor.String(),
			Owner:      set.supervisor.String(),
			Meta: &types.PlaneMeta{
				Name:       "name",
				Color:      "white",
				Experience: 100,
			},
		})
		require.NoError(t, err)

		act, err := set.keeper.GetPlane(ctx.Context, id.String())
		require.NoError(t, err)

		require.Equal(t, &types.Plane{
			Id:    id.String(),
			Owner: set.supervisor.String(),
			Meta: &types.PlaneMeta{
				Name:       "name",
				Color:      "white",
				Experience: 100,
			},
		}, act)

		require.Equal(t, set.supervisor, set.nftKeeper.GetOwner(ctx.Context, types.NftClassID, id.String()))
	})

	t.Run("supervisor_to_addr", func(t *testing.T) {
		id := uuid.New()

		_, err := s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
			Id:         id.String(),
			Supervisor: set.supervisor.String(),
			Owner:      addr.String(),
			Meta: &types.PlaneMeta{
				Name:       "name2",
				Color:      "white",
				Experience: 100,
			},
		})
		require.NoError(t, err)

		act, err := set.keeper.GetPlane(ctx.Context, id.String())
		require.NoError(t, err)

		require.Equal(t, &types.Plane{
			Id:    id.String(),
			Owner: addr.String(),
			Meta: &types.PlaneMeta{
				Name:       "name2",
				Color:      "white",
				Experience: 100,
			},
		}, act)

		require.Equal(t, addr, set.nftKeeper.GetOwner(ctx.Context, types.NftClassID, id.String()))
	})
}
