package pgstore

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"messagio_assignment/internal/domain"
	"messagio_assignment/internal/domain/message"
)

type MessageRepoPG struct {
	db *pgxpool.Pool
}

func NewMessageRepoPG(db *pgxpool.Pool) *MessageRepoPG {
	return &MessageRepoPG{db: db}
}

func (r *MessageRepoPG) Create(ctx context.Context, msg *message.Message) error {
	q := "insert into messages(content, processed) values($1, $2) returning id"

	var id int
	err := r.db.QueryRow(ctx, q, msg.Content, msg.Processed).Scan(&id)
	if err != nil {
		return &message.Error{Err: ErrCreateIntoDomain(err)}
	}

	msg.ID = id
	return nil
}

func (r *MessageRepoPG) GetByID(ctx context.Context, id int) (*message.Message, error) {
	q := "select m.id, m.content, m.processed from messages as m where m.id = $1"

	var msg message.Message
	err := r.db.QueryRow(ctx, q, id).Scan(&msg.ID, &msg.Content, &msg.Processed)
	if err != nil {
		return nil, &message.ErrorWithID{ID: id, Err: ErrGetIntoDomain(err)}
	}

	return &msg, nil
}

func (r *MessageRepoPG) GetStats(ctx context.Context) (*message.Stats, error) {
	q := `select (select count(*) from messages as m), 
       (select count(*) from messages as m where m.processed = true)`

	var stats message.Stats
	err := r.db.QueryRow(ctx, q).Scan(&stats.All, &stats.Processed)
	if err != nil {
		return nil, &message.StatsError{Err: err}
	}

	return &stats, nil
}

func (r *MessageRepoPG) UpdateProcessed(ctx context.Context, msg *message.Message) error {
	q := "update messages as m set processed = $1 where m.id = $2"

	ct, err := r.db.Exec(ctx, q, msg.Processed, msg.ID)
	if err != nil {
		return &message.ErrorWithID{ID: msg.ID, Err: err}
	}
	if ct.RowsAffected() == 0 {
		return &message.ErrorWithID{ID: msg.ID, Err: domain.ErrNotFound}
	}

	return nil
}
