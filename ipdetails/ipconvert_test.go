package ipdetails

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIpToFloat(t *testing.T) {
	result := IpToFloat("1.0.103.211")
	fmt.Println(result)
	require.NotEqual(t, 0, result)
}
