package closer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type closerFn struct {
	name string
	fn   func(ctx context.Context) error
}

type closer struct {
	mu     sync.Mutex
	once   sync.Once
	funcs  []closerFn
	result error
}

var globalCloser = &closer{}

func Add(name string, fn func(ctx context.Context) error) {
	globalCloser.add(name, fn)
}

func CloseAll(ctx context.Context, log *slog.Logger) error {
	return globalCloser.closeAll(ctx, log)
}

func (c *closer) add(name string, fn func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, closerFn{name, fn})
}

func (c *closer) closeAll(ctx context.Context, log *slog.Logger) error {
	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		log.InfoContext(ctx, "начинаем плавное завершение", "count", len(funcs))

		var errs []error

		// Идём от конца к началу - LIFO, как defer.
		for _, v := range slices.Backward(funcs) {
			f := v

			// Бюджет времени исчерпан — пропускаем оставшиеся ресурсы.
			select {
			case <-ctx.Done():
				log.WarnContext(ctx, "контекст истёк, пропускаем ресурс", "name", f.name)
				errs = append(errs, fmt.Errorf("ресурс %q пропущен: %w", f.name, ctx.Err()))
				continue
			default:
			}

			start := time.Now()
			log.InfoContext(ctx, "закрываем ресурс", "name", f.name)

			// Запускаем закрытие в горутине: если fn завис и игнорирует ctx,
			// select ниже всё равно снимется по ctx.Done() и мы идём дальше.
			// Буферизованный канал гарантирует, что горутина не утечёт навсегда
			// в ожидании получателя — она запишет результат и завершится сама,
			// когда fn в итоге вернётся (или при завершении процесса).
			done := make(chan error, 1)
			go func() {
				done <- f.fn(ctx)
			}()

			select {
			case err := <-done:
				if err != nil {
					log.ErrorContext(ctx, "ошибка при закрытии ресурса",
						"name", f.name,
						"err", err,
						"duration", time.Since(start),
					)
					errs = append(errs, fmt.Errorf("ресурс %q: %w", f.name, err))
				} else {
					log.InfoContext(ctx, "ресурс закрыт", "name", f.name, "duration", time.Since(start))
				}
			case <-ctx.Done():
				log.ErrorContext(ctx, "ресурс завис при закрытии",
					"name", f.name,
					"duration", time.Since(start),
				)
				errs = append(errs, fmt.Errorf("ресурс %q завис при закрытии: %w", f.name, ctx.Err()))

				// Не бросаем зависшую горутину. Дожидаемся её завершения(ошибки) и логируем
				go func(name string, start time.Time, done <-chan error) {
					err := <-done
					if err != nil {
						log.WarnContext(ctx, "ресурс закрылся с опозданием/ошибкой",
							"name", name, "duration", time.Since(start), "err", err)
					} else {
						log.WarnContext(ctx, "ресурс закрылся с опозданием",
							"name", name, "duration", time.Since(start))
					}
				}(f.name, start, done)
			}
		}

		log.InfoContext(ctx, "плавное завершение завершено")

		c.result = errors.Join(errs...)
	})

	return c.result
}
