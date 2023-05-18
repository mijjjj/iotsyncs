// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/location"
	"hotgo/internal/model"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/validate"
	"time"
)

type sSysSmsLog struct{}

func NewSysSmsLog() *sSysSmsLog {
	return &sSysSmsLog{}
}

func init() {
	service.RegisterSysSmsLog(NewSysSmsLog())
}

// Delete 删除
func (s *sSysSmsLog) Delete(ctx context.Context, in sysin.SmsLogDeleteInp) (err error) {
	_, err = dao.SysSmsLog.Ctx(ctx).Where("id", in.Id).Delete()
	return
}

// Edit 修改/新增
func (s *sSysSmsLog) Edit(ctx context.Context, in sysin.SmsLogEditInp) (err error) {
	if in.Ip == "" {
		err = gerror.New("ip不能为空")
		return
	}

	// 修改
	if in.Id > 0 {
		_, err = dao.SysSmsLog.Ctx(ctx).Where("id", in.Id).Data(in).Update()
		return
	}

	// 新增
	_, err = dao.SysSmsLog.Ctx(ctx).Data(in).Insert()
	return
}

// Status 更新部门状态
func (s *sSysSmsLog) Status(ctx context.Context, in sysin.SmsLogStatusInp) (err error) {
	if in.Id <= 0 {
		err = gerror.New("ID不能为空")
		return
	}

	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}

	if !validate.InSliceInt(consts.StatusSlice, in.Status) {
		err = gerror.New("状态不正确")
		return
	}

	// 修改
	_, err = dao.SysSmsLog.Ctx(ctx).Where("id", in.Id).Data("status", in.Status).Update()
	return
}

// MaxSort 最大排序
func (s *sSysSmsLog) MaxSort(ctx context.Context, in sysin.SmsLogMaxSortInp) (res *sysin.SmsLogMaxSortModel, err error) {
	if in.Id > 0 {
		if err = dao.SysSmsLog.Ctx(ctx).Where("id", in.Id).Order("sort desc").Scan(&res); err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return
		}
	}

	if res == nil {
		res = new(sysin.SmsLogMaxSortModel)
	}

	res.Sort = form.DefaultMaxSort(ctx, res.Sort)
	return
}

// View 获取指定字典类型信息
func (s *sSysSmsLog) View(ctx context.Context, in sysin.SmsLogViewInp) (res *sysin.SmsLogViewModel, err error) {
	if err = dao.SysSmsLog.Ctx(ctx).Where("id", in.Id).Scan(&res); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}
	return
}

// List 获取列表
func (s *sSysSmsLog) List(ctx context.Context, in sysin.SmsLogListInp) (list []*sysin.SmsLogListModel, totalCount int, err error) {
	mod := dao.SysSmsLog.Ctx(ctx)

	if in.Mobile != "" {
		mod = mod.WhereLike("mobile", "%"+in.Mobile+"%")
	}

	if in.Ip != "" {
		mod = mod.Where("ip", in.Ip)
	}

	if in.Event != "" {
		mod = mod.Where("event", in.Event)
	}

	if in.Status > 0 {
		mod = mod.Where("status", in.Status)
	}

	totalCount, err = mod.Count()
	if err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return list, totalCount, err
	}

	if totalCount == 0 {
		return list, totalCount, nil
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id desc").Scan(&list); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return list, totalCount, err
	}

	return list, totalCount, err
}

// SendCode 发送验证码
func (s *sSysSmsLog) SendCode(ctx context.Context, in sysin.SendCodeInp) (err error) {
	if in.Event == "" {
		err = gerror.New("事件不能为空")
		return
	}

	if in.Mobile == "" {
		err = gerror.New("手机号不能为空")
		return
	}

	var models *entity.SysSmsLog
	if err = dao.SysSmsLog.Ctx(ctx).Where("event", in.Event).Where("mobile", in.Mobile).Scan(&models); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	config, err := service.SysConfig().GetSms(ctx)
	if err != nil {
		return
	}

	if in.Template, err = s.GetTemplate(ctx, in.Event, config); err != nil {
		return
	}

	if err = s.AllowSend(ctx, models, config); err != nil {
		return
	}

	if in.Code == "" {
		in.Code = grand.Digits(4)
	}

	//if err = sms.New(config.SmsDrive).SendCode(ctx, in, config); err != nil {
	//	return
	//}

	var data = new(entity.SysSmsLog)
	data.Event = in.Event
	data.Mobile = in.Mobile
	data.Code = in.Code
	data.Ip = location.GetClientIp(ghttp.RequestFromCtx(ctx))
	data.Status = consts.SmsStatusNotUsed
	data.CreatedAt = gtime.Now()
	data.UpdatedAt = gtime.Now()

	_, err = dao.SysSmsLog.Ctx(ctx).Data(data).Insert()
	return
}

// GetTemplate 获取指定短信模板
func (s *sSysSmsLog) GetTemplate(ctx context.Context, template string, config *model.SmsConfig) (val string, err error) {
	if template == "" {
		err = gerror.New("模板不能为空")
		return
	}
	if config == nil {
		config, err = service.SysConfig().GetSms(ctx)
		if err != nil {
			return
		}
	}

	switch config.SmsDrive {
	case consts.SmsDriveAliYun:
		if len(config.AliYunTemplate) == 0 {
			err = gerror.New("管理员还没有配置任何阿里云短信模板！")
			return
		}

		for _, v := range config.AliYunTemplate {
			if v.Key == template {
				return v.Value, nil
			}
		}

	case consts.SmsDriveTencent:
		if len(config.TencentTemplate) == 0 {
			err = gerror.New("管理员还没有配置任何腾讯云短信模板！")
			return
		}

		for _, v := range config.TencentTemplate {
			if v.Key == template {
				return v.Value, nil
			}
		}
	default:
		err = gerror.Newf("暂不支持短信驱动:%v", config.SmsDrive)
		return
	}

	return
}

// AllowSend 是否允许发送
func (s *sSysSmsLog) AllowSend(ctx context.Context, models *entity.SysSmsLog, config *model.SmsConfig) (err error) {
	if models == nil {
		return
	}

	if config == nil {
		if config, err = service.SysConfig().GetSms(ctx); err != nil {
			return
		}
	}

	if gtime.Now().Before(models.CreatedAt.Add(time.Second * time.Duration(config.SmsMinInterval))) {
		err = gerror.New("发送频繁，请稍后再试！")
		return
	}

	if config.SmsMaxIpLimit > 0 {
		count, err := dao.SysSmsLog.NowDayCount(ctx, models.Event, models.Mobile)
		if err != nil {
			return err
		}

		if count >= config.SmsMaxIpLimit {
			err = gerror.New("今天发送短信过多，请次日后再试！")
			return err
		}
	}

	return
}

// VerifyCode 效验验证码
func (s *sSysSmsLog) VerifyCode(ctx context.Context, in sysin.VerifyCodeInp) (err error) {
	if in.Event == "" {
		err = gerror.New("事件不能为空")
		return
	}

	if in.Mobile == "" {
		err = gerror.New("手机号不能为空")
		return
	}

	config, err := service.SysConfig().GetSms(ctx)
	if err != nil {
		return
	}

	var models *entity.SysSmsLog
	if err = dao.SysSmsLog.Ctx(ctx).Where("event", in.Event).Where("mobile", in.Mobile).Order("id desc").Scan(&models); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if models == nil {
		err = gerror.New("验证码错误")
		return
	}

	if models.Times >= 10 {
		err = gerror.New("验证码错误次数过多，请重新发送！")
		return
	}

	if in.Event != consts.SmsTemplateCode {
		if models.Status == consts.SmsStatusUsed {
			err = gerror.New("验证码已使用，请重新发送！")
			return
		}
	}

	if gtime.Now().After(models.CreatedAt.Add(time.Second * time.Duration(config.SmsCodeExpire))) {
		err = gerror.New("验证码已过期，请重新发送")
		return
	}

	if models.Code != in.Code {
		dao.SysSmsLog.Ctx(ctx).Where("id", models.Id).Increment("times", 1)
		err = gerror.New("验证码错误！")
		return
	}

	_, err = dao.SysSmsLog.Ctx(ctx).Where("id", models.Id).Data(g.Map{
		"times":      models.Times + 1,
		"status":     consts.SmsStatusUsed,
		"updated_at": gtime.Now(),
	}).Update()

	return
}
