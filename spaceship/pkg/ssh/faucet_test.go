package ssh

import (
	"context"
	"fmt"
	"testing"
)

func Test_faucetBinary(t *testing.T) {
	got, err := fetchFaucetBinary(context.Background(), "darwin_amd64")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(got)
}
