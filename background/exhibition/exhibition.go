package exhibition

import (
	"context"
	"fmt"
	"time"
)

func Create(ctx context.Context, index int) error {
	fmt.Println("Creating an exhibition", index)
	time.Sleep(1 * time.Second)
	return nil
}
