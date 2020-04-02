package logging

import (
	"reflect"
	"time"
	"zack.wtf/lib/chrono"
)

type (
	Log struct {
		Level   Level
		Context Data
		Sink    Sink
		Clock   chrono.Clock
	}

	Event struct {
		Level   Level
		Message string
		Context Data
		Details Data
		Time    time.Time
	}

	Sink     interface{ Handle(*Event) }
	SinkFunc func(*Event)
	Data     map[string]interface{}
	Level    int
)

const (
	None Level = iota
	Error
	Warn
	Info
	Debug
)

var (
	Default = New(Console{})

	KeyLevel   = "level"
	KeyMessage = "message"
	KeyTime    = "time"
	KeyType    = "type"

	MessagesUTC = true
)

func (this *Event) Data() Data {
	result := this.Context.Clone()

	for k, v := range this.Details {
		result[k] = v
	}

	result[KeyLevel] = this.Level
	result[KeyMessage] = this.Message
	result[KeyTime] = this.Time

	return result
}

func New(sink Sink) *Log {
	return &Log{
		Level: Info,
		Sink:  sink,
		Clock: chrono.SystemClock,
	}
}

func (this *Log) Error(template string, props ...interface{}) {
	this.Emit(Error, template, props)
}

func (this *Log) Warn(template string, props ...interface{}) {
	this.Emit(Warn, template, props)
}

func (this *Log) Info(template string, props ...interface{}) {
	this.Emit(Info, template, props)
}

func (this *Log) Debug(template string, props ...interface{}) {
	this.Emit(Debug, template, props)
}

func (this *Log) For(key string, value interface{}) *Log {
	context := this.Context.Clone()
	context[key] = value

	return &Log{
		Level:   this.Level,
		Context: context,
		Sink:    this.Sink,
		Clock:   this.Clock,
	}
}

func (this *Log) ForType(t interface{}) *Log {
	return this.For(KeyType, reflect.TypeOf(t))
}

func (this *Log) Emit(level Level, template string, props []interface{}) {
	if this.Level < level {
		return
	}

	message, details := convert(template, props)

	this.Sink.Handle(&Event{
		Level:   level,
		Message: message,
		Context: this.Context,
		Details: details,
		Time:    this.Clock.Now(),
	})
}

func (this Data) Clone() Data {
	clone := make(Data)
	for k, v := range this {
		clone[k] = v
	}
	return clone
}

func (this SinkFunc) Handle(e *Event) { this(e) }
