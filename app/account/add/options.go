package add

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	AliasesDatabaseURI    string
	PropertiesDatabaseURI string
	AccountId             int64
	AccountName           string
	Aliases               []string
	AccountType           string
	Discoverable          bool
	DisplayName           string
	Blurb                 string
	URL                   string
	PublicKeyURI          string
	PrivateKeyURI         string
	AccountIconURI        string
	AllowRemoteIconURI    bool
	EmbedIconURI          bool
	Properties            map[string]string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	properties_map := make(map[string]string)

	for _, kv := range properties_kv {
		properties_map[kv.Key()] = kv.Value().(string)
	}

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		AliasesDatabaseURI:    aliases_database_uri,
		PropertiesDatabaseURI: properties_database_uri,
		AccountId:             account_id,
		AccountName:           account_name,
		Aliases:               aliases_list,
		AccountType:           account_type,
		AccountIconURI:        account_icon_uri,
		AllowRemoteIconURI:    allow_remote_icon_uri,
		EmbedIconURI:          embed_icon_uri,
		Discoverable:          discoverable,
		DisplayName:           display_name,
		Blurb:                 blurb,
		URL:                   account_url,
		PublicKeyURI:          public_key_uri,
		PrivateKeyURI:         private_key_uri,
		Properties:            properties_map,
	}

	return opts, nil
}
