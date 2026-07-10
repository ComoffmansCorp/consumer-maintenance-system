package photo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

var (
	ErrPhotoNotFound   = errors.New("photo not found")
	ErrInvalidTarget   = errors.New("exactly one of inspectionActId or replacementActId is required")
	ErrFileRequired    = errors.New("file is required")
	ErrFileTooLarge    = errors.New("file exceeds the maximum allowed size")
	ErrUnsupportedType = errors.New("only JPEG and PNG images are supported")
)

// MaxFileSize matches the legacy Spring upload limit (20MB).
const MaxFileSize = 20 << 20

var allowedContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

type Service struct {
	repo      *Repository
	acts      ActStore
	uploadDir string
}

func NewService(repo *Repository, acts ActStore, uploadDir string) *Service {
	return &Service{repo: repo, acts: acts, uploadDir: uploadDir}
}

type UploadInput struct {
	TenantID         int64
	UploadedBy       int64
	InspectionActID  *int64
	ReplacementActID *int64
	Note             string
	OriginalFilename string
	ContentType      string
	Size             int64
	Content          io.Reader
}

func (s *Service) Upload(ctx context.Context, in UploadInput) (DTO, error) {
	if (in.InspectionActID == nil) == (in.ReplacementActID == nil) {
		return DTO{}, ErrInvalidTarget
	}
	if in.InspectionActID != nil {
		if err := s.acts.EnsureInspectionAct(ctx, *in.InspectionActID, in.TenantID); err != nil {
			return DTO{}, err
		}
	} else {
		if err := s.acts.EnsureReplacementAct(ctx, *in.ReplacementActID, in.TenantID); err != nil {
			return DTO{}, err
		}
	}

	ext, ok := allowedContentTypes[in.ContentType]
	if !ok {
		return DTO{}, ErrUnsupportedType
	}
	if in.Size <= 0 {
		return DTO{}, ErrFileRequired
	}
	if in.Size > MaxFileSize {
		return DTO{}, ErrFileTooLarge
	}

	storedName, err := s.store(in.TenantID, ext, in.Content)
	if err != nil {
		return DTO{}, err
	}

	uploadedBy := in.UploadedBy
	p, err := s.repo.Create(ctx, CreateParams{
		Filename:         storedName,
		Note:             in.Note,
		TenantID:         in.TenantID,
		InspectionActID:  in.InspectionActID,
		ReplacementActID: in.ReplacementActID,
		OriginalFilename: sanitizeFilename(in.OriginalFilename),
		ContentType:      in.ContentType,
		SizeBytes:        in.Size,
		UploadedBy:       &uploadedBy,
	})
	if err != nil {
		_ = os.Remove(filepath.Join(s.tenantDir(in.TenantID), storedName))
		return DTO{}, fmt.Errorf("save photo: %w", err)
	}
	return ToDTO(p), nil
}

func (s *Service) ListInspection(ctx context.Context, actID, tenantID int64) ([]DTO, error) {
	if err := s.acts.EnsureInspectionAct(ctx, actID, tenantID); err != nil {
		return nil, err
	}
	photos, err := s.repo.ListByInspectionAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return toDTOs(photos), nil
}

func (s *Service) ListReplacement(ctx context.Context, actID, tenantID int64) ([]DTO, error) {
	if err := s.acts.EnsureReplacementAct(ctx, actID, tenantID); err != nil {
		return nil, err
	}
	photos, err := s.repo.ListByReplacementAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return toDTOs(photos), nil
}

// GetFile returns the photo record and its absolute path on disk for
// download handlers.
func (s *Service) GetFile(ctx context.Context, id, tenantID int64) (Photo, string, error) {
	p, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Photo{}, "", ErrPhotoNotFound
		}
		return Photo{}, "", err
	}
	return p, s.filePath(p), nil
}

func (s *Service) Delete(ctx context.Context, id, tenantID int64) error {
	p, err := s.repo.Delete(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPhotoNotFound
		}
		return fmt.Errorf("delete photo: %w", err)
	}
	_ = os.Remove(s.filePath(p))
	return nil
}

// --- Ports exposed to the act domain ---

func (s *Service) ListForInspectionAct(ctx context.Context, actID int64) ([]Summary, error) {
	photos, err := s.repo.ListByInspectionAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return s.toSummaries(photos), nil
}

func (s *Service) ListForReplacementAct(ctx context.Context, actID int64) ([]Summary, error) {
	photos, err := s.repo.ListByReplacementAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return s.toSummaries(photos), nil
}

// --- helpers ---

func (s *Service) toSummaries(photos []Photo) []Summary {
	out := make([]Summary, 0, len(photos))
	for _, p := range photos {
		out = append(out, Summary{
			ID:               p.ID,
			OriginalFilename: p.OriginalFilename,
			Note:             p.Note,
			FilePath:         s.filePath(p),
			CreatedAt:        p.CreatedAt,
		})
	}
	return out
}

func (s *Service) tenantDir(tenantID int64) string {
	return filepath.Join(s.uploadDir, strconv.FormatInt(tenantID, 10))
}

func (s *Service) filePath(p Photo) string {
	return filepath.Join(s.tenantDir(p.TenantID), p.Filename)
}

func (s *Service) store(tenantID int64, ext string, r io.Reader) (string, error) {
	dir := s.tenantDir(tenantID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}
	name, err := randomFilename(ext)
	if err != nil {
		return "", err
	}
	dest := filepath.Join(dir, name)
	f, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = f.Close() }()
	// Defensive re-check of size in case the declared Content-Length lied.
	if _, err := io.Copy(f, io.LimitReader(r, MaxFileSize+1)); err != nil {
		_ = os.Remove(dest)
		return "", fmt.Errorf("write file: %w", err)
	}
	return name, nil
}

func randomFilename(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate filename: %w", err)
	}
	return hex.EncodeToString(buf) + ext, nil
}

func sanitizeFilename(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "photo"
	}
	return name
}

func toDTOs(photos []Photo) []DTO {
	out := make([]DTO, 0, len(photos))
	for _, p := range photos {
		out = append(out, ToDTO(p))
	}
	return out
}
