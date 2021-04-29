package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	// bytes.BufferはWriterメソッドを持っている
	// bytes.Buffer
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		// bufに保持されているので検証してしてみる
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が検出されました", buf.String())
		}
	}
}

func TraceOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
