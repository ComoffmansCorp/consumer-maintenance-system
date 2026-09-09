package chat

import (
	"context"
	"errors"
	"strings"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
)

const defaultPageSize = 50

var (
	ErrRequestNotFound = errors.New("request not found")
	ErrForbidden       = errors.New("only a participant of the request can access its chat")
	ErrNotAssignedYet  = errors.New("chat is only available once the request has been assigned to a master")
	ErrTextRequired    = errors.New("message text is required")
)

type Service struct {
	repo     *Repository
	requests RequestPort
	events   *broker.Bus
}

func NewService(repo *Repository, requests RequestPort, events *broker.Bus) *Service {
	return &Service{repo: repo, requests: requests, events: events}
}

func (s *Service) SendMessage(ctx context.Context, requestID, senderID int64, req SendMessageRequest) (MessageDTO, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return MessageDTO{}, ErrTextRequired
	}
	if err := s.authorize(ctx, requestID, senderID); err != nil {
		return MessageDTO{}, err
	}

	m, err := s.repo.Create(ctx, requestID, senderID, text)
	if err != nil {
		return MessageDTO{}, err
	}
	dto := ToMessageDTO(m)

	s.events.Publish(ctx, broker.Event{
		Type:    broker.EventChatMessageCreated,
		Payload: map[string]any{"request_id": requestID, "message": dto},
	})

	return dto, nil
}

func (s *Service) ListMessages(ctx context.Context, requestID, callerID, sinceID int64) ([]MessageDTO, error) {
	if err := s.authorize(ctx, requestID, callerID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListSince(ctx, requestID, sinceID, defaultPageSize)
	if err != nil {
		return nil, err
	}
	out := make([]MessageDTO, 0, len(rows))
	for _, m := range rows {
		out = append(out, ToMessageDTO(m))
	}
	return out, nil
}

// MarkRead marks every message the caller didn't send in this thread as
// read -- called when a client/master actually opens the chat, not on
// every poll, so read receipts reflect "seen the thread" rather than "the
// client happened to have polled this endpoint".
func (s *Service) MarkRead(ctx context.Context, requestID, readerID int64) error {
	if err := s.authorize(ctx, requestID, readerID); err != nil {
		return err
	}
	return s.repo.MarkRead(ctx, requestID, readerID)
}

// CountUnreadConversations returns how many of the caller's request threads
// have at least one unread message from the other participant -- the
// number a header badge shows, not a raw message count.
func (s *Service) CountUnreadConversations(ctx context.Context, userID int64) (int64, error) {
	requestIDs, err := s.requests.ListParticipantRequestIDs(ctx, userID)
	if err != nil {
		return 0, err
	}
	if len(requestIDs) == 0 {
		return 0, nil
	}
	return s.repo.CountUnreadRequestIDs(ctx, requestIDs, userID)
}

// authorize enforces both rules the chat thread depends on: the caller must
// be a participant, and the request must be at least ASSIGNED (master_id
// IS NOT NULL) -- there's no one to chat with on a still-OPEN request.
func (s *Service) authorize(ctx context.Context, requestID, callerID int64) error {
	clientID, masterID, found, err := s.requests.GetParticipants(ctx, requestID)
	if err != nil {
		return err
	}
	if !found {
		return ErrRequestNotFound
	}
	if masterID == nil {
		return ErrNotAssignedYet
	}
	if callerID != clientID && callerID != *masterID {
		return ErrForbidden
	}
	return nil
}
