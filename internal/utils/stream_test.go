package utils

import (
	"bytes"
	"io"
	"testing"
)

func setup() {

}

type TestContext struct {
	writerStreamEnabled  *bytes.Buffer
	writerStreamDisabled *bytes.Buffer
}

func setupTest(t *testing.T) *TestContext {
	return &TestContext{
		writerStreamEnabled:  new(bytes.Buffer),
		writerStreamDisabled: new(bytes.Buffer),
	}
}

func TestGetOutputStream(t *testing.T) {
	t.Run("Enabled write stream", func(t *testing.T) {
		ctx := setupTest(t)

		data := "hello world!"

		outStream := GetOutputStream(true, ctx.writerStreamEnabled, ctx.writerStreamDisabled)

		outStream.Write([]byte(data))

		writeBytesEnabled, err := io.ReadAll(ctx.writerStreamEnabled)
		if err != nil {
			t.Fatalf("unexpected err reading enabled writer bytes: %v", err)
		}

		writeBytesDisabled, err := io.ReadAll(ctx.writerStreamDisabled)
		if err != nil {
			t.Fatalf("unexpected err reading disabled writer bytes: %v", err)
		}

		writeStrEnabled := string(writeBytesEnabled)
		writeStrDisabled := string(writeBytesDisabled)

		if data != writeStrEnabled {
			t.Fatalf("unexpected data in enabled write stream, got '%s' want '%s'", writeStrEnabled, data)
		}

		if len(writeStrDisabled) > 0 {
			t.Fatalf("unexpected data in disabled write stream, got '%s' want ''", writeStrDisabled)
		}
	})

	t.Run("Disabled write stream", func(t *testing.T) {
		ctx := setupTest(t)

		data := "hello world disabled!"

		outStream := GetOutputStream(false, ctx.writerStreamEnabled, ctx.writerStreamDisabled)

		outStream.Write([]byte(data))

		writeBytesEnabled, err := io.ReadAll(ctx.writerStreamEnabled)
		if err != nil {
			t.Fatalf("unexpected err reading enabled writer bytes: %v", err)
		}

		writeBytesDisabled, err := io.ReadAll(ctx.writerStreamDisabled)
		if err != nil {
			t.Fatalf("unexpected err reading disabled writer bytes: %v", err)
		}

		writeStrEnabled := string(writeBytesEnabled)
		writeStrDisabled := string(writeBytesDisabled)

		if data != writeStrDisabled {
			t.Fatalf("unexpected data in disabled write stream, got '%s' want '%s'", writeStrEnabled, data)
		}

		if len(writeStrEnabled) > 0 {
			t.Fatalf("unexpected data in enabled write stream, got '%s' want ''", writeStrDisabled)
		}
	})
}
