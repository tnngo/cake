package cake

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapConsole() {
	consoleWrite := zapcore.AddSync(io.MultiWriter(os.Stdout))
	consoleConfig := zap.NewProductionEncoderConfig()
	consoleConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	// color.
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfig),
		consoleWrite,
		zapcore.DebugLevel,
	)

	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
}
