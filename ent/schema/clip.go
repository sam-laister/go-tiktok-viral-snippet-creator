package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Clip holds the schema definition for the Clip entity.
type Clip struct {
	ent.Schema
}

// Fields of the Clip.
func (Clip) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash"),
		field.String("audio_path"),
		field.String("video_path"),
		field.String("gen_captions_path").
			Optional().
			Nillable(),
		field.String("gen_raw_video_path").
			Optional().
			Nillable(),
		field.String("gen_trimmed_video_path").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the Clip.
func (Clip) Edges() []ent.Edge {
	return nil
}
