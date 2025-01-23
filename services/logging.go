package services

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"

	"github.com/jmarren/deepfried/util"
)

type LogService struct {
}

func NewLogService() *LogService {
	l := new(LogService)
	return l
}

type key int

const requestIdKey key = 0

var requestNo int = 0

func (l *LogService) generateRequestId() int {
	requestNo = requestNo + 1
	return requestNo
}

func (l *LogService) WithId(ctx context.Context) context.Context {
	requestId := l.generateRequestId()
	ctx = context.WithValue(ctx, requestIdKey, requestId)
	return ctx
}

func (l *LogService) GetReqId(ctx context.Context) (int, error) {
	id, ok := ctx.Value(requestIdKey).(int)
	if !ok {
		return -1, fmt.Errorf("value not of type int")
	}
	return id, nil
}

func (l *LogService) LogId(ctx context.Context) {
	id, err := l.GetReqId(ctx)
	util.EMsg(err, "getting request Id")
	if err == nil {
		color.Blue("[ReqId: %d]\n", id)
		color.Unset()
	}
}

func (l *LogService) Log(ctx context.Context, msg string) {
	if os.Getenv("loglevel") == "high" {
		log.Println()
		l.LogId(ctx)
		log.Printf("[Msg: %s]\n", msg)
	}
}
