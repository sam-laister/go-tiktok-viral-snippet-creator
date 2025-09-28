package service

import (
	"context"

	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/helper"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/repository"
)

type ClipServiceImpl struct {
	clipRepo *repository.ClipRepository
}

func NewClipServiceImpl(
	clipRepo *repository.ClipRepository,
) *ClipServiceImpl {
	return &ClipServiceImpl{
		clipRepo: clipRepo,
	}
}

func (r *ClipServiceImpl) Create(ctx context.Context, clip *model.ClipDTO) error {
	_, err := r.clipRepo.Create(ctx, helper.DTOToClip(clip))
	return err
}

func (r *ClipServiceImpl) Update(ctx context.Context, clip *model.ClipDTO) error {
	_, err := r.clipRepo.Update(ctx, helper.DTOToClip(clip))
	return err
}

func (r *ClipServiceImpl) GetOrCreateWithHash(
	ctx context.Context,
	hash,
	audioPath,
	videoPath string,
) (*model.ClipDTO, error) {
	clipEntity, err := r.clipRepo.GetOrCreateWithHash(ctx, hash, audioPath, videoPath)
	if err != nil {
		return nil, err
	}
	return helper.ClipToDTO(clipEntity), nil
}
