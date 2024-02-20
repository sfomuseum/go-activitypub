package id

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var snowflake_node *snowflake.Node

var setupSnowflakeOnce sync.Once
var setupSnowflakeErr error

func setupSnowflake() {

	node, err := snowflake.NewNode(1)

	if err != nil {
		setupSnowflakeErr = fmt.Errorf("Failed to create snowflake node, %w", err)
		return
	}

	snowflake_node = node
}

func NewId() (int64, error) {

	setupSnowflakeOnce.Do(setupSnowflake)

	if setupSnowflakeErr != nil {
		return 0, setupSnowflakeErr
	}

	id := snowflake_node.Generate()
	return id.Int64(), nil
}

func NewUUID() string {
	guid := uuid.New()
	return guid.String()
}
