package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/jackskj/carta"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type pgRepository struct {
	client     *sql.DB
	connString string
}

func newPgClient(connectionString string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	_, err = db.Query(`SELECT * FROM "ticket" WHERE "id"=$1`, "69709b7e-b769-4587-ac0b-5bb99e122c27")
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}

// NewPgRepository : create a new postgres repository
func NewPgRepository(connectionString string) (ticket.Repository, error) {
	repo := &pgRepository{
		connString: connectionString,
	}

	attempt := 1
	for attempt < 100 {
		client, err := newPgClient(connectionString)
		if err != nil {
			attempt++
			log.Println(err)
			time.Sleep(2 * time.Second)
		} else {
			attempt = 100
			repo.client = client
		}
	}

	if repo.client == nil {
		log.Fatal("can't connect to database for the 100th time, exit")
	}

	return repo, nil
}

// Find : find a ticket in the ticket table by column name
func (r *pgRepository) Find(column string, value string) (*ticket.Ticket, error) {
	ticketModel := ticket.Ticket{}

	query := fmt.Sprintf(`
	 SELECT 
	 "ticket"."id" AS ticket_id,
	 created_at AS ticket_created_at,
	 updated_at AS ticket_updated_at,
     user_id as ticket_user_id,
     CASE WHEN worker_id::varchar IS NULL THEN '' ELSE worker_id::varchar END as ticket_worker_id,
     fault_type as ticket_fault_type,
     address as ticket_address,
     full_name AS ticket_full_name,
	 phone as ticket_phone,
     ST_AsLatLonText(ST_AsText(geo_location)) as ticket_geo_location,
     CASE WHEN image_url IS NULL THEN '' ELSE image_url END as ticket_image_url,
     status as ticket_status
	 FROM "public"."ticket"
     WHERE "ticket"."%v" = '%v';`, column, value)

	rows, err := r.client.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Ticket.Find")
	}

	err = carta.Map(rows, &ticketModel)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Ticket.Find")
	}
	if ticketModel.Id == "" {
		return nil, errors.Wrap(ticket.ErrTicketNotFound, "repository.Ticket.Find")
	}
	return &ticketModel, nil
}

func (r *pgRepository) FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string, workerId string, status string, lat string, long string) (*[]ticket.Ticket, int, error) {
	var tickets []ticket.Ticket
	var desc string
	var whereQuery string

	offsetPage := page - 1
	offset := offsetPage * rowsPerPage

	if descending {
		desc = "DESC"
	} else {
		desc = "ASC"
	}

	if filter != "" {
		whereQuery = fmt.Sprintf(`WHERE "ticket"."full_name" ILIKE '%%%v%%' OR "ticket"."address" ILIKE '%%%v%%' OR "ticket"."phone" ILIKE '%%%v%%'`, filter, filter, filter)
	} else {
		whereQuery = `WHERE true`
	}

	if workerId != "" {
		if workerId == "NULL" {
			whereQuery += fmt.Sprintf(` AND "ticket"."worker_id" IS NULL`)
		} else {
			whereQuery += fmt.Sprintf(` AND "ticket"."worker_id" = '%v'`, workerId)
		}
	}

	if status != "" {
		whereQuery += fmt.Sprintf(` AND "ticket"."status" = '%v'`, status)
	}

	if lat != "" && long != "" {
		whereQuery += fmt.Sprintf(
			` AND ST_DWithin("ticket".geo_location, ST_MakePoint(%v,%v)::geography, 10000)`,
			lat, long)
	}

	countSql := fmt.Sprintf(`SELECT COUNT(id) FROM "ticket" %v`, whereQuery)

	sql := fmt.Sprintf(`
	WITH cte AS (SELECT
	 "ticket"."id" AS ticket_id,
	 created_at AS ticket_created_at,
	 updated_at AS ticket_updated_at,
     user_id as ticket_user_id,
     CASE WHEN worker_id::varchar IS NULL THEN '' ELSE worker_id::varchar END as ticket_worker_id,
     fault_type as ticket_fault_type,
     address as ticket_address,
     full_name AS ticket_full_name,
	 phone as ticket_phone,
     ST_AsLatLonText(ST_AsText(geo_location)) as ticket_geo_location,
     CASE WHEN image_url IS NULL THEN '' ELSE image_url END as ticket_image_url,
     status as ticket_status
	FROM "ticket"
	%v
	)
	SELECT *
	FROM(
	   TABLE  cte
	   ORDER  BY "cte"."%v" %v
	   LIMIT  %v
	   OFFSET %v
	   ) sub;
	`, whereQuery, sortBy, desc, rowsPerPage, offset)

	rows, err := r.client.Query(sql)

	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.Ticket.FindAll")
	}

	err = carta.Map(rows, &tickets)

	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.Ticket.FindAll")
	}

	result, err := r.client.Query(countSql)

	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.Ticket.FindAll")
	}

	var count int
	for result.Next() {

		if err := result.Scan(&count); err != nil {
			return nil, 0, errors.Wrap(err, "repository.Ticket.FindAll")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "repository.Ticket.FindAll")
	}

	return &tickets, count, nil

}

func (r *pgRepository) Store(t *ticket.Ticket) (*ticket.Ticket, error) {

	tx, err := r.client.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Ticket.Store")
	}

	geo := strings.Split(t.GeoLocation, ",")
	lat := geo[0]
	long := geo[1]

	_, err = tx.Exec(`INSERT INTO "ticket" (id, user_id, fault_type, address, full_name,phone,geo_location,image_url,status,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,ST_MakePoint($7,$8),$9,$10,$11,$12)`, t.Id, t.UserId, t.FaultType, t.Address, t.FullName, t.Phone, lat, long, t.ImageUrl, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, errors.Wrap(rollbackErr, "repository.Ticket.Store")
		}
		return nil, errors.Wrap(err, "repository.Ticket.Store")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "repository.Ticket.Store")
	}

	return t, nil
}

func (r *pgRepository) Update(t *ticket.Ticket) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.Ticket.Update")
	}

	lat := "0"
	long := "0"
	geo := strings.Split(t.GeoLocation, ",")

	if len(geo) == 2 {
		lat = geo[0]
		long = geo[1]
	}

	response, err := tx.Exec(`UPDATE "ticket" SET 
	user_id = case when $1 = '' THEN "ticket"."user_id" ELSE $1::uuid END,
	worker_id = case when $2 = '' THEN "ticket"."worker_id" ELSE $2::uuid END,
	fault_type = case when $3 = '' THEN "ticket"."fault_type" ELSE $3::"FaultType" END,
	address = case when $4 = '' THEN "ticket"."address" ELSE $4 END,
	full_name = case when $5 = '' THEN "ticket"."full_name" ELSE $5 END,
	phone = case when $6 = '' THEN "ticket"."phone" ELSE $6 END,
	geo_location = case when $7 = '0' THEN "ticket"."geo_location" ELSE ST_MakePoint($7::double precision, $8::double precision) END,
	image_url = case when $9 = '' THEN "ticket"."image_url" ELSE $9 END,
	status = case when $10 = '' THEN "ticket"."status" ELSE $10::"StatusType" END,
	updated_at = now()
FROM ticket t WHERE "ticket"."id" = $11`, t.UserId, t.WorkerId, t.FaultType, t.Address, t.FullName, t.Phone, lat, long, t.ImageUrl, t.Status, t.Id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "repository.Ticket.Update")
		}
		return errors.Wrap(err, "repository.Ticket.Update")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.Ticket.Update")
	}

	rowsAffected, err := response.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "repository.Ticket.Update")
	}

	if rowsAffected == 0 {
		return errors.Wrap(ticket.ErrTicketNotFound, "repository.Ticket.Update")
	}

	return nil

}

func (r *pgRepository) Delete(id string) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.Ticket.Delete")
	}

	res, err := tx.Exec(`DELETE FROM "ticket" WHERE "ticket"."id" = $1`, id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "repository.Ticket.Delete")
		}
		return errors.Wrap(err, "repository.Ticket.Delete")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.Ticket.Delete")
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "repository.Ticket.Delete")
	}

	if rowsAffected == 0 {
		return errors.Wrap(ticket.ErrTicketNotFound, "repository.Ticket.Delete")
	}

	return nil
}
