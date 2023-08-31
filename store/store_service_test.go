package store

import(
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStore = &StorageService{}

func init(){
	testStoreService = InitializeStore()
}

func testStoreInit(t *testing.T){
	assert.True(t, testStoreService.redisClient != nil)
}

func testInsertionAndRetreival(t *testing.T){
	// TODO: insertion and retrieve test cases
}

