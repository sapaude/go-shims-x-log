package test

import (
    "context"
    "fmt"
    "io"
    "testing"
    "time"

    "github.com/sapaude/go-shims/x/log"
    "github.com/sirupsen/logrus"
)

// simulateServiceCall 模拟一个服务调用，其中包含日志
func simulateServiceCall(ctx context.Context, data string) {
    log.InfoContextf(ctx, "Processing data: %s", data)
    log.DebugContextf(ctx, "Detailed debug info for data: %s", data)

    // 模拟一个错误
    err := fmt.Errorf("something went wrong with %s", data)
    log.ErrorContextf(ctx, "Failed to process data: %v", err)
}

func TestLogX(t *testing.T) {
    t.Log("--- Starting Logger Demo ---")

    // --- 1. 初始化全局 Logger (文本格式，输出到 stdout) ---
    t.Log("\n--- Demo 1: Global Logger (Text Format, Stdout) ---")
    cfg1 := log.DefaultConfig()
    cfg1.Level = logrus.DebugLevel // 设置为 Debug 级别以查看所有日志
    cfg1.ReportCaller = true       // 确保报告调用者信息
    cfg1.Format = log.FormatJSON
    log.InitGlobalLogger(cfg1)

    log.Debugf("This is a debug message from main function.")
    log.Infof("This is an info message.")
    log.Warnf("This is a warning message.")
    log.Errorf("This is an error message.")

    // 演示带上下文的日志
    ctx1 := log.WithRequestID(context.Background(), "req-12345")
    ctx1 = log.WithTraceID(ctx1, "trace-xyz")
    ctx1 = log.WithUserID(ctx1, "12345")
    ctx1 = log.WithCustomField(ctx1, "foo", "bar")
    ctx1 = log.WithCustomField(ctx1, "url", "https://x.com")

    simulateServiceCall(ctx1, "payload-A")

    // --- 2. 创建一个独立的 Logger 实例 (JSON 格式，输出到文件) ---
    t.Log("\n--- Demo 2: Independent Logger (JSON Format, File) ---")
    logFilePath := "app.log"
    cfg2 := log.DefaultConfig()
    cfg2.Level = logrus.InfoLevel
    cfg2.Format = log.FormatJSON
    cfg2.FilePath = logFilePath
    cfg2.ReportCaller = true
    cfg2.TimestampFormat = time.RFC3339 // 更短的时间戳格式

    fileLogger, err := log.NewLogger(cfg2)
    if err != nil {
        fmt.Printf("Error creating file logger: %v\n", err)
        return
    }
    defer func() {
        // 关闭文件句柄 (如果 NewLogger 内部打开了文件)
        if closer, ok := fileLogger.(*log.LogrusLogger).Out.(io.Closer); ok {
            closer.Close()
        }
        fmt.Printf("\nCheck log file: %s\n", logFilePath)
    }()

    fileLogger.Infof("This info message goes to the file.")
    fileLogger.Warnf("This warning message also goes to the file.")
    fileLogger.Debugf("This debug message will NOT appear in file (level is Info).") // 不会显示

    ctx2 := context.WithValue(context.Background(), log.TraceIDKey, "trace-XYZ")
    fileLogger.ErrorContextf(ctx2, "An error occurred in file logger context.")

    // --- 3. 动态调整全局 Logger 配置 ---
    t.Log("\n--- Demo 3: Dynamic Configuration Change ---")
    t.Log("Changing global logger level to Warn and format to JSON...")
    globalLog := log.GetGlobalLogger()
    globalLog.SetLevel(logrus.WarnLevel)
    globalLog.SetFormatter(log.FormatJSON)

    log.Debugf("This debug message will NOT appear now (level is Warn).") // 不会显示
    log.Infof("This info message will NOT appear now (level is Warn).")   // 不会显示
    log.Warnf("This warning message WILL appear in JSON format.")
    log.Errorf("This error message WILL appear in JSON format.")

    // 演示 Fatalf (会退出程序，所以放在最后或注释掉)
    // Fatalf("This is a fatal message, program will exit.")
    // FatalContextf(ctx1, "Fatal with context, program will exit.")

    t.Log("\n--- Logger Demo Finished ---")
}
