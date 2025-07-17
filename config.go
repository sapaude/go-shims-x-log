package log

import (
    "io"
    "os"

    "github.com/sirupsen/logrus"
)

// LogFormat 定义日志输出格式
type LogFormat string

const (
    FormatText LogFormat = "text"
    FormatJSON LogFormat = "json"
)

// Config 定义日志库的配置参数
type Config struct {
    Level           logrus.Level // 日志级别
    Format          LogFormat    // 日志输出格式 (text/json)
    Output          io.Writer    // 日志输出目标 (例如 os.Stdout, 文件)
    FilePath        string       // 如果输出到文件，指定文件路径
    EnableJSON      bool         // 是否启用 JSON 格式输出
    JSONPretty      bool         // JSON美化输出
    ReportCaller    bool         // 是否报告调用者信息 (文件, 行号, 函数名)
    TimestampFormat string       // 时间戳格式，默认为 time.RFC3339Nano
}

// DefaultConfig 返回一个默认的日志配置
func DefaultConfig() Config {
    return Config{
        Level:           logrus.InfoLevel,
        Format:          FormatJSON,
        Output:          os.Stdout,
        FilePath:        "",
        EnableJSON:      false,
        ReportCaller:    true, // 默认开启调用者信息
        TimestampFormat: "2006/01/02 15:04:05.000",
    }
}
