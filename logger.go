package log

import (
    "context"
    "io"
    "os"
    "sync"

    "github.com/sirupsen/logrus"
)

// Logger 定义了自定义日志库的核心接口
type Logger interface {
    // Debugf 标准方法
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
    Fatalf(format string, args ...any)

    // DebugContextf 带上下文（Context）方法
    DebugContextf(ctx context.Context, format string, args ...any)
    InfoContextf(ctx context.Context, format string, args ...any)
    WarnContextf(ctx context.Context, format string, args ...any)
    ErrorContextf(ctx context.Context, format string, args ...any)
    FatalContextf(ctx context.Context, format string, args ...any)

    // 动态配置方法
    SetLevel(level logrus.Level)
    SetOutput(output io.Writer)
    SetFormatter(format LogFormat)
}

// LogrusLogger 是 Logger 接口的 Logrus 实现
type LogrusLogger struct {
    *logrus.Logger
    config Config
    mu     sync.RWMutex // 用于保护配置修改
}

// NewLogger 创建并返回一个新的 Logger 实例
func NewLogger(cfg Config) (Logger, error) {
    l := logrus.New()

    // 设置日志级别
    l.SetLevel(cfg.Level)

    // 设置输出目标
    if cfg.FilePath != "" {
        file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            return nil, err
        }
        l.SetOutput(file)
    } else {
        l.SetOutput(cfg.Output)
    }

    // 设置日志格式
    if cfg.EnableJSON || cfg.Format == FormatJSON {
        jsonFormatter := &logrus.JSONFormatter{
            TimestampFormat:   cfg.TimestampFormat,
            DisableTimestamp:  false,
            DisableHTMLEscape: true,
            FieldMap:          nil,
            CallerPrettyfier:  nil,
            PrettyPrint:       cfg.JSONPretty, // JSON格式美化输出
        }
        l.SetFormatter(jsonFormatter)

    } else {
        l.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:   true,
            TimestampFormat: cfg.TimestampFormat,
            ForceColors:     true, // 强制终端颜色
            DisableColors:   false,
        })
    }

    // 添加 Caller Hook,
    if cfg.ReportCaller {
        l.AddHook(NewCallerHook(CallerSkipFrames))
    }

    return &LogrusLogger{
        Logger: l,
        config: cfg,
    }, nil
}

// Debugf --- Logger 接口实现 ---
// 为了 SkipFrames 一致，需要保持和 XXContextf 一样的调用方式
func (l *LogrusLogger) Debugf(format string, args ...any) {
    l.Logger.WithContext(context.Background()).Debugf(format, args...)
}

func (l *LogrusLogger) Infof(format string, args ...any) {
    l.Logger.WithContext(context.Background()).Infof(format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...any) {
    l.Logger.WithContext(context.Background()).Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...any) {
    l.Logger.WithContext(context.Background()).Errorf(format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...any) {
    l.Logger.WithContext(context.Background()).Fatalf(format, args...)
}

// --- 带上下文（Context）方法实现 ---

// Note: Logrus 本身没有直接的 WithContext 方法来传递 context 到 formatter 或 hook。
// 通常，我们会将 context 中的特定值（如 request ID, trace ID）提取出来作为日志字段。
// 这里我们通过 WithField("context", ctx) 来简单演示，实际应用中会更精细地处理。
// 更好的做法是：从 context 中提取 trace/span ID，并使用 WithField 添加。

// addContextFields 从 Context 中提取预定义的字段并添加到 Logrus Entry
func (l *LogrusLogger) addContextFields(ctx context.Context, entry *logrus.Entry) *logrus.Entry {
    if reqID, ok := GetRequestID(ctx); ok {
        entry = entry.WithField(string(RequestIDKey), reqID)
    }
    if userID, ok := GetUserID(ctx); ok {
        entry = entry.WithField(string(UserIDKey), userID)
    }
    if traceID, ok := GetTraceID(ctx); ok {
        entry = entry.WithField(string(TraceIDKey), traceID)
    }
    if spanID, ok := GetSpanID(ctx); ok {
        entry = entry.WithField(string(SpanIDKey), spanID)
    }
    // 处理自定义字段
    if customFields, ok := GetCustomFields(ctx); ok {
        for k, v := range customFields {
            entry = entry.WithField(k, v)
        }
    }
    return entry
}

func (l *LogrusLogger) DebugContextf(ctx context.Context, format string, args ...any) {
    entry := l.Logger.WithContext(ctx)
    entry = l.addContextFields(ctx, entry) // 添加上下文字段
    entry.Debugf(format, args...)
}

func (l *LogrusLogger) InfoContextf(ctx context.Context, format string, args ...any) {
    entry := l.Logger.WithContext(ctx)
    entry = l.addContextFields(ctx, entry) // 添加上下文字段
    entry.Infof(format, args...)
}

func (l *LogrusLogger) WarnContextf(ctx context.Context, format string, args ...any) {
    entry := l.Logger.WithContext(ctx)
    entry = l.addContextFields(ctx, entry) // 添加上下文字段
    entry.Warnf(format, args...)
}

func (l *LogrusLogger) ErrorContextf(ctx context.Context, format string, args ...any) {
    entry := l.Logger.WithContext(ctx)
    entry = l.addContextFields(ctx, entry) // 添加上下文字段
    entry.Errorf(format, args...)
}

func (l *LogrusLogger) FatalContextf(ctx context.Context, format string, args ...any) {
    entry := l.Logger.WithContext(ctx)
    entry = l.addContextFields(ctx, entry) // 添加上下文字段
    entry.Fatalf(format, args...)
}

// --- 动态配置方法实现 ---

func (l *LogrusLogger) SetLevel(level logrus.Level) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.Logger.SetLevel(level)
    l.config.Level = level
}

func (l *LogrusLogger) SetOutput(output io.Writer) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.Logger.SetOutput(output)
    l.config.Output = output
    l.config.FilePath = "" // 如果手动设置了输出，则清空文件路径
}

func (l *LogrusLogger) SetFormatter(format LogFormat) {
    l.mu.Lock()
    defer l.mu.Unlock()

    l.config.Format = format
    if format == FormatJSON {
        l.Logger.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: l.config.TimestampFormat,
            PrettyPrint:     l.config.JSONPretty,
        })
        l.config.EnableJSON = true
    } else {
        l.Logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:   true,
            TimestampFormat: l.config.TimestampFormat,
            ForceColors:     true,
            DisableColors:   false,
        })
        l.config.EnableJSON = false
    }

}
