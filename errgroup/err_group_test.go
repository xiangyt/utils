package errgroup

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestErrGroup(t *testing.T) {

	group, _ := WithContext(context.Background())
	group.SetLimit(10)
	for i := 0; i < 100; i++ {
		i := i
		group.Go(func() error {
			time.Sleep(time.Second * 1)
			fmt.Println(i)
			if i == 10 {
				panic("1")
			}
			return nil
		})
	}

	fmt.Println("start wait")
	if err := group.Wait(); err != nil {
		t.Log(err)
	}

}
