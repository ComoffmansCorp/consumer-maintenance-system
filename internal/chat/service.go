package chat

import (
	"context"
	"errors"
	"strings"
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
}

func NewService(repo *Repository, requests RequestPort) *Service {
	return &Service{repo: repo, requests: requests}
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
	return ToMessageDTO(m), nil
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
