package queue

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/go-logr/logr"
	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/mongo"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
	"github.com/pkg/errors"
)

const randomFactor = 5

var ErrNoData = errors.New("no data")

var location *time.Location

func init() {
	location, _ = time.LoadLocation("Asia/Shanghai")
}

type Interface interface {
	Sync(item models.QueueItem)
	Dequeue() (string, error)
	Ready() bool
	Run(stopCh <-chan struct{})
}

func New() Interface {
	return &Type{
		buf:    make(chan *models.QueueItem, 1000),
		logger: log.New().WithName("queue"),
	}
}

type Type struct {
	items  items
	index  map[string]*models.QueueItem
	buf    chan *models.QueueItem
	ready  bool
	once   sync.Once
	lock   sync.Mutex
	logger logr.Logger
}

func (t *Type) Sync(item models.QueueItem) {
	if t.ready {
		t.buf <- &item
	}
}

func (t *Type) Dequeue() (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	count := len(t.items)
	var i int
	if count <= randomFactor {
		if count == 0 {
			return "", ErrNoData
		}
		i = rand.Intn(count)
	} else {
		i = count - 1 - rand.Intn(randomFactor)
	}
	now := time.Now().In(location)
	if err := mongo.DB().C(mongo.CollectionBox).UpdateId(t.items[i].Id, bson.M{
		"$push":      bson.M{"Viewed": now},
		"LastViewed": now,
	}); err != nil {
		return "", errors.Wrap(err, "update view record")
	}
	t.items[i].Score = now.UnixNano()
	sort.Sort(t.items)
	return t.items[i].Id, nil
}

func (t *Type) Ready() bool {
	return t.ready
}

func (t *Type) Run(stopCh <-chan struct{}) {
	const (
		resync = 15 * time.Minute
		retry  = 10 * time.Second
	)
	resyncTimer := time.NewTimer(0)
	defer resyncTimer.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-resyncTimer.C:
			t.logger.Info("resync initialized")
			go func() {
				if err := t.resync(); err != nil {
					t.logger.Error(err, "resync")
					resyncTimer.Reset(retry)
				} else {
					t.logger.Info("resync successful")
					resyncTimer.Reset(resync)
				}
			}()
		case item := <-t.buf:
			t.logger.Info("sync item received", "id", item.Id, "score", item.Score)
			func() {
				if item != nil {
					t.lock.Lock()
					defer t.lock.Unlock()
					if existing, ok := t.index[item.Id]; ok {
						existing.Score = item.Score
					} else {
						t.items = append(t.items, item)
						t.index[item.Id] = t.items[len(t.items)-1]
					}
					sort.Sort(t.items)
				}
			}()
		}
	}
}

func (t *Type) resync() error {
	var buf []models.Box
	if err := mongo.DB().C(mongo.CollectionBox).Find(nil).All(&buf); err != nil {
		return errors.Wrap(err, "list Box")
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.items = t.items[:0]
	t.index = make(map[string]*models.QueueItem)
	for i := range buf {
		if buf[i].LastViewed != nil {
			item := models.QueueItem{
				Id:    buf[i].Id,
				Score: buf[i].LastViewed.UnixNano(),
			}
			t.items = append(t.items, &item)
			t.index[item.Id] = &item
		}
	}
	sort.Sort(t.items)
	t.once.Do(func() {
		t.ready = true
	})
	return nil
}

type items []*models.QueueItem

func (it items) Len() int {
	return len(it)
}

func (it items) Less(i, j int) bool {
	return it[i].Score < it[j].Score
}

func (it items) Swap(i, j int) {
	cpy := it[i]
	it[i] = it[j]
	it[j] = cpy
}
