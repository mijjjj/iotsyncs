// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"hotgo/internal/dao/internal"
)

// internalAdminNoticeReadDao is internal type for wrapping internal DAO implements.
type internalAdminNoticeReadDao = *internal.AdminNoticeReadDao

// adminNoticeReadDao is the data access object for table admin_notice_read.
// You can define custom methods on it to extend its functionality as you wish.
type adminNoticeReadDao struct {
	internalAdminNoticeReadDao
}

var (
	// AdminNoticeRead is globally public accessible object for table admin_notice_read operations.
	AdminNoticeRead = adminNoticeReadDao{
		internal.NewAdminNoticeReadDao(),
	}
)

// Fill with you ideas below.
