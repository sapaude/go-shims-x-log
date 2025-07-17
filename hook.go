package log

import (
    "fmt"
    "runtime"
    "strings"

    "github.com/sirupsen/logrus"
)

const (
    // CallerFileFieldKey 是日志中存储调用者信息的字段名
    CallerFileFieldKey = "file"
    CallerFuncFieldKey = "func"

    EchoCallerSkipFrames = 9
    CallerSkipFrames     = EchoCallerSkipFrames
)

// CallerHook 是一个 Logrus Hook，用于添加调用者信息（文件、行号、函数名）
type CallerHook struct {
    // SkipFrames 决定向上跳过多少个栈帧来找到真正的调用者
    // 默认情况下，我们需要跳过 Logrus 内部调用和我们自己的封装层
    SkipFrames int
}

// NewCallerHook 创建一个新的 CallerHook 实例
func NewCallerHook(skipFrames int) *CallerHook {
    return &CallerHook{
        SkipFrames: skipFrames,
    }
}

// Levels 返回 Hook 应该触发的日志级别
func (hook *CallerHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// Fire 在日志事件发生时被调用
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
    // 向上跳过 hook.Fire, logrus.Entry.log, my_logger.Logger 方法, 以及 Logrus 内部的调用
    // 具体的跳过帧数可能需要根据实际封装层级进行微调
    pc, file, line, ok := runtime.Caller(hook.SkipFrames)
    if !ok {
        return nil
    }

    funcName := runtime.FuncForPC(pc).Name()
    // 简化函数名，去除包路径
    lastSlash := strings.LastIndex(funcName, "/")
    if lastSlash != -1 {
        funcName = funcName[lastSlash+1:]
    }
    lastDot := strings.LastIndex(funcName, ".")
    if lastDot != -1 {
        funcName = funcName[lastDot+1:]
    }

    // 格式化调用者信息
    entry.Data[CallerFileFieldKey] = fmt.Sprintf("file://%s:%d", file, line)
    entry.Data[CallerFuncFieldKey] = fmt.Sprintf("%s()", funcName)
    return nil

}
