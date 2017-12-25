package postgres

import (
	"context"
	"log"
	"time"
)

type Alert struct {
	ID         int
	CurrencyID int
	Currency   string
	Price      float64
	LessThan   bool
	Active     bool
	Alerting   bool
}

// return an array of all active alerts
func (p *Postgres) GetAlerts() ([]*Alert, error) {
	query := `SELECT alert.id, currency_id, code, price, less from alert join currency on alert.currency_id = currency.id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := p.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*Alert
	for rows.Next() {
		alert := &Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.CurrencyID,
			&alert.Currency,
			&alert.Price,
			&alert.LessThan,
		)
		if err != nil {
			log.Printf("Error reading alert: %s", err)
			continue // skip this record
		}

		alert.Active = true
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
