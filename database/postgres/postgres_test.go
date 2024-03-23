package postgres

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {

	opts := &Options{
		Username:           "merchant",
		Password:           "merchant",
		Name:               "merchants",
		Schema:             "cico",
		Host:               "159.89.205.12",
		Port:               5432,
		MinIdleConnections: 10,
		MaxOpenConnections: 30,
		MaxLifetime:        10,
		LogMode:            true,
	}
	conn, err := Connect("master", opts)

	fmt.Println(err)
	fmt.Printf("%#v", conn)

	rows, err := conn.Query("SELECT id, created_date, reference_number, status, third_party_http_status, mfs_reference_number FROM cico.emoney")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var (
			id, thirdPartyHTTPStatus   int64
			createdDate, status, refNo string
			mfsRefNo                   sql.NullString
		)
		if err := rows.Scan(
			&id,
			&createdDate,
			&refNo,
			&status,
			&thirdPartyHTTPStatus,
			&mfsRefNo,
		); err != nil {
			continue
		}

		fmt.Printf("%v %v %v %v %v %v \n\n", id, createdDate, refNo, status, thirdPartyHTTPStatus, mfsRefNo)

	}

}
