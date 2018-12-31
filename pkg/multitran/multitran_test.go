package multitran

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWord(t *testing.T) {
	res, err := GetWord("put out")
	assert.NoError(t, err)

	fmt.Println(res.String(100))

}
