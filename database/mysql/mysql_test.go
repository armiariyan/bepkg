package mysql

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestMysqlSelect(t *testing.T) {

	opts := &Options{
		DSN:                "root:password@tcp(0.0.0.0:3306)/mysql",
		MinIdleConnections: 10,
		MaxOpenConnections: 30,
		MaxLifetime:        10,
		LogMode:            true,
	}
	conn, err := Connect("master", opts)

	rows, err := conn.Query("SELECT host, user FROM user")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var (
			host, user string
		)
		if err := rows.Scan(
			&host,
			&user,
		); err != nil {
			continue
		}

		fmt.Printf("%v %v \n\n", host, user)

	}

}

func TestMysqlCreateTable(t *testing.T) {

	opts := &Options{
		DSN:                "root:password@tcp(0.0.0.0:3306)/dbtest",
		MinIdleConnections: 10,
		MaxOpenConnections: 30,
		MaxLifetime:        10,
		LogMode:            true,
	}
	conn, err := Connect("master", opts)

	conn.Exec(`DROP TABLE IF EXISTS places;`)

	_, err = conn.Exec(`CREATE TABLE places (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		name varchar(25) NOT NULL,
		city varchar(25) DEFAULT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1;
	  `)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func TestMysqlInsert(t *testing.T) {

	opts := &Options{
		DSN:                "root:password@tcp(0.0.0.0:3306)/dbtest",
		MinIdleConnections: 10,
		MaxOpenConnections: 30,
		MaxLifetime:        10,
		LogMode:            true,
	}
	conn, err := Connect("master", opts)

	_, err = conn.Exec(`INSERT INTO places (id, name, city) 
						VALUES ('1', 'sudirman', 'jakarta'),
							   ('2', 'jonggol', 'jakarta'),
							   ('3', 'judir', NULL),
							   ('4', 'asdfa', NULL);
						`)

	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := conn.Query("SELECT id, name, city FROM places")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var (
			id   int64
			name sql.NullString
			city string
		)
		if err := rows.Scan(
			&id,
			&name,
			&city,
		); err != nil {
			continue
		}

		fmt.Printf("%v %#v %v\n\n", id, name, city)

	}

}
