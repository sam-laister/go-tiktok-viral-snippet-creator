package repository

import (
	"context"
	"time"

	"github.com/sam-laister/tiktok-creator/ent"
	"github.com/sam-laister/tiktok-creator/ent/clip"
)

type ClipRepository struct {
	client *ent.Client
}

func NewClipRepository(
	client *ent.Client,
) *ClipRepository {
	return &ClipRepository{
		client: client,
	}
}

func (r *ClipRepository) Create(ctx context.Context, clip *ent.Clip) (*ent.Clip, error) {
	c, err := r.client.Clip.
		Create().
		SetHash(clip.Hash).
		SetVideoPath(clip.VideoPath).
		SetAudioPath(clip.AudioPath).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClipRepository) GetClipByID(ctx context.Context, id int) (*ent.Clip, error) {
	c, err := r.client.Clip.
		Query().
		Where(clip.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClipRepository) Update(ctx context.Context, clip *ent.Clip) (*ent.Clip, error) {
	c, err := r.client.Clip.
		UpdateOne(clip).
		SetNillableGenCaptionsPath(clip.GenCaptionsPath).
		SetVideoPath(clip.VideoPath).
		SetAudioPath(clip.AudioPath).
		SetNillableGenRawVideoPath(clip.GenRawVideoPath).
		SetNillableGenTrimmedVideoPath(clip.GenTrimmedVideoPath).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClipRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Clip.
		UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClipRepository) GetOrCreateWithHash(
	ctx context.Context,
	hash, audioPath, videoPath string,
) (*ent.Clip, error) {
	c, err := r.client.Clip.
		Query().
		Where(clip.Hash(hash)).
		Only(ctx)

	switch {
	case err == nil:
		return c, nil
	case !ent.IsNotFound(err):
		return nil, err
	}

	return r.Create(ctx, &ent.Clip{
		AudioPath: audioPath,
		VideoPath: videoPath,
		Hash:      hash,
	})
}
