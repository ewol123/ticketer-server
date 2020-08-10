package postgres

import (
    "database/sql"
    "fmt"
    "log"
	"github.com/ewol123/ticketer-server/user-service/user"
    _ "github.com/lib/pq"
    "github.com/jackskj/carta"
)

type pgRepository struct {
	client *sql.DB
	connString string
}

func newPgClient(connectionString string) (*sql.DB, error) {

	db, err := sqlx.Connect("postgres", connectionString)
    if err != nil {
        return nil, err
	}

	user := user.User{}
	err = db.Get(user, "SELECT * FROM user WHERE id=$1", 1)
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
		}
		else {
			attempt = 100
			repo.client = client
		}
	}

	return repo, nil
}

// Find : find a user in the user db by id
func (r *pgRepository) Find(id int) (*user.User, error) {
	 user := &user.User{}
	 err := r.client.Get(user, `
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
	 if user.Id == 0 {
		 return nil, errors.Wrap(user.ErrRedirectNotFound, "repository.user.Find")
	 }
	 return user, nil
} 

func (r *pgRepository) Store(user *user.User) error {

	tx := db.MustBegin()
	_, err =  tx.NamedExec("INSERT INTO user (full_name, email, password,created_at,updated_at) VALUES (:full_name, :email, :password, :created_at, :updated_at)", &user)
	if err != nil {
		return errors.Wrap(err, "repository.User.Store")
	}
	tx.Commit()
	
	return nil
}