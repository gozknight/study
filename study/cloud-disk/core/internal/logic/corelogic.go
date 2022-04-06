package logic

import (
	"bytes"
	"cloud-disk/core/model"
	"cloud-disk/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoreLogic {
	return &CoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoreLogic) Core(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	resp = &types.Response{}
	user := make([]*models.UserBasic, 0)
	err = model.Engine.Find(&user)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "")
	if err != nil {
		log.Fatalln(err)
	}
	resp.Message = dst.String()
	fmt.Println(dst.String())
	return
}
