// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao/internal"
	"hotgo/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
)

// internalAdminRoleDao is internal type for wrapping internal DAO implements.
type internalAdminRoleDao = *internal.AdminRoleDao

// adminRoleDao is the data access object for table admin_role.
// You can define custom methods on it to extend its functionality as you wish.
type adminRoleDao struct {
	internalAdminRoleDao
}

var (
	// AdminRole is globally common accessible object for table admin_role operations.
	AdminRole = adminRoleDao{
		internal.NewAdminRoleDao(),
	}
)

// IsUniqueName 判断名称是否唯一
func (dao *adminRoleDao) IsUniqueName(ctx context.Context, id int64, name string) (bool, error) {
	var data *entity.AdminRole
	m := dao.Ctx(ctx).Where("name", name)

	if id > 0 {
		m = m.WhereNot("id", id)
	}

	if err := m.Scan(&data); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return false, err
	}

	if data == nil {
		return true, nil
	}

	return false, nil
}

// IsUniqueCode 判断编码是否唯一
func (dao *adminRoleDao) IsUniqueCode(ctx context.Context, id int64, code string) (bool, error) {
	var data *entity.AdminRole
	m := dao.Ctx(ctx).Where("key", code)

	if id > 0 {
		m = m.WhereNot("id", id)
	}

	if err := m.Scan(&data); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return false, err
	}

	if data == nil {
		return true, nil
	}

	return false, nil
}
