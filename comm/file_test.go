package comm

import (
	"fmt"
	"testing"
)

func TestStakeMark(t *testing.T) {
	var distance float64 = TransformStakeMarkToDistance("K sb 193     + a234.7")
	fmt.Println(distance)
	distance = TransformStakeMarkToDistance("K193     + a234.7")
	fmt.Println(distance)
	distance = TransformStakeMarkToDistance("K193  dd   + a234.7")
	fmt.Println(distance)
	distance = TransformStakeMarkToDistance("  K193  dd   + a234z.7")
	fmt.Println(distance)
}
