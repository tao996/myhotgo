package pay

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/payment"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/model/input/payin"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Query 订单查询
func (s *sPay) Query(ctx context.Context, in adminin.OrderQueryInp) (res *payin.PayQueryModel, err error) {
	if err = in.Filter(ctx); err != nil {
		return nil, err
	}

	where := g.Map{}
	if in.Id > 0 {
		where[dao.PayLog.Columns().Id] = in.Id
	}
	if in.OrderSn != "" {
		where[dao.PayLog.Columns().OrderSn] = in.OrderSn
	}
	if in.OutTradeNo != "" {
		where[dao.PayLog.Columns().OutTradeNo] = in.OutTradeNo
	}
	if in.TransactionId != "" {
		where[dao.PayLog.Columns().TransactionId] = in.TransactionId
	}

	var models *entity.PayLog
	if err = s.Model(ctx).Where(where).Scan(&models); err != nil {
		return
	}

	if models == nil {
		err = gerror.Newf("订单不存在，请检查")
		return
	}

	if models.PayStatus != consts.PayStatusWait {
		err = gerror.Newf("订单已被处理，请勿重复操作")
		return
	}

	data, err := payment.New(models.PayType).Query(ctx, payin.PayQueryInp{
		OutTradeNo:    models.OutTradeNo,
		TransactionId: models.TransactionId,
	})
	if err != nil {
		return nil, err
	}
	if err = s.doWithNotifyModel(ctx, models, data); err != nil {
		return nil, err
	}
	return &payin.PayQueryModel{PayLog: *models}, nil
}
