package cake

import (
	"context"
	"testing"
	"time"

	"github.com/tnngo/cake/gredis"
	"go.uber.org/zap"
)

func Test_loadRedis(t *testing.T) {
	type args struct {
		rc *redisConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				rc: &redisConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gredis.RDB.HSet(context.Background(), "aaa", "bbb", time.Now())
			sc := gredis.RDB.HGet(context.Background(), "aaa", "bbb")
			getTime, err := sc.Time()
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			zap.S().Debug(getTime)
		})
	}
}
