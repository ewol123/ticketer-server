package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/jackskj/carta"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
)

type pgRepository struct {
	client *sql.DB
	connString string
}


func newPgClient(connectionString string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
	}

	_, err = db.Query(`SELECT * FROM "user" WHERE "id"=$1`, "69709b7e-b769-4587-ac0b-5bb99e122c27")
    if err != nil {
		return nil, err
    }

	return db, nil

}

// NewPgRepository : create a new postgres repository
func NewPgRepository(connectionString string) (user.Repository, error) {
	repo := &pgRepository{
		connString: connectionString,
	}

	attempt := 1
	for attempt < 100 {	
		client,err := newPgClient(connectionString)
		if err != nil {
			attempt++
			log.Fatalln(err)
		} else {
			attempt = 100
			repo.client = client
		}
	}

	return repo, nil
}

// Find : find a user in the user migrations by column
func (r *pgRepository) Find(column string, value string) (*user.User, error) {
	 userModel := user.User{}


	 query := fmt.Sprintf(`
	 SELECT 
	 "user"."id" AS user_id,
	 created_at AS user_created_at,
	 updated_at AS user_updated_at,
	 full_name AS user_full_name,
	 email AS user_email,
	 password AS user_password,
	 CASE WHEN registration_code IS NULL THEN '' ELSE registration_code END AS user_registration_code,
	 CASE WHEN reset_password_code IS NULL THEN '' ELSE reset_password_code END AS user_reset_password_code,
	 status AS user_status,
	 r.id AS roles_id,
	 r.name AS roles_name
	 FROM "public"."user"
	 LEFT JOIN 
	 (SELECT 
		id,
		name, 
		user_role.user_id, 
		user_role.role_id 
		FROM "public"."role"
		INNER JOIN 
		"public"."user_role" ON role.id = user_role.role_id )r 
		ON "user"."id" = r.user_id
		WHERE "user"."%v" = '%v';`,column,value)

	 rows, err := r.client.Query(query)
	 if err != nil {
		return nil, errors.Wrap(err, "repository.User.Find")
	 }

	 err = carta.Map(rows, &userModel)
	 if err != nil {
	 	return nil, errors.Wrap(err,"repository.user.Find")
	 }
	 if userModel.Id == "" {
		 return nil, errors.Wrap(user.ErrUserNotFound, "repository.user.Find")
	 }
	 return &userModel, nil
}

func (r *pgRepository) FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string) (*[]user.User, int, error) {
	var users []user.User
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
		whereQuery = `WHERE ("user"."full_name" ILIKE '%`+filter+`%')`
	} else {
		whereQuery = `WHERE true`
	}

	countSql := fmt.Sprintf(`SELECT COUNT(id) FROM "user" %v`, whereQuery)

	sql := fmt.Sprintf(`
	WITH cte AS (SELECT
	"id" AS "user_id",
	"created_at" AS "user_created_at",
	"updated_at" AS "user_updated_at",
	"full_name" AS "user_full_name",
	"email" AS "user_email",
	"password" AS "user_password",
    CASE WHEN registration_code IS NULL THEN '' ELSE registration_code END AS user_registration_code,
    CASE WHEN reset_password_code IS NULL THEN '' ELSE reset_password_code END AS user_reset_password_code,
	"status" AS "user_status"
	FROM "user"
	%v
	)
	SELECT *
	FROM(
	   TABLE  cte
	   ORDER  BY "cte"."%v" %v
	   LIMIT  %v
	   OFFSET %v
	   ) sub;
	`,whereQuery,sortBy, desc,rowsPerPage,offset)

	rows, err := r.client.Query(sql)

	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.User.FindAll")
	}

	err = carta.Map(rows, &users )

	if err != nil {
		return nil, 0, errors.Wrap(err,"repository.user.FindAll")
	}

	result, err := r.client.Query(countSql)

	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.User.FindAll")
	}

	var count int
	for result.Next() {

		if err := result.Scan(&count); err != nil {
			return nil, 0, errors.Wrap(err, "repository.User.FindAll")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "repository.User.FindAll")
	}

	return &users, count, nil

}

func (r *pgRepository) Store(u *user.User) (*user.User, error) {

	tx, err := r.client.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "repository.User.Store")
	}

	_, err =  tx.Exec(`INSERT INTO "user" (id,full_name, email, password,registration_code,status,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,u.Id, u.FullName,u.Email,u.Password,u.RegistrationCode, u.Status,u.CreatedAt,u.UpdatedAt)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, errors.Wrap(rollbackErr, "repository.User.Store")
		}
		return nil, errors.Wrap(err, "repository.User.Store")
	}




	_, err = tx.Exec(`INSERT INTO user_role (user_id,role_id) VALUES($1,$2)`,u.Id, user.USER.Id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, errors.Wrap(rollbackErr, "repository.User.Store")
		}
		return nil, errors.Wrap(err, "repository.User.Store")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "repository.User.Store")
	}
	
	return u, nil
}

func (r *pgRepository) Update(u *user.User) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.User.Update")
	}

	response, err := tx.Exec(`UPDATE "user" SET 
	full_name = case when $1 = '' THEN "user"."full_name" ELSE $1 END,
	email = case when $2 = '' THEN "user"."email" ELSE $2 END,
	password = case when $3 = '' THEN "user"."password" ELSE $3 END,
	registration_code = case when $4 = '' THEN "user"."registration_code" ELSE $4 END,
	status = case when $5 = '' THEN "user"."status" ELSE $5 END
FROM user u WHERE "user"."id" = $6`, u.FullName,u.Email,u.Password,u.RegistrationCode, u.Status,u.Id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "repository.User.Update")
		}
		return errors.Wrap(err, "repository.User.Update")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.User.Update")
	}

	rowsAffected,err := response.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "repository.User.Update")
	}

	if rowsAffected == 0 {
		return errors.Wrap(user.ErrUserNotFound, "repository.User.Update")
	}

	return nil

}

func (r *pgRepository) Delete(id string) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.User.Delete")
	}

	_, err = tx.Exec(`DELETE FROM "user_role" WHERE "user_id" = $1`,id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "repository.User.Delete")
		}
		return errors.Wrap(err, "repository.User.Delete")
	}

	res, err := tx.Exec(`DELETE FROM "user" WHERE "user"."id" = $1`, id)
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "repository.User.Delete")
		}
		return errors.Wrap(err, "repository.User.Delete")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.User.Delete")
	}

	rowsAffected,err := res.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "repository.User.Delete")
	}

	if rowsAffected == 0 {
		return errors.Wrap(user.ErrUserNotFound, "repository.User.Delete")
	}

	return nil
}