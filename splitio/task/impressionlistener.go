package task

import (
	"errors"
	"github.com/splitio/go-agent/log"
	"github.com/splitio/go-agent/splitio/util"
	"time"
)

var impressionListenerStream = make(chan *ImpressionBulk, util.ImpressionListenerMainQueueSize)

func QueueImpressionsForListener(impressions *ImpressionBulk) error {
	select {
	case impressionListenerStream <- impressions:
		return nil
	default:
		return errors.New("Impression listener queue is full. Last bulk not added")
	}

}

func taskPostImpressionsToListener(ilSubmitter util.ImpressionListenerSubmitter, failedQueue chan *ImpressionBulk) {
	failedImpressions := true
	for failedImpressions {
		select {
		case msg := <-failedQueue:
			err := ilSubmitter.Post(msg.Data, msg.SdkVersion, msg.MachineIP, "")
			if err != nil {
				msg.attempt++
				if msg.attempt < 3 {
					failedQueue <- msg
				}
				time.Sleep(time.Millisecond * 100)
			}
		default:
			failedImpressions = false
		}
	}
	msg := <-impressionListenerStream
	err := ilSubmitter.Post(msg.Data, msg.SdkVersion, msg.MachineIP, msg.MachineName)
	if err != nil {
		select {
		case failedQueue <- msg:
		default:
			log.Error.Println("Impression listener queue is full. " +
				"Impressions will be dropped until the listener enpoint is restored.")
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func PostImpressionsToListener(ilSubmitter util.ImpressionListenerSubmitter) {
	var failedQueue = make(chan *ImpressionBulk, util.ImpressionListenerFailedQueueSize)
	for {
		taskPostImpressionsToListener(ilSubmitter, failedQueue)
	}
}
