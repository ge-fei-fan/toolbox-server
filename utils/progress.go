package utils

import (
	"fmt"
	"io"
)

// ProgressEventType defines transfer progress event type
type ProgressEventType int

const (
	// TransferStartedEvent transfer started, set TotalBytes
	TransferStartedEvent ProgressEventType = 1 + iota
	// TransferDataEvent transfer data, set ConsumedBytes and TotalBytes
	TransferDataEvent
	// TransferCompletedEvent transfer completed
	TransferCompletedEvent
	// TransferFailedEvent transfer encounters an error
	TransferFailedEvent
)

// ProgressEvent defines progress event
type ProgressEvent struct {
	ConsumedBytes int64
	TotalBytes    int64
	RwBytes       int64
	EventType     ProgressEventType
}

// ProgressListener listens progress change
type ProgressListener interface {
	ProgressChanged(event *ProgressEvent)
}

func newProgressEvent(eventType ProgressEventType, consumed, total int64, rwBytes int64) *ProgressEvent {
	return &ProgressEvent{
		ConsumedBytes: consumed,
		TotalBytes:    total,
		RwBytes:       rwBytes,
		EventType:     eventType}
}

// publishProgress
func publishProgress(listener ProgressListener, event *ProgressEvent) {
	if listener != nil && event != nil {
		listener.ProgressChanged(event)
	}
}

type teeReader struct {
	reader io.Reader
	//writer        io.Writer
	listener      ProgressListener
	consumedBytes int64
	totalBytes    int64
}

func TeeReader(reader io.Reader, totalBytes int64, listener ProgressListener) io.ReadCloser {
	return &teeReader{
		reader: reader,
		//writer:        writer,
		listener:      listener,
		consumedBytes: 0,
		totalBytes:    totalBytes,
	}
}
func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)

	// Read encountered error
	if err != nil && err != io.EOF {
		event := newProgressEvent(TransferFailedEvent, t.consumedBytes, t.totalBytes, 0)
		publishProgress(t.listener, event)
	}

	if n > 0 {
		t.consumedBytes += int64(n)
		// CRC
		//if t.writer != nil {
		//	if n, err := t.writer.Write(p[:n]); err != nil {
		//		return n, err
		//	}
		//}
		// Progress
		if t.listener != nil {
			event := newProgressEvent(TransferDataEvent, t.consumedBytes, t.totalBytes, int64(n))
			publishProgress(t.listener, event)
		}
		// Track
		//if t.tracker != nil {
		//	t.tracker.completedBytes = t.consumedBytes
		//}
	}

	return
}

func (t *teeReader) Close() error {
	if rc, ok := t.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}

// 定义进度条监听器。
type MyProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *MyProgressListener) ProgressChanged(event *ProgressEvent) {
	switch event.EventType {
	case TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case TransferDataEvent:
		if event.TotalBytes != 0 {
			fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
				event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
		}
	case TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}
