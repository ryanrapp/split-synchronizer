package task

import (
	"errors"
	"sync"
	"time"

	"github.com/splitio/split-synchronizer/log"
	"github.com/splitio/split-synchronizer/splitio/api"
	"github.com/splitio/split-synchronizer/splitio/nethelper"
	"github.com/splitio/split-synchronizer/splitio/recorder"
	"github.com/splitio/split-synchronizer/splitio/stats/counter"
	"github.com/splitio/split-synchronizer/splitio/stats/latency"
	"github.com/splitio/split-synchronizer/splitio/storage"
	"github.com/splitio/split-synchronizer/splitio/storage/redis"
)

var eventsIncoming chan string

var postEventsLatencies = latency.NewLatencyBucket()
var postEventsCounters = counter.NewCounter()
var postEventsLocalCounters = counter.NewCounter()

const totalPostAttemps = 3

// InitializeEvents initialiaze events task
func InitializeEvents(threads int) {
	eventsIncoming = make(chan string, threads)
}

// StopPostEvents stops PostEvents task sendding signal
func StopPostEvents() {
	select {
	case eventsIncoming <- "STOP":
	default:
	}
}

func taskPostEvents(tid int,
	recorderAdapter recorder.EventsRecorder,
	storageAdapter storage.EventStorage,
	bulkSize int64,
) {
	//[SDKVersion][MachineIP][MachineName]
	toSend := make(map[string]map[string]map[string][]api.EventDTO)

	storedEvents, err := storageAdapter.PopN(bulkSize) //PopN has a mutex, so this function can be async without issues
	if err != nil {
		log.Error.Println("(Task) Post Events fails fetching events from storage", err.Error())
		return
	}

	for _, stored := range storedEvents {

		if stored.Metadata.SDKVersion == "" {
			continue
		}

		sdk := stored.Metadata.SDKVersion
		ip := stored.Metadata.MachineIP
		mname := stored.Metadata.MachineName

		if ip == "" {
			ip = "unknown"
		}

		if mname == "" {
			mname = "unknown"
		}

		if toSend[sdk] == nil {
			toSend[sdk] = make(map[string]map[string][]api.EventDTO)
		}

		if toSend[sdk][ip] == nil {
			toSend[sdk][ip] = make(map[string][]api.EventDTO)
		}

		if toSend[sdk][ip][mname] == nil {
			toSend[sdk][ip][mname] = make([]api.EventDTO, 0)
		}

		toSend[sdk][ip][mname] = append(toSend[sdk][ip][mname], stored.Event)
	}

	for s, byIP := range toSend {
		for i, byName := range byIP {
			for n, bulk := range byName {

				var err = errors.New("") // forcing error to start "for" attempts
				attemps := 0
				for err != nil && attemps < totalPostAttemps {
					startTime := postEventsLatencies.StartMeasuringLatency()
					err = recorderAdapter.Post(bulk, s, i, n)
					if err != nil {
						log.Error.Println("Error posting events", err)
						postEventsLocalCounters.Increment("backend::request.error")
					} else {
						postEventsLatencies.RegisterLatency("backend::/api/events/bulk", startTime)
						postEventsLatencies.RegisterLatency("events.time", startTime)
						postEventsLocalCounters.Increment("backend::request.ok")
					}
					attemps++
					time.Sleep(nethelper.WaitForNextAttempt() * time.Second)
				}

			}
		}
	}
}

// PostEvents post events to Split Server task
func PostEvents(
	tid int,
	eventsRecorderAdapter recorder.EventsRecorder,
	eventsStorageAdapter storage.EventStorage,
	eventsRefreshRate int,
	eventsBulkSize int,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
	keepLoop := true
	for keepLoop {
		// Checks if a current eviction is already running
		if IsOperationRunning(EventsOperation) {
			log.Debug.Println("Another task executed by the user is performing operations on Events. Skipping.")
		} else {
			taskPostEvents(tid, eventsRecorderAdapter, eventsStorageAdapter, int64(eventsBulkSize))
		}

		select {
		case msg := <-eventsIncoming:
			if msg == "STOP" {
				log.Debug.Println("Stopping task: post_events")
				keepLoop = false
			}
		case <-time.After(time.Duration(eventsRefreshRate) * time.Second):
		}
	}
	wg.Done()
}

// EventsFlush Task to flush cached events.
func EventsFlush(
	eventsRecorderAdapter recorder.EventsRecorder,
	eventsStorageAdapter *redis.EventStorageAdapter,
	size *int64,
) error {
	if RequestOperation(EventsOperation) {
		defer FinishOperation(EventsOperation)
	} else {
		log.Debug.Println("Cannot execute flush. Another operation is performing operations on Events.")
		return errors.New("Cannot execute flush. Another operation is performing operations on Events")
	}
	elementsToFlush := api.MaxSizeToFlush
	if size != nil {
		elementsToFlush = *size
	}

	for elementsToFlush > 0 && eventsStorageAdapter.Size() > 0 {
		maxSize := api.DefaultSize
		if elementsToFlush < api.DefaultSize {
			maxSize = elementsToFlush
		}
		taskPostEvents(0, eventsRecorderAdapter, eventsStorageAdapter, maxSize)
		elementsToFlush = elementsToFlush - api.DefaultSize
	}
	return nil
}
