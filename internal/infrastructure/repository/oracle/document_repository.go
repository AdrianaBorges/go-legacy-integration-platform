package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/application/service"
	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/domain/document"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (r *DocumentRepository) Create(
	ctx context.Context,
	doc document.Document,
	idempotencyKey string,
) (document.Document, bool, error) {

	if idempotencyKey != "" {
		existing, found, err := r.findByIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			return document.Document{}, false, err
		}

		if found {
			return existing, true, nil
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return document.Document{}, false,
			fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO DOCUMENT (
			ID,
			NAME,
			CONTENT_TYPE,
			STATUS,
			IDEMPOTENCY_KEY,
			CREATED_AT,
			UPDATED_AT
		) VALUES (
			:1, :2, :3, :4, :5, :6, :7
		)`,
		doc.ID,
		doc.Name,
		doc.ContentType,
		doc.Status,
		nullableString(idempotencyKey),
		doc.CreatedAt,
		doc.UpdatedAt,
	)

	if err != nil {
		return document.Document{}, false,
			fmt.Errorf("erro ao inserir documento: %w", err)
	}

	event := document.Event{
		ID:            doc.ID + "-created",
		DocumentID:    doc.ID,
		Type:          "DOCUMENT_CREATED",
		OccurredAt:    doc.CreatedAt,
		CorrelationID: idempotencyKey,
	}

	if err := insertEvent(ctx, tx, event); err != nil {
		return document.Document{}, false, err
	}

	if err := tx.Commit(); err != nil {
		return document.Document{}, false,
			fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return doc, false, nil
}

func (r *DocumentRepository) Get(
	ctx context.Context,
	id string,
) (document.Document, error) {

	var doc document.Document

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
			ID,
			NAME,
			CONTENT_TYPE,
			STATUS,
			CREATED_AT,
			UPDATED_AT
		FROM DOCUMENT
		WHERE ID = :1`,
		id,
	).Scan(
		&doc.ID,
		&doc.Name,
		&doc.ContentType,
		&doc.Status,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return document.Document{}, service.ErrNotFound
	}

	if err != nil {
		return document.Document{},
			fmt.Errorf("erro ao consultar documento: %w", err)
	}

	return doc, nil
}

func (r *DocumentRepository) Delete(
	ctx context.Context,
	id string,
	event document.Event,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE DOCUMENT
		 SET STATUS = :1,
		     UPDATED_AT = :2
		 WHERE ID = :3`,
		document.StatusDeleted,
		event.OccurredAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("erro ao excluir documento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"erro ao verificar documento atualizado: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return service.ErrNotFound
	}

	if err := insertEvent(ctx, tx, event); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

func (r *DocumentRepository) History(
	ctx context.Context,
	id string,
) ([]document.Event, error) {

	if _, err := r.Get(ctx, id); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			ID,
			DOCUMENT_ID,
			EVENT_TYPE,
			OCCURRED_AT,
			CORRELATION_ID
		FROM DOCUMENT_EVENT
		WHERE DOCUMENT_ID = :1
		ORDER BY OCCURRED_AT`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar histórico do documento: %w",
			err,
		)
	}

	defer rows.Close()

	var events []document.Event

	for rows.Next() {
		var event document.Event
		var correlationID sql.NullString

		if err := rows.Scan(
			&event.ID,
			&event.DocumentID,
			&event.Type,
			&event.OccurredAt,
			&correlationID,
		); err != nil {
			return nil, fmt.Errorf(
				"erro ao ler histórico do documento: %w",
				err,
			)
		}

		if correlationID.Valid {
			event.CorrelationID = correlationID.String
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"erro durante leitura do histórico: %w",
			err,
		)
	}

	return events, nil
}

func (r *DocumentRepository) findByIdempotencyKey(
	ctx context.Context,
	idempotencyKey string,
) (document.Document, bool, error) {

	var doc document.Document

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
			ID,
			NAME,
			CONTENT_TYPE,
			STATUS,
			CREATED_AT,
			UPDATED_AT
		FROM DOCUMENT
		WHERE IDEMPOTENCY_KEY = :1`,
		idempotencyKey,
	).Scan(
		&doc.ID,
		&doc.Name,
		&doc.ContentType,
		&doc.Status,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return document.Document{}, false, nil
	}

	if err != nil {
		return document.Document{}, false,
			fmt.Errorf(
				"erro ao consultar chave de idempotência: %w",
				err,
			)
	}

	return doc, true, nil
}

func insertEvent(
	ctx context.Context,
	tx *sql.Tx,
	event document.Event,
) error {

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO DOCUMENT_EVENT (
			ID,
			DOCUMENT_ID,
			EVENT_TYPE,
			OCCURRED_AT,
			CORRELATION_ID
		) VALUES (
			:1, :2, :3, :4, :5
		)`,
		event.ID,
		event.DocumentID,
		event.Type,
		event.OccurredAt,
		nullableString(event.CorrelationID),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir evento: %w", err)
	}

	return nil
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}

	return value
}
