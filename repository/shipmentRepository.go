package repository

import (
	"feklistova/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// CreateShipment registers a new shipment in the database and sets it's ID
func (r *Repository) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	query := `
		INSERT INTO shipments (user_id, projectName, modelType, algorithm, targetColumn, status, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING shipment_id
	`

	err := r.Db.QueryRowContext(ctx, query,
		shipment.UserID,
		shipment.ProjectName,
		shipment.ModelType,
		shipment.Algorithm,
		shipment.TargetColumn,
		shipment.Status,
		shipment.Timestamp,
	).Scan(
		&shipment.ShipmentID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to scan created shipment")
	}

	return nil
}

func (r *Repository) GetShipmentByID(ctx context.Context, shipmentID int) (*models.Shipment, error) {
	var shipment models.Shipment

	query := `
        SELECT user_id, projectName, modelType, algorithm, targetColumn, status, timestamp
        FROM shipments
        WHERE shipment_id = $1
    `

	err := r.Db.QueryRowContext(ctx, query, shipmentID).Scan(
		&shipment.UserID,
		&shipment.ProjectName,
		&shipment.ModelType,
		&shipment.Algorithm,
		&shipment.TargetColumn,
		&shipment.Status,
		&shipment.Timestamp,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("shipment with ID %d not found", shipmentID)
		}
		return nil, errors.Wrap(err, "failed to scan shipment")
	}

	shipment.ShipmentID = shipmentID
	return &shipment, nil
}

func (r *Repository) UpdateShipmentStatus(ctx context.Context, shipment *models.Shipment) error {
	query := `
        UPDATE shipments
        SET status = $1
        WHERE shipment_id = $2
    `

	_, err := r.Db.ExecContext(ctx, query, shipment.Status, shipment.ShipmentID)
	if err != nil {
		return errors.Wrap(err, "failed to update shipment status")
	}

	return nil
}

// DeleteShipment deletes a shipment from the database by its shipment ID
func (r *Repository) DeleteShipment(ctx context.Context, shipmentID int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM shipments WHERE shipment_id = $1", shipmentID)
	if err != nil {
		return err
	}
	return nil
}
