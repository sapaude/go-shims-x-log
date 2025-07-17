package log

import (
    "context"
    "sync"

    "github.com/sirupsen/logrus"
)

var (
    globalLogger     Logger
    globalLoggerOnce sync.Once
)

// InitGlobalLogger 初始化全局 Logger 实例。
// 只能被调用一次，后续调用将被忽略。
func InitGlobalLogger(cfg Config) {
    globalLoggerOnce.Do(func() {
        l, err := NewLogger(cfg)
        if err != nil {
            // 如果初始化失败，退回到一个最简单的 Logrus 实例，并打印错误
            logrus.SetOutput(cfg.Output)
            logrus.SetLevel(logrus.ErrorLevel)
            logrus.Errorf("Failed to initialize custom logger: %v. Falling back to basic logrus.", err)
            globalLogger = &LogrusLogger{Logger: logrus.StandardLogger(), config: cfg} // 使用标准 Logrus 作为回退
            return
        }
        globalLogger = l
    })
}

// GetGlobalLogger 获取全局 Logger 实例。
// 如果尚未初始化，将使用 DefaultConfig() 进行初始化。
func GetGlobalLogger() Logger {
    globalLoggerOnce.Do(func() {
        cfg := DefaultConfig()
        l, err := NewLogger(cfg)
        if err != nil {
            logrus.SetOutput(cfg.Output)
            logrus.SetLevel(logrus.ErrorLevel)
            logrus.Errorf("Failed to initialize default global logger: %v. Falling back to basic logrus.", err)
            globalLogger = &LogrusLogger{Logger: logrus.StandardLogger(), config: cfg}
            return
        }
        globalLogger = l
    })
    return globalLogger
}

// --- 全局日志方法 (方便直接调用) ---

func Debugf(format string, args ...any) {
    GetGlobalLogger().Debugf(format, args...)
}

func Infof(format string, args ...any) {
    GetGlobalLogger().Infof(format, args...)
}

func Warnf(format string, args ...any) {
    GetGlobalLogger().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
    GetGlobalLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
    GetGlobalLogger().Fatalf(format, args...)
}

func DebugContextf(ctx context.Context, format string, args ...any) {
    GetGlobalLogger().DebugContextf(ctx, format, args...)
}

func InfoContextf(ctx context.Context, format string, args ...any) {
    GetGlobalLogger().InfoContextf(ctx, format, args...)
}

func WarnContextf(ctx context.Context, format string, args ...any) {
    GetGlobalLogger().WarnContextf(ctx, format, args...)
}

func ErrorContextf(ctx context.Context, format string, args ...any) {
    GetGlobalLogger().ErrorContextf(ctx, format, args...)
}

func FatalContextf(ctx context.Context, format string, args ...any) {
    GetGlobalLogger().FatalContextf(ctx, format, args...)
}
