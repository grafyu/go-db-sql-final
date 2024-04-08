package main

import (
	"database/sql"
	"fmt"
	"log"
)

// struct хранящая ссылку на struct sql.DB
type ParcelStore struct {
	db *sql.DB
}

// создает структуру хранящую ссылку на открытое подключение к db
func NewParcelStore(driverName string, dataSourceName string) ParcelStore {
	var pStore ParcelStore

	if driverName == "sqlite" {
		db, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			fmt.Println("not connect")
			log.Println(err)
		}
		pStore.db = db
	}

	return pStore
}


// реализует добавление строки в таблицу parcel
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(number), nil
}


// реализует чтение строки по заданному number
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

// реализует чтение строк из таблицы parcel по заданному client
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return []Parcel{}, err
	}
	defer rows.Close()

	var res []Parcel

	for i := 0; rows.Next(); i++ {
		var temp Parcel
		err = rows.Scan(&temp.Number, &temp.Client, &temp.Status, &temp.Address, &temp.CreatedAt)
		if err != nil {
			return []Parcel{}, err
		}
		res = append(res, temp)
	}
	return res, nil
}

// реализует обновление статуса в таблице parcel
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	return nil
}

// реализует обновление адреса в таблице parcel
func (s ParcelStore) SetAddress(number int, address string) error {
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number))

	var status string

	err := row.Scan(&status)
	if err != nil {
		return err
	}

	if status == "registered" {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}


// реализует удаление строки из таблицы parcel
func (s ParcelStore) Delete(number int) error {
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number))

	var status string

	err := row.Scan(&status)
	if err != nil {
		return err
	}

	if status == "registered" {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))
	}
	if err != nil {
		return err
	}

	return nil
}
