package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/psych-senior/biz/adaptor/cmd"
	"github.com/xh-polaris/psych-senior/biz/infrastructure/mapper/history"
)

type IHistoryService interface {
	ListHistory(ctx context.Context, req *cmd.ListHistoryReq) (*cmd.ListHistoryResp, error)
}

type HistoryService struct {
	HistoryMapper *history.MongoMapper
}

var HistoryServiceSet = wire.NewSet(
	wire.Struct(new(HistoryService), "*"),
	wire.Bind(new(IHistoryService), new(*HistoryService)),
)

func (s *HistoryService) ListHistory(ctx context.Context, req *cmd.ListHistoryReq) (*cmd.ListHistoryResp, error) {
	data, total, err := s.HistoryMapper.FindMany(ctx, &req.Paging)
	if err != nil {
		return nil, err
	}

	his := make([]*cmd.History, 0, len(data))
	for _, h := range data {
		dia := make([]*cmd.Dialog, 0, len(h.Dialogs))
		for _, d := range h.Dialogs {
			if d == nil {
				continue
			}
			dia = append(dia, &cmd.Dialog{
				Role:    d.Role,
				Content: d.Content,
			})
		}
		ch := &cmd.History{
			ID:        h.ID.Hex(),
			Dialogs:   dia,
			StartTime: h.StartTime.Unix(),
			EndTime:   h.EndTime.Unix(),
			Report:    &cmd.Report{},
		}
		if h.Report != nil {
			if err := copier.Copy(ch.Report, h.Report); err != nil {
				return nil, err
			}
		}

		his = append(his, ch)
	}
	return &cmd.ListHistoryResp{
		Code:    0,
		Msg:     "success",
		History: his,
		Total:   total,
	}, nil
}
