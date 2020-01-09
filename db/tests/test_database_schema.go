package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/mcctor/marauders/db"
)

const (
	passMark = "\u2713"
	failMark = "\u2717"
)

var (
	testDB *sqlx.DB
	schema = `
DROP DATABASE IF EXISTS marauders_test;
CREATE DATABASE marauders_test;
USE marauders_test;

CREATE TABLE users (
	username VARCHAR(20),
	fname VARCHAR(20),
	lname VARCHAR(20),
	email VARCHAR(30) NOT NULL,
	phone VARCHAR(15),
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT pk_users PRIMARY KEY (username)
);

CREATE TABLE passwords (
	user VARCHAR(20),
	salt VARCHAR(50) UNIQUE NOT NULL,
	hash VARCHAR(64) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT pk_passwords PRIMARY KEY (user),
	CONSTRAINT fk_passwords_user FOREIGN KEY (user) REFERENCES users (username) ON DELETE CASCADE
);

CREATE TABLE auth_tokens (
	user VARCHAR(20),
	token VARCHAR(40) NOT NULL,
	refresh_token VARCHAR(20) NOT NULL,
	expiry DATETIME NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT  pk_auth_tokens PRIMARY KEY (user),
	CONSTRAINT fk_auth_tokens_user FOREIGN KEY (user) REFERENCES users (username) ON DELETE CASCADE 
);

CREATE TABLE billings (
	time_stamp DATETIME,
	user VARCHAR(20),
	debit FLOAT NOT NULL DEFAULT 0.0,
	credit FLOAT NOT NULL DEFAULT 0.0,
	CONSTRAINT pk_billings PRIMARY KEY (user, time_stamp),
	CONSTRAINT fk_billings_user FOREIGN KEY (user) REFERENCES users (username) ON DELETE CASCADE
);

CREATE TABLE devices (
	id INT,
	user VARCHAR(20),
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_devices PRIMARY KEY (id, user),
	CONSTRAINT fk_devices_user FOREIGN KEY (user) REFERENCES users (username) ON DELETE CASCADE
);

CREATE TABLE location_snapshots (
	device_id INT,
	time_stamp DATETIME,
	latitude FLOAT NOT NULL,
	longitude FLOAT NOT NULL,
	CONSTRAINT pk_location_snapshots PRIMARY KEY (device_id, time_stamp),
	CONSTRAINT fk_location_snapshots_device FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
);

CREATE TABLE cloaks (
	id VARCHAR(20),
	user VARCHAR(20) NOT NULL,
	name VARCHAR(20) NOT NULL,
	description MEDIUMTEXT NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	wake TIME NOT NULL ,
	sleep TIME NOT NULL ,
	accuracy ENUM('pinpoint', 'street', 'city', 'country') NOT NULL DEFAULT 'street',
	duration DATETIME NOT NULL,
	member_limit INT NOT NULL,
	member_visible BOOLEAN NOT NULL DEFAULT TRUE,
	creator_visible BOOLEAN NOT NULL DEFAULT TRUE,
	everyone_visible BOOLEAN NOT NULL DEFAULT FALSE,
	private BOOLEAN NOT NULL DEFAULT TRUE,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT pk_cloaks PRIMARY KEY (id),
	CONSTRAINT fk_cloaks_user FOREIGN KEY (user) REFERENCES users (username)
);

CREATE TABLE associated_cloaks (
	cloak_id VARCHAR(20),
	device_id INT,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_associated_cloaks PRIMARY KEY (cloak_id, device_id),
	CONSTRAINT fk_associated_cloaks_cloak_id FOREIGN KEY (cloak_id) REFERENCES cloaks (id) ON DELETE CASCADE,
	CONSTRAINT fk_associated_cloaks_device_id FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
);

CREATE TABLE permitted_cloaks (
    id INT AUTO_INCREMENT,
	cloak_id VARCHAR(20),
	permitted_cloak_id VARCHAR(20),
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_permitted_cloaks PRIMARY KEY (id, cloak_id),
	CONSTRAINT fk_permitted_cloak_owning_cloak FOREIGN KEY (cloak_id) REFERENCES cloaks (id) ON DELETE CASCADE,
	CONSTRAINT fk_permitted_cloaks_cloak_id FOREIGN KEY (permitted_cloak_id) REFERENCES cloaks (id) ON DELETE CASCADE
);

CREATE TABLE cloak_invite_links (
	link VARCHAR(12),
	cloak_id VARCHAR(20),
	created_by VARCHAR(20),
	expiry DATETIME NOT NULL,
	added INT NOT NULL DEFAULT 0,
	count_limit INT NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT pk_cloak_invite_links PRIMARY KEY (link, cloak_id),
	CONSTRAINT fk_cloak_invite_links_cloak FOREIGN KEY (cloak_id) REFERENCES cloaks (id) ON DELETE CASCADE,
	CONSTRAINT pk_cloak_invite_link_creator FOREIGN KEY (created_by) REFERENCES users (username) ON DELETE CASCADE 
);

`
)

func init() {
	testDB = sqlx.MustConnect(
		"mysql",
		"mcctor:@lienmwanga01@(localhost:3306)/?multiStatements=True")
	testDB.MustExec(schema)

	// use declared methods against the mock database instead of the actual one
	db.ChangeDB(testDB)
}

