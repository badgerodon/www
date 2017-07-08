package rbsa

import (
	//. "github.com/badgerodon/lalg"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
/*
	alg := New()
	alg.AddIndex("i1", Vector{1,1.1,3})
	alg.AddIndex("i2", Vector{3,2,1})
	alg.AddIndex("i3", Vector{1,1,1})
	solution, err := alg.Run(Vector{1,2,1})
	
	fmt.Print(solution, err)
*/
	solution, err := Analyze("IVV")
	fmt.Println(solution, err)
}