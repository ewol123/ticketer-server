package postgres

import (
	"database/sql"
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

	_, err = db.Query("SELECT * FROM user WHERE id=$1", 1)
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

// Find : find a user in the user db by id
func (r *pgRepository) Find(id string) (*user.User, error) {
	 userModel := user.User{}

	 rows, err := r.client.Query(`
	 SELECT 
	 id AS user_id
	 created_at AS user_created_at
	 updated_at AS user_updated_at
	 full_name AS user_full_name
	 email AS user_email
	 password AS user_password
	 r.id AS roles_id
	 r.name AS roles_name
	 FROM user 
	 INNER JOIN 
	 (SELECT 
		id,
		name, 
		ur.user_id, 
		ur.role_id 
		FROM role 
		INNER JOIN 
		user_role ON role.id = user_role.role_id )r 
		ON user.id = r.user_id
		WHERE user.id = $1`, id)
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
		whereQuery = `WHERE (user.first_name ILIKE '%`+filter+`%') OR (user.last_name ILIKE '%`+filter+`%')`
	} else {
		whereQuery = `WHERE true`
	}

	rows, err := r.client.Query(`
	WITH cte AS (SELECT
	"user".*
	FROM user
	$1
	)
	SELECT *
	FROM(
	   TABLE  cte
	   ORDER  BY "cte"."$2" $3
	   LIMIT  $4
	   OFFSET $5
	   ) sub  
	RIGHT JOIN (SELECT count(*) FROM cte) c(full_count) ON true;
	`, whereQuery, sortBy, desc,rowsPerPage,offset)


	if err != nil {
		return nil, 0, errors.Wrap(err, "repository.User.Find")
	}

	err = carta.Map(rows, &users )

	if err != nil {
		return nil, 0, errors.Wrap(err,"repository.user.Find")
	}

	return &users, 1, nil

}

func (r *pgRepository) Store(user *user.User) error {

	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	_, err =  tx.Exec("INSERT INTO user (full_name, email, password,created_at,updated_at) VALUES ($1,$2,$3,$4)", user.FullName,user.Email,user.Password,user.CreatedAt,user.UpdatedAt )
	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(err, "repository.User.Store")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}
	
	return nil
}

func (r *pgRepository) Update(u *user.User) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	response, err := tx.Exec(`UPDATE user 
								SET 
									full_name = case when $1 = '' THEN user.full_name ELSE $1 END,
								    email = case when $2 = '' THEN user.email ELSE $2 END,
								    password = case when $3 = '' THEN user.password ELSE $3 END 
								FROM user WHERE user.id = $4`, u.FullName,u.Email,u.Password,u.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(err, "repository.User.Store")
		}
		return errors.Wrap(err, "repository.User.Store")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	rowsAffected,err := response.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	if rowsAffected == 0 {
		return errors.Wrap(user.ErrUserNotFound, "repository.User.Store")
	}

	return nil

}

func (r *pgRepository) Delete(id string) error {
	tx, err := r.client.Begin()
	if err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	_, err = tx.Exec("DELETE FROM user WHERE user.id = $1", id)

	if err != nil {
		log.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(err, "repository.User.Store")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}

	return nil
}