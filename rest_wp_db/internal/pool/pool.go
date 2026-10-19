package pool

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	errs "github.com/arlengur/go_assessment/internal/errors"
	"github.com/arlengur/go_assessment/internal/model"
	"github.com/arlengur/go_assessment/internal/repository"
)

type workerPool struct {
	repo         repository.Repository
	numWorkers   int
	taskChan     chan string
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.Mutex
	startOnce    sync.Once
	activeTracks map[string]context.CancelFunc
}

func NewWorkerPool(repo repository.Repository, numWorkers int, chanSize int) *workerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &workerPool{
		repo:         repo,
		numWorkers:   numWorkers,
		taskChan:     make(chan string, chanSize),
		ctx:          ctx,
		cancel:       cancel,
		activeTracks: make(map[string]context.CancelFunc),
	}
}

func (p *workerPool) Start() {
	p.startOnce.Do(func() {
		for range p.numWorkers {
			p.wg.Add(1)
			go p.worker()
		}
		p.wg.Add(1)
		go p.poller()
	})
}

func (p *workerPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case id, ok := <-p.taskChan:
			if !ok {
				return
			}
			p.processTask(id)
		}
	}
}

func (p *workerPool) processTask(id string) {
	log.Printf("Processing task: %s", id)
	err := p.repo.UpdatePendingTask(p.ctx, id)
	if err != nil {
		if dbErr, ok := errors.AsType[*errs.DbError](err); ok {
			log.Fatalf("Critical error updatePendingTask: %v\n", dbErr.Error())
		}
		log.Printf("UpdatePendingTask: failed to update task status: %v", err)
		return
	}

	taskCtx, taskCancel := context.WithCancel(p.ctx)
	p.mu.Lock()
	p.activeTracks[id] = taskCancel
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.activeTracks, id)
		p.mu.Unlock()
		taskCancel()
	}()

	// Simulating processing work logic
	err = p.executePayloadWork(taskCtx, id)

	finalStatus := model.StatusCompleted
	var errStr string
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return // State updated directly via CancelTask
		}
		finalStatus = model.StatusFailed
		errStr = err.Error()
	}

	err = p.repo.UpdateProcessingTask(p.ctx, finalStatus, errStr, id)
	if err != nil {
		if dbErr, ok := errors.AsType[*errs.DbError](err); ok {
			log.Fatalf("Critical error updateProcessingTask: %v\n", dbErr.Error())
		}
		log.Printf("UpdateProcessingTask: failed to update task status: %v", err)
		return
	}
	log.Printf("Task completed: %s", id)
}

func (p *workerPool) submit(id string) {
	select {
	case p.taskChan <- id:
	case <-p.ctx.Done():
	}
}

func (p *workerPool) CancelTask(id string) error {
	p.mu.Lock()
	cancelFn, tracking := p.activeTracks[id]
	if tracking && cancelFn != nil {
		cancelFn() // Interrupt the runtime execution
	}
	p.mu.Unlock()

	err := p.repo.CancelTask(p.ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *workerPool) poller() {
	defer p.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			ids, err := p.repo.GetPendingTasks(p.ctx)
			if err != nil {
				if dbErr, ok := errors.AsType[*errs.DbError](err); ok {
					log.Fatalf("Critical error polling for tasks: %v\n", dbErr.Error())
				}
				log.Printf("Error polling for tasks: %v", err)
				continue
			}
			for _, id := range ids {
				p.submit(id)
			}
		}
	}
}

func (p *workerPool) executePayloadWork(ctx context.Context, id string) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return p.ctx.Err()
	}
}

func (p *workerPool) Stop() {
	p.cancel()
	close(p.taskChan)
	p.wg.Wait()
}
