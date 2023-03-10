package repository

import (
	"database/sql"
	"errors"
	"final-project-ticketing-api/dto"
	"final-project-ticketing-api/structs"
	"strconv"
	"time"
)

func GetAllTransaction(db *sql.DB) (err error, results []dto.TransactionGet) {
	sqlQuery := `SELECT t.id, t.qr_code, t.created_at, c.id, c.full_name, c.email, c.phone_number, tic.id, tic."name" , tic ."date", tic.price, e.id, e."name" FROM transaction t 
         		INNER JOIN customer c on c.id = t.customer_id
                INNER JOIN ticket tic on tic.id = t.ticket_id
                INNER JOIN "event" e  on tic.event_id  = e.id `
	rows, err := db.Query(sqlQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction = dto.TransactionGet{}
		err = rows.Scan(
			&transaction.ID,
			&transaction.QrCode,
			&transaction.CreatedAt,
			&transaction.CustomerId,
			&transaction.CustomerName,
			&transaction.CustomerEmail,
			&transaction.CustomerPhoneNumber,
			&transaction.TicketId,
			&transaction.TicketName,
			&transaction.TicketDate,
			&transaction.Price,
			&transaction.EventId,
			&transaction.EventName,
		)
		if err != nil {
			panic(err)
		}
		results = append(results, transaction)
	}
	return
}

func GetByTransactionId(db *sql.DB, transactionId int) (err error, result dto.TransactionGet) {
	sqlQuery := `SELECT t.id, t.qr_code, t.created_at, c.id, c.full_name, c.email, c.phone_number, tic.id, tic."name" , tic ."date", tic.price, e.id, e."name" FROM transaction t 
         		INNER JOIN customer c on c.id = t.customer_id
                INNER JOIN ticket tic on tic.id = t.ticket_id
                INNER JOIN "event" e  on tic.event_id  = e.id 
				WHERE t.id = $1`
	var transaction = dto.TransactionGet{}
	rows, err := db.Query(sqlQuery, transactionId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&transaction.ID,
			&transaction.QrCode,
			&transaction.CreatedAt,
			&transaction.CustomerId,
			&transaction.CustomerName,
			&transaction.CustomerEmail,
			&transaction.CustomerPhoneNumber,
			&transaction.TicketId,
			&transaction.TicketName,
			&transaction.TicketDate,
			&transaction.Price,
			&transaction.EventId,
			&transaction.EventName,
		)
		if err != nil {
			panic(err)
		}
		result = transaction
		return nil, transaction
	}
	err = errors.New("transaction with ID : " + strconv.Itoa(transactionId) + " not found")
	return err, transaction
}

func GetTransactionsByCustomerId(db *sql.DB, customerId int) (err error, results []dto.TransactionGet) {
	sqlQuery := `SELECT t.id, t.qr_code, t.created_at, c.id, c.full_name, c.email, c.phone_number, tic.id, tic."name" , tic ."date", tic.price, e.id, e."name" FROM transaction t 
         		INNER JOIN customer c on c.id = t.customer_id
                INNER JOIN ticket tic on tic.id = t.ticket_id
                INNER JOIN "event" e  on tic.event_id  = e.id 
				WHERE t.customer_id = $1`
	var transaction = dto.TransactionGet{}
	rows, err := db.Query(sqlQuery, customerId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&transaction.ID,
			&transaction.QrCode,
			&transaction.CreatedAt,
			&transaction.CustomerId,
			&transaction.CustomerName,
			&transaction.CustomerEmail,
			&transaction.CustomerPhoneNumber,
			&transaction.TicketId,
			&transaction.TicketName,
			&transaction.TicketDate,
			&transaction.Price,
			&transaction.EventId,
			&transaction.EventName,
		)
		if err != nil {
			panic(err)
		}
		results = append(results, transaction)
		return nil, results
	}
	err = errors.New("transaction with customer ID : " + strconv.Itoa(customerId) + " not found")
	return err, results
}

func InsertTransaction(db *sql.DB, transaction structs.Transaction) (structs.Transaction, []error) {
	var errs []error
	sqlQuery := `INSERT INTO transaction (date, qr_code, created_at, updated_at, ticket_id, customer_id) 
				VALUES ($1, $2, $3, $4, $5, $6) 
				Returning *`
	err := db.QueryRow(sqlQuery,
		transaction.Date,
		transaction.QrCode,
		time.Now(),
		time.Now(),
		transaction.TicketId,
		transaction.CustomerId).Scan(
		&transaction.ID,
		&transaction.Date,
		&transaction.QrCode,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.TicketId,
		&transaction.CustomerId)
	if err != nil {
		errs = append(errs, err)
		return transaction, errs
	}
	return transaction, nil
}
