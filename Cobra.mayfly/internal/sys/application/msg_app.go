package application

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"learn_zinx/Cobra.mayfly/pkg/ws"
	"time"
)

type Msg interface {
	GetPageList(condition *entity.Msg, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Create(msg *entity.Msg)

	// 创建消息，并通过ws发送
	CreateAndSend(la *model.LoginAccount, msg *ws.Msg)
}

func newMsgApp(msgRepo repository.Msg) Msg {
	return &msgAppImpl{
		msgRepo: msgRepo,
	}
}

type msgAppImpl struct {
	msgRepo repository.Msg
}

// 查询消息列表
func (a *msgAppImpl) GetPageList(condition *entity.Msg, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return a.msgRepo.GetPageList(condition, pageParam, toEntity)
}

func (a *msgAppImpl) Create(msg *entity.Msg) {
	a.msgRepo.Insert(msg)
}

func (a *msgAppImpl) CreateAndSend(la *model.LoginAccount, wmsg *ws.Msg) {
	now := time.Now()
	msg := &entity.Msg{Type: 2, Msg: wmsg.Msg, RecipientId: int64(la.Id), CreateTime: &now, CreatorId: la.Id, Creator: la.Username}
	a.msgRepo.Insert(msg)
	ws.SendMsg(la.Id, wmsg)
}
