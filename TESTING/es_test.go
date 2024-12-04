package TESTING

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"koriebruh/find/conf"
	"testing"
)

func TestEsConnection(t *testing.T) {
	client := conf.ElasticClient()
	assert.NotNil(t, client)

	res, err := client.Info()
	assert.Nil(t, err)
	defer res.Body.Close()

	fmt.Println(res)
	assert.NotNil(t, res)
}

func TestMMK(t *testing.T) {
	fmt.Println("MEMEK")
}
