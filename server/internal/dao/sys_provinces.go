// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"hotgo/internal/dao/internal"
)

// internalSysProvincesDao is internal type for wrapping internal DAO implements.
type internalSysProvincesDao = *internal.SysProvincesDao

// sysProvincesDao is the data access object for table sys_provinces.
// You can define custom methods on it to extend its functionality as you wish.
type sysProvincesDao struct {
	internalSysProvincesDao
}

var (
	// SysProvinces is globally common accessible object for table sys_provinces operations.
	SysProvinces = sysProvincesDao{
		internal.NewSysProvincesDao(),
	}
)

// Fill with you ideas below.
