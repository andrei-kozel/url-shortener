package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/andrei-kozel/url-shortener/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mng struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Client) *mng {
	return &mng{
		db: client.Database("url_shortener"),
	}
}

func (m *mng) col() *mongo.Collection {
	return m.db.Collection("shortenings")
}

func (m *mng) Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error) {
	const op = "storage.mng.Put"

	shortening.CreatedAt = time.Now().UTC()

	count, err := m.col().CountDocuments(ctx, bson.M{"_id": shortening.Identifier})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", op, model.ErrorIdentifierAlreadyExists)
	}

	_, err = m.col().InsertOne(ctx, mgnShorteningFromModel(shortening))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortening, nil
}

func (m *mng) Get(ctx context.Context, identifier string) (*model.Shortening, error) {
	const op = "storage.mng.Get"

	var shortening mgnShortening
	err := m.col().FindOne(ctx, bson.M{"_id": identifier}).Decode(&shortening)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Shortening{
		Identifier:  shortening.Identifier,
		OriginalUrl: shortening.OriginalUrl,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}, nil
}

func (m *mng) IncrementVisits(ctx context.Context, identifier string) error {
	const op = "storage.mng.IncrementVisits"

	_, err := m.col().UpdateOne(ctx,
		bson.M{"_id": identifier},
		bson.M{"$inc": bson.M{"visits": 1}},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type mgnShortening struct {
	Identifier  string    `bson:"_id"`
	OriginalUrl string    `bson:"original_url"`
	Visits      int64     `bson:"visits"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func mgnShorteningFromModel(shortneing model.Shortening) mgnShortening {
	return mgnShortening{
		Identifier:  shortneing.Identifier,
		OriginalUrl: shortneing.OriginalUrl,
		Visits:      shortneing.Visits,
		CreatedAt:   shortneing.CreatedAt,
		UpdatedAt:   shortneing.UpdatedAt,
	}
}
