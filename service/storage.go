package gomoney

import (
	"database/sql"
	"stock-reservation/domain"
)

// Storage ...
type Storage struct {
	connection *sql.DB
}

// NewStorage ...
func NewStorage(connection *sql.DB) *Storage {
	return &Storage{
		connection: connection,
	}
}

// Get ...
func (repository *Storage) GetUser(id string) (string, error) {
	rows, err := repository.connection.Query(`
			SELECT
				r.reservation_id,
				r.order_id,
				COALESCE(r.origin, ''),
				ri.offer_id,
				ri.channel_id,
				ri.warehouse_id,
				ri.status,
				ri.quantity,
				ri.created_at,
				ri.updated_at
			FROM reservations.reservations r
			RIGHT JOIN reservations.reservation_items ri ON ri.reservation_id = r.reservation_id
			WHERE r.order_id = $1
		`, id)

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return "", err
	}

	reservation := &domain.Reservation{}
	for rows.Next() {
		item := &domain.ReservationItem{}

		if err = rows.Scan(
			&reservation.ReservationID,
			&reservation.OrderID,
			&reservation.Origin,
			&item.OfferID,
			&item.ChannelID,
			&item.WarehouseID,
			&item.Status,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return "", err
		}

		reservation.Items = append(reservation.Items, item)
	}

	if len(reservation.Items) > 0 {
		return "", nil
	}

	return "", nil
}
