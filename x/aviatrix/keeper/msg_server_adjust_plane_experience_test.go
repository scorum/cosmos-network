package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func TestMsgServer_AdjustPlaneExperience(t *testing.T) {
	supervisorAddr := sample.AccAddress()
	addr := sample.AccAddress()
	id := uuid.New().String()

	testCases := []struct {
		name string

		owner string
		msg   *types.MsgAdjustPlaneExperience

		isError          bool
		resultExperience uint64
	}{
		{
			name:  "not allowed",
			owner: sample.AccAddress().String(),
			msg: &types.MsgAdjustPlaneExperience{
				Supervisor: sample.AccAddress().String(),
				Id:         id,
				Amount:     1,
			},
			isError: true,
		},
		{
			name:  "not supervisor but owner",
			owner: addr.String(),
			msg: &types.MsgAdjustPlaneExperience{
				Supervisor: addr.String(),
				Id:         id,
				Amount:     1,
			},

			isError: true,
		},
		{
			name:  "add experience by supervisor",
			owner: supervisorAddr.String(),
			msg: &types.MsgAdjustPlaneExperience{
				Supervisor: supervisorAddr.String(),
				Id:         id,
				Amount:     1,
			},
			resultExperience: 11,
		},
		{
			name:  "reduce experience by supervisor",
			owner: supervisorAddr.String(),
			msg: &types.MsgAdjustPlaneExperience{
				Supervisor: supervisorAddr.String(),
				Id:         id,
				Amount:     -1,
			},
			resultExperience: 9,
		},
		{
			name:  "potential experience overflow",
			owner: supervisorAddr.String(),
			msg: &types.MsgAdjustPlaneExperience{
				Supervisor: supervisorAddr.String(),
				Id:         id,
				Amount:     -1000,
			},
			resultExperience: 0,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			set, ctx := setupKeeper(t)
			set.scorumKeeper.SetParams(ctx.Context, scorumtypes.Params{
				Supervisors:            []string{supervisorAddr.String()},
				GasLimit:               math.NewInt(1000),
				GasUnconditionedAmount: math.NewInt(500),
				GasAdjustCoefficient:   math.LegacyNewDec(1),
			})

			s := keeper.NewMsgServerImpl(set.keeper)
			_, err := s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
				Id:         id,
				Supervisor: supervisorAddr.String(),
				Owner:      addr.String(),
				Meta: &types.PlaneMeta{
					Experience: 10,
				},
			})
			require.NoError(t, err)

			_, err = s.AdjustPlaneExperience(ctx.Context, tc.msg)

			if tc.isError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			act, err := set.keeper.GetPlane(ctx.Context, id)
			require.NoError(t, err)
			require.Equal(t, &types.Plane{
				Id:    id,
				Owner: addr.String(),
				Meta: &types.PlaneMeta{
					Experience: tc.resultExperience,
				},
			}, act)
		})
	}
}
