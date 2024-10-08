package date

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/slog"
	"github.com/sfomuseum/go-activitypub/stats"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	blocks_db, err := database.NewBlocksDatabase(ctx, opts.BlocksDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer blocks_db.Close(ctx)

	boosts_db, err := database.NewBoostsDatabase(ctx, opts.BoostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer boosts_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	following_db, err := database.NewFollowingDatabase(ctx, opts.FollowingDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer following_db.Close(ctx)

	likes_db, err := database.NewLikesDatabase(ctx, opts.LikesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate likes database, %w", err)
	}

	defer likes_db.Close(ctx)

	messages_db, err := database.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate messages database, %w", err)
	}

	defer messages_db.Close(ctx)

	notes_db, err := database.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate notes database, %w", err)
	}

	defer notes_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	counts_opts := &stats.CountsForDateOptions{
		Date:               opts.Date,
		Location:           opts.Location,
		AccountsDatabase:   accounts_db,
		BlocksDatabase:     blocks_db,
		BoostsDatabase:     boosts_db,
		DeliveriesDatabase: deliveries_db,
		FollowersDatabase:  followers_db,
		FollowingDatabase:  following_db,
		LikesDatabase:      likes_db,
		MessagesDatabase:   messages_db,
		NotesDatabase:      notes_db,
		PostsDatabase:      posts_db,
	}

	counts, err := stats.CountsForDate(ctx, counts_opts)

	if err != nil {
		return fmt.Errorf("Failed to derive counts for date, %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(counts)

	if err != nil {
		return fmt.Errorf("Failed to encode counts, %w", err)
	}

	return nil
}
