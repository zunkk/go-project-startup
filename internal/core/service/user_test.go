package service

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/zunkk/go-project-startup/internal/core/dao"
	"github.com/zunkk/go-project-startup/internal/core/model"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-sidecar/db/memory"
)

func PrepareDB(t *testing.T) (*base.CustomSidecar, *dao.SQLConnector) {
	sidecar := base.NewMockCustomSidecar(t)
	memoryDB, err := memory.OpenSQLDB()
	require.Nil(t, err)
	sqlConnector, err := dao.NewSQLConnectorWithDB(sidecar, memoryDB)
	require.Nil(t, err)
	err = sqlConnector.Start()
	require.Nil(t, err)
	return sidecar, sqlConnector
}

func TestUserService_QueryByID(t *testing.T) {
	sidecar, sqlConnector := PrepareDB(t)

	userSrv, err := NewUserService(sidecar, sqlConnector)
	require.Nil(t, err)

	ctx := sidecar.BackgroundContext()
	userID := int64(1)
	now := time.Now()
	id, err := model.Users.Insert(&model.UserSetter{
		ID:         lo.ToPtr(userID),
		CreateTime: lo.ToPtr(now),
		UpdateTime: lo.ToPtr(now),
		DeleteTime: lo.ToPtr(time.Time{}),
		DelState:   lo.ToPtr(int64(0)),
		Version:    lo.ToPtr(int64(0)),
		Nickname:   lo.ToPtr("test"),
		Info:       lo.ToPtr("test"),
		Role:       lo.ToPtr("test"),
	}).Exec(ctx.Ctx, sqlConnector.DB)
	require.Nil(t, err)
	require.Equal(t, userID, id)

	user, err := userSrv.QueryByID(ctx.Ctx, userID)
	require.Nil(t, err)
	require.Equal(t, userID, user.ID)
	require.Equal(t, now.Unix(), user.CreateTime.Unix())
	require.Equal(t, now.Unix(), user.UpdateTime.Unix())
	require.Equal(t, time.Time{}.Unix(), user.DeleteTime.Unix())
	require.Equal(t, int64(0), user.DelState)
	require.Equal(t, int64(0), user.Version)
	require.Equal(t, "test", user.Nickname)
	require.Equal(t, "test", user.Info)
	require.Equal(t, "test", user.Role)
}
