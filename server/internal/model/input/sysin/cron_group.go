// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
//
package sysin

import (
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
)

// CronGroupMaxSortInp 最大排序
type CronGroupMaxSortInp struct {
	Id int64
}

type CronGroupMaxSortModel struct {
	Sort int
}

// CronGroupEditInp 修改/新增字典数据
type CronGroupEditInp struct {
	entity.SysCronGroup
}
type CronGroupEditModel struct{}

// CronGroupDeleteInp 删除字典类型
type CronGroupDeleteInp struct {
	Id interface{}
}
type CronGroupDeleteModel struct{}

// CronGroupViewInp 获取信息
type CronGroupViewInp struct {
	Id int64
}

type CronGroupViewModel struct {
	entity.SysCronGroup
}

// CronGroupListInp 获取列表
type CronGroupListInp struct {
	form.PageReq
	form.RangeDateReq
	form.StatusReq
	Name string
}

type CronGroupListModel struct {
	entity.SysCronGroup
}

// CronGroupStatusInp 更新状态
type CronGroupStatusInp struct {
	entity.SysCronGroup
}
type CronGroupStatusModel struct{}

// CronGroupSelectInp 选项
type CronGroupSelectInp struct {
}

type CronGroupSelectModel struct {
	List []*CronGroupTree `json:"list"`
}

type CronGroupTree struct {
	entity.SysCronGroup
	Disabled bool             `json:"disabled"  dc:"是否禁用"`
	Label    string           `json:"label"     dc:"标签"`
	Value    int64            `json:"value"     dc:"键值"`
	Key      int64            `json:"key"       dc:"键名"`
	Children []*CronGroupTree `json:"children"  dc:"子级"`
}
