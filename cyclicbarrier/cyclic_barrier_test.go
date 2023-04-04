package cyclicbarrier

import (
	"context"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

type H2O struct {
	semaH *semaphore.Weighted // 氢原子的信号量
	semaO *semaphore.Weighted // 氧原子的信号量
	cb    CyclicBarrier       // 循环栅栏，用来控制合成}
}

func (h *H2O) hydrogen(releaseHydrogen func()) {

	h.semaH.Acquire(context.Background(), 1)

	releaseHydrogen()
	h.cb.Await(context.Background())
	h.semaH.Release(1)
}

func (h *H2O) oxygen(releaseOxygen func()) {
	h.semaO.Acquire(context.Background(), 1)

	releaseOxygen()
	h.cb.Await(context.Background())
	h.semaO.Release(1)
}

func TestCyclicBarrier(t *testing.T) {

	// 有一个名叫大自然的搬运工的工厂，生产一种叫做一氧化二氢的神秘液体。
	// 这种液体的分子是由一个氧原子和两个氢原子组成的，也就是水。
	// 这个工厂有多条生产线，每条生产线负责生产氧原子或者是氢原子，每条生产线由一个 goroutine 负责。
	// 这些生产线会通过一个栅栏，只有一个氧原子生产线和两个氢原子生产线都准备好，才能生成出一个水分子，
	// 否则所有的生产线都会处于等待状态。也就是说，一个水分子必须由三个不同的生产线提供原子，
	// 而且水分子是一个一个按照顺序产生的，每生产一个水分子，就会打印出 HHO、HOH、OHH 三种形式的其中一种。
	// HHH、OOH、OHO、HOO、OOO 都是不允许的。生产线中氢原子的生产线为 2N 条，氧原子的生产线为 N 条。
	h2o := &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		cb:    New(3),
	}

	//用来存放水分子结果的channel
	var ch chan string

	releaseHydrogen := func() { ch <- "H" }
	releaseOxygen := func() { ch <- "O" }
	// 300个原子，300个goroutine,每个goroutine并发的产生一个原子
	var N = 100
	ch = make(chan string, N*3)

	// 用来等待所有的goroutine完成
	var wg sync.WaitGroup
	wg.Add(N * 3)

	for i := 0; i < N*2; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}

	//等待所有的goroutine执行完
	wg.Wait()

	// 结果中肯定是300个原子
	if len(ch) != N*3 {
		t.Fatalf("expect %d atom but got %d", N*3, len(ch))
	}

	// 每三个原子一组，分别进行检查。要求这一组原子中必须包含两个氢原子和一个氧原子，这样才能正确组成一个水分子。
	var arr = make([]string, 3)
	for i := 0; i < N; i++ {
		arr[0] = <-ch
		arr[1] = <-ch
		arr[2] = <-ch
		sort.Strings(arr)

		water := strings.Join(arr, "")
		if water != "HHO" {
			t.Fatalf("expect a water molecule but got %s", water)
		}
	}

}
