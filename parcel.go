package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (address, client, created_at, status) VALUES (:address, :client, :created_at, :status)",
		sql.Named("address", p.Address),
		sql.Named("client", p.Client),
		sql.Named("created_at", p.CreatedAt),
		sql.Named("status", p.Status),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, address, client, created_at, status FROM parcel WHERE number = :n", sql.Named("n", number))
	err := row.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Status)
	if err != nil {
		return Parcel{}, err
	}

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel
	rows, err := s.db.Query("SELECT number, address, client, created_at, status FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Status)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status != ParcelStatusRegistered {
		return err
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status != ParcelStatusRegistered {
		return err
	}
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}
