package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(int) error
	GetAccountByID(int) (*Account, error) // nice
	GetAccounts() ([]*Account, error)

}

type PostgresStore struct {
	db *sql.DB

}


func (s* PostgresStore) Init() error {
	err:= s.createAccountTable()

	return err
}
func (s* PostgresStore) createAccountTable() error {
query := `create table if not exists account (
	id serial primary key, 
	first_name varchar(30), 
	last_name varchar(30), 
	number serial, 
	balance serial,
	created_at timestamp
	 )`

	 _,err := s.db.Exec(query)
	 return err
}

func NewPostgresStore () (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}


	if err:= db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}


func (s* PostgresStore) CreateAccount(account *Account) error {
	query := `
	insert into account (first_name, last_name, number, balance, created_at)
    values ($1, $2, $3, $4, $5);
	`

 if	_, err := s.db.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt); err != nil {
	return err
 }

return nil
}

func (s* PostgresStore) DeleteAccount(id int) error {
	query := `
	delete from account 
	where id= $1
	`;

if _, err := s.db.Query(query, id); err != nil {
	return err
}


return nil
}

func (s* PostgresStore) GetAccountByID(id int) (*Account, error) {
return nil,nil
}


func (s* PostgresStore) UpdateAccount(id int) error {
return nil
}

func (s* PostgresStore) GetAccounts () ([]*Account, error) {
	rows, err := s.db.Query("select * from account;")

	if err != nil {
		return  nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt); err != nil {
				return  nil, err
		}
		accounts = append(accounts, account)
	}

return accounts, nil
}