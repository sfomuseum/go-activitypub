package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Block struct {
	Id           int64  `json:"id"`
	AccountId    int64  `json:"account_id"`
	Name         string `json:"name"`
	Host         string `json:"host"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func NewBlock(ctx context.Context, account_id int64, block_host string, block_name string) (*Block, error) {

	block_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new block ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	b := &Block{
		Id:           block_id,
		AccountId:    account_id,
		Host:         block_host,
		Name:         block_name,
		Created:      ts,
		LastModified: ts,
	}

	return b, nil
}

func IsBlockedByAccount(ctx context.Context, db BlocksDatabase, account_id int64, host string, name string) (bool, error) {

	_, err := db.GetBlockWithAccountIdAndAddress(ctx, account_id, host, name)

	if err == nil {
		return true, nil
	}

	if err != ErrNotFound {
		return false, fmt.Errorf("Failed to retrieve block with account and address, %w", err)
	}

	if name == "*" {
		return false, nil
	}

	return IsBlockedByAccount(ctx, db, account_id, host, "*")
}
