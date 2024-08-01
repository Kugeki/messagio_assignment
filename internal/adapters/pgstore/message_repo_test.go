package pgstore

import (
	"context"
	"messagio_assignment/internal/domain"
	"messagio_assignment/internal/domain/message"
)

func (su *PGStoreTestSuite) MsgRepo() *MessageRepoPG {
	return su.store.Message()
}

func (su *PGStoreTestSuite) TestMessageRepo() {
	su.Run("get repo", func() {
		repo := su.store.Message()
		su.Require().NotNil(repo)
	})

	su.Run("get message that not exists", func() {
		id := 1
		_, err := su.MsgRepo().GetByID(context.Background(), id)
		su.ErrorIs(err, domain.ErrNotFound)

		var wantErr *message.ErrorWithID
		su.ErrorAs(err, &wantErr)
	})

	su.Run("create, get message", func() {
		tcases := []struct {
			Name string
			Msg  *message.Message
		}{
			{
				Name: "simple",
				Msg: &message.Message{
					ID:        0,
					Content:   "this is content that is string",
					Processed: false,
				},
			},
			{
				Name: "without default fields",
				Msg: &message.Message{
					ID:        1,
					Content:   "default fields not exist, you know",
					Processed: true,
				},
			},
		}

		for _, tc := range tcases {
			su.Run(tc.Name, func() {
				err := su.MsgRepo().Create(context.Background(), tc.Msg)
				su.NoError(err)

				gotMsg, err := su.MsgRepo().GetByID(context.Background(), tc.Msg.ID)
				su.NoError(err)
				su.Equal(tc.Msg, gotMsg)
			})
		}
	})

	su.Run("stats", func() {
		notProcessed := func() *message.Message {
			return &message.Message{ID: 0, Content: "not processed", Processed: false}
		}
		processed := func() *message.Message {
			return &message.Message{ID: 0, Content: "processed", Processed: true}
		}

		su.Run("with create", func() {
			tcases := []struct {
				Name      string
				Messages  []*message.Message
				WantStats *message.Stats
			}{
				{
					Name:     "0 all, 0 processed",
					Messages: []*message.Message{},
					WantStats: &message.Stats{
						All:       0,
						Processed: 0,
					},
				},
				{
					Name:     "1 all, 0 processed",
					Messages: []*message.Message{notProcessed()},
					WantStats: &message.Stats{
						All:       1,
						Processed: 0,
					},
				},
				{
					Name:     "1 all, 1 processed",
					Messages: []*message.Message{processed()},
					WantStats: &message.Stats{
						All:       1,
						Processed: 1,
					},
				},
				{
					Name:     "2 all, 1 processed",
					Messages: []*message.Message{processed(), notProcessed()},
					WantStats: &message.Stats{
						All:       2,
						Processed: 1,
					},
				},
				{
					Name: "7 all, 3 processed",
					Messages: []*message.Message{processed(), notProcessed(), processed(),
						processed(), notProcessed(), notProcessed(), notProcessed()},
					WantStats: &message.Stats{
						All:       7,
						Processed: 3,
					},
				},
			}

			for _, tc := range tcases {
				su.Run(tc.Name, func() {
					for _, m := range tc.Messages {
						err := su.MsgRepo().Create(context.Background(), m)
						su.NoError(err)
					}

					gotStats, err := su.MsgRepo().GetStats(context.Background())
					su.NoError(err)
					su.Equal(tc.WantStats, gotStats)
				})
			}
		})

		su.Run("with update", func() {
			err := su.MsgRepo().UpdateProcessed(context.Background(), &message.Message{
				ID:        1,
				Content:   "123",
				Processed: true,
			})
			su.ErrorIs(err, domain.ErrNotFound)

			messages := []*message.Message{
				notProcessed(), notProcessed(),
				notProcessed(), notProcessed(),
			}

			for _, msg := range messages {
				err := su.MsgRepo().Create(context.Background(), msg)
				su.NoError(err)
			}

			processedCount := 0

			checkProcessedUpdate := func(i int, processed bool) {
				if messages[i].Processed != processed {
					if processed {
						processedCount++
					} else {
						processedCount--
					}
				}

				messages[i].Processed = processed
				err := su.MsgRepo().UpdateProcessed(context.Background(), messages[i])
				su.NoError(err)

				gotMsg, err := su.MsgRepo().GetByID(context.Background(), messages[i].ID)
				su.NoError(err)
				su.NotNil(gotMsg)
				su.Equal(messages[i], gotMsg)

				wantStats := &message.Stats{
					All:       len(messages),
					Processed: processedCount,
				}
				gotStats, err := su.MsgRepo().GetStats(context.Background())
				su.NoError(err)
				su.NotNil(gotMsg)
				su.Equal(wantStats, gotStats)
			}

			for i := range messages {
				checkProcessedUpdate(i, true)
			}

			for i := range messages {
				checkProcessedUpdate(i, false)
			}
		})
	})
}
