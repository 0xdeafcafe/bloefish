package app

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

func (a *App) ListConversationsWithInteractions(ctx context.Context, req *conversation.ListConversationsWithInteractionsRequest) (*conversation.ListConversationsWithInteractionsResponse, error) {
	conversations, err := a.ConversationRepository.ListByOwner(ctx, models.Actor{
		Type:       models.ActorType(req.Owner.Type),
		Identifier: req.Owner.Identifier,
	})
	if err != nil {
		return nil, err
	}

	errGroup, egCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(runtime.NumCPU() * 4)

	relics := map[string][]*models.Interaction{}
	var mu sync.Mutex

	for _, convo := range conversations {
		errGroup.Go(func() error {
			interactions, err := a.InteractionRepository.GetAllByConversationID(egCtx, convo.ID)
			if err != nil {
				return err
			}

			for _, interaction := range interactions {
				mu.Lock()
				if relics[convo.ID] == nil {
					relics[convo.ID] = make([]*models.Interaction, 0)
				}
				mu.Unlock()

				mu.Lock()
				relics[convo.ID] = append(relics[convo.ID], interaction)
				mu.Unlock()
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	resp := &conversation.ListConversationsWithInteractionsResponse{
		Conversations: make([]*conversation.ListConversationsWithInteractionsResponseConversation, len(conversations)),
	}
	for i, convo := range conversations {
		resp.Conversations[i] = &conversation.ListConversationsWithInteractionsResponseConversation{
			ID: convo.ID,
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(convo.Owner.Type),
				Identifier: convo.Owner.Identifier,
			},
			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: convo.AIRelayOptions.ProviderID,
				ModelID:    convo.AIRelayOptions.ModelID,
			},

			Title:           convo.Title,
			StreamChannelID: convo.ID,

			Interactions: make([]*conversation.ListConversationsWithInteractionsResponseConversationInteraction, len(relics[convo.ID])),

			CreatedAt: convo.CreatedAt,
			UpdatedAt: convo.UpdatedAt,
			DeletedAt: convo.DeletedAt,
		}

		for j, interaction := range relics[convo.ID] {
			resp.Conversations[i].Interactions[j] = &conversation.ListConversationsWithInteractionsResponseConversationInteraction{
				ID: interaction.ID,
				Owner: &conversation.Actor{
					Type:       conversation.ActorType(interaction.Owner.Type),
					Identifier: interaction.Owner.Identifier,
				},
				MessageContent:  interaction.MessageContent,
				FileIDs:         interaction.FileIDs,
				StreamChannelID: fmt.Sprintf("%s/%s", convo.ID, interaction.ID),
				AIRelayOptions: &conversation.AIRelayOptions{
					ProviderID: interaction.AIRelayOptions.ProviderID,
					ModelID:    interaction.AIRelayOptions.ModelID,
				},

				CreatedAt:   interaction.CreatedAt,
				CompletedAt: interaction.CompletedAt,
				UpdatedAt:   interaction.UpdatedAt,
				DeletedAt:   interaction.DeletedAt,
			}
		}
	}

	return resp, nil
}
