package event_test

import (
	"testing"
	"time"

	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/stretchr/testify/assert"
)

func TestNewObservable(t *testing.T) {
	t.Run("observable life cycle", func(t *testing.T) {
		// Создаем наблюдаемый
		observable := event.NewObservable()

		assert.NotNil(t, observable)

		// Канал для отслеживания отработки подписки
		subCh := make(chan *event.Event)
		defer close(subCh)

		// Подписываемся на событие
		subscriber := observable.Subscribe(func(e *event.Event) {
			if e.Name == "test" {
				subCh <- e
			}
		})

		// Испускаем событие
		observable.Next(&event.Event{
			Name: "test",
		})

		// Ждем, обработки события
		timeout := time.After(2 * time.Second)
		select {
		case <-subCh:
			break
		case <-timeout:
			t.Error("timed out waiting for subscriber")
		}

		// Отписываемся от события
		subscriber.Unsubscribe()

		// Испускаем событие
		observable.Next(&event.Event{
			Name: "test",
		})

		// Ждем, обработки события
		timeout = time.After(2 * time.Second)
		select {
		case <-subCh:
			t.Error("timed out waiting for subscriber")
		case <-timeout:
			break
		}
	})
}
