package example

import (
	"fmt"
	"testing"
	"toolbox-server/utils"
)

func TestSave(t *testing.T) {
	v := &VideoInstance{
		Bvid:     "BV1Ha4y1g76j",
		Cid:      1138698813,
		Title:    "test",
		SavePath: "D:\\自己写的工具\\toolbox-server\\aaa",
		Status:   1,
	}
	u := "https://cn-zjhz-cm-01-12.bilivideo.com/upgcxcode/13/88/1138698813/1138698813-1-30120.m4s?e=ig8euxZM2rNcNbdlhoNvNC8BqJIzNbfqXBvEqxTEto8BTrNvN0GvT90W5JZMkX_YN0MvXg8gNEV4NC8xNEV4N03eN0B5tZlqNxTEto8BTrNvNeZVuJ10Kj_g2UB02J0mN0B5tZlqNCNEto8BTrNvNC7MTX502C8f2jmMQJ6mqF2fka1mqx6gqj0eN0B599M=&uipk=5&nbs=1&deadline=1684948514&gen=playurlv2&os=bcache&oi=1882208457&trid=00005221b2e36ba044ef89cf6a68a5bb6618u&mid=34801693&platform=pc&upsig=cfcb9fe1b2740f5ebd7a2816623c08c2&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,mid,platform&cdnid=4070&bvc=vod&nettype=0&orderid=0,3&buvid=32BB4676-369D-6273-4209-76C3E20E619D99657infoc&build=0&agrr=1&bw=2061468&logo=80000000"
	err := v.saveToFile(u, ".mp4", &MyProgressListener{})
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 定义进度条监听器。
type MyProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *MyProgressListener) ProgressChanged(event *utils.ProgressEvent) {
	switch event.EventType {
	case utils.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case utils.TransferDataEvent:
		if event.TotalBytes != 0 {
			fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
				event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
		}
	case utils.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case utils.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}
