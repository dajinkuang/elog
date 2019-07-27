package elog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"runtime"
	"time"

	"github.com/dajinkuang/util"
	"github.com/labstack/gommon/color"
)

const defaultTopic = "default_topic"

// 在main中修改
func SetTopic(topic string, absolutePath string) {
	// 考虑重入
	if __eJsonLog != nil {
		__eJsonLog.Close()
		__eJsonLogErrorAbove.Close()
	}
	dir := "/tmp/go/log"
	if len(absolutePath) > 0 {
		dir = absolutePath
	}

	file, err := NewFileBackend(dir, topic+".log_json_std")
	if err != nil {
		panic(err)
	}
	__eJsonLog = NewEJsonLog(file, topic)
	SetLogger(__eJsonLog)

	fileErrorAbove, err := NewFileBackend(dir, topic+".log_json_error_above")
	if err != nil {
		panic(err)
	}
	__eJsonLogErrorAbove = NewEJsonLog(fileErrorAbove, topic)
	SetLoggerErrorAbove(__eJsonLogErrorAbove)
}

var _ Logger = __eJsonLog

// 根据级别打印所有日志
var __eJsonLog *eJsonLog

// 只打印 ERROR FATAL 日志
var __eJsonLogErrorAbove *eJsonLog

func GetJsonELog() *eJsonLog {
	if __eJsonLog == nil {
		SetTopic(defaultTopic, "")
	}
	return __eJsonLog
}

func GetJsonELogErrorAbove() *eJsonLog {
	if __eJsonLogErrorAbove == nil {
		SetTopic(defaultTopic, "")
	}
	return __eJsonLogErrorAbove
}

type eJsonLog struct {
	prefix string
	level  Lvl
	output io.Writer
	levels []string
	color  *color.Color
	dw     *elogWriter
}

func NewEJsonLog(w io.WriteCloser, topic string) *eJsonLog {
	if len(topic) <= 0 {
		topic = defaultTopic
	}
	l := &eJsonLog{
		level:  INFO,
		prefix: topic,
		color:  color.New(),
	}
	l.initLevels()
	l.dw = NewElogWriter(w)
	l.SetOutput(l.dw)
	l.SetLevel(INFO)
	return l
}

func (p *eJsonLog) With(ctx context.Context, kv ...interface{}) context.Context {
	om := FromContext(ctx)
	if om == nil {
		om = NewOrderMap()
	}
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	for i := 0; i < len(kv); i += 2 {
		om.Set(fmt.Sprintf("%v", kv[i]), kv[i+1])
	}
	return setContext(ctx, om)
}

func (p *eJsonLog) Debug(ctx context.Context, msg interface{}, kv ...interface{}) {
	p.logJSON(DEBUG, ctx, msg, kv...)
}

func (p *eJsonLog) Info(ctx context.Context, msg interface{}, kv ...interface{}) {
	p.logJSON(INFO, ctx, msg, kv...)
}

func (p *eJsonLog) Warn(ctx context.Context, msg interface{}, kv ...interface{}) {
	p.logJSON(WARN, ctx, msg, kv...)
}

func (p *eJsonLog) Error(ctx context.Context, msg interface{}, kv ...interface{}) {
	p.logJSON(ERROR, ctx, msg, kv...)
}

func (p *eJsonLog) Fatal(ctx context.Context, msg interface{}, kv ...interface{}) {
	p.logJSON(FATAL, ctx, msg, kv...)
	panic(msg)
}

// kv 应该是成对的 数据, 类似: name,张三,age,10,...
func (p *eJsonLog) logJSON(v Lvl, ctx context.Context, msg interface{}, kv ...interface{}) (err error) {
	if v < p.level {
		return nil
	}
	om := NewOrderMap()
	_, file, line, _ := runtime.Caller(3)
	file = p.getFilePath(file)
	om.Set("prefix", p.Prefix())
	om.Set("level", p.levels[v])
	om.Set("cur_time", time.Now().Format(time.RFC3339Nano))
	om.Set("file", file)
	om.Set("line", line)
	localMachineIPV4, _ := util.LocalMachineIPV4()
	om.Set("local_machine_ipv4", localMachineIPV4)
	om.Set(TraceId, ValueFromOM(ctx, TraceId))
	om.Set(SpanId, ValueFromOM(ctx, SpanId))
	om.Set(ParentId, ValueFromOM(ctx, ParentId))
	om.Set(UserRequestIp, ValueFromOM(ctx, UserRequestIp))
	om.Set("log_desc", msg)
	om.AddVals(FromContext(ctx))
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	for i := 0; i < len(kv); i += 2 {
		om.Set(fmt.Sprintf("%v", kv[i]), kv[i+1])
	}
	str, _ := json.Marshal(om)
	str = append(str, []byte("\n")...)
	_, err = p.Output().Write(str)
	return
}

func (p *eJsonLog) getFilePath(file string) string {
	dir, base := path.Dir(file), path.Base(file)
	return path.Join(path.Base(dir), base)
}

func (p *eJsonLog) Close() error {
	if p.dw != nil {
		p.dw.Close()
		p.dw = nil
	}
	return nil
}

func (p *eJsonLog) EnableDebug(b bool) {
	if b {
		p.SetLevel(DEBUG)
	} else {
		p.SetLevel(INFO)
	}
}

type Lvl uint8

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

func (l *eJsonLog) initLevels() {
	l.levels = []string{
		"-",
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
	}
}

func (l *eJsonLog) Prefix() string {
	return l.prefix
}

func (l *eJsonLog) SetPrefix(p string) {
	l.prefix = p
}

func (l *eJsonLog) Level() Lvl {
	return l.level
}

func (l *eJsonLog) SetLevel(v Lvl) {
	l.level = v
}

func (l *eJsonLog) Output() io.Writer {
	return l.output
}

func (l *eJsonLog) SetOutput(w io.Writer) {
	l.output = w
}

func (l *eJsonLog) Color() *color.Color {
	return l.color
}
