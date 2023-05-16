package storage

import (
	"testing"

	"github.com/immune-gmbh/AttestationFailureAnalysisService/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestCompileWhereConds(t *testing.T) {
	{
		whereConds, whereArgs := compileFirmwareImageWhereConds(FindFirmwareFilter{})
		require.Empty(t, whereConds)
		require.Nil(t, whereArgs)
	}
	{
		whereConds, whereArgs := compileFirmwareImageWhereConds(FindFirmwareFilter{
			ImageID: &types.ImageID{1, 2, 3},
		})
		require.Equal(t, "`image_id` = ?", whereConds)
		require.Equal(t, []interface{}{types.ImageID{1, 2, 3}}, whereArgs)
	}
	{
		whereConds, whereArgs := compileFirmwareImageWhereConds(FindFirmwareFilter{
			ImageID:  &types.ImageID{1, 2, 3},
			Filename: &[]string{"unit-test"}[0],
		})
		require.Equal(t, "`image_id` = ? AND `filename` = ?", whereConds)
		require.Equal(t, []interface{}{types.ImageID{1, 2, 3}, "unit-test"}, whereArgs)
	}
}