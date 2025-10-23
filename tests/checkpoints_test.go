package examples

import (
	"context"
	"testing"

	sui "github.com/0xdraco/sui-go-sdk/grpc"
)

func TestGetCheckpoint(t *testing.T) {
	ctx := context.Background()
	client, err := sui.NewMainnetClient(ctx)
	requireNoError(t, err, "NewMainnetClient")
	t.Cleanup(func() {
		client.Close()
	})

	const (
		sequence uint64 = 201477601
		digest          = "CEpRBP5xcdBZYG8q1sxEkm2vLyDYaa8Rf3fAkC3zhZ9j"
	)

	checkpointBySeq, err := client.GetCheckpointBySequence(ctx, sequence, nil)
	requireNoError(t, err, "GetCheckpointBySequence")
	requireNotNil(t, checkpointBySeq, "GetCheckpointBySequence")
	requireEqual(t, checkpointBySeq.GetDigest(), digest, "checkpoint digest")

	checkpointByDigest, err := client.GetCheckpointByDigest(ctx, digest, nil)
	requireNoError(t, err, "GetCheckpointByDigest")
	requireNotNil(t, checkpointByDigest, "GetCheckpointByDigest")
	requireEqual(t, checkpointByDigest.GetSequenceNumber(), sequence, "checkpoint sequence")
}
