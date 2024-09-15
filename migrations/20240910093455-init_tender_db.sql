-- +migrate Up
CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
    );

CREATE TABLE tenders
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(50),
    service_type VARCHAR(50),
    status      tender_status,
    version     INTEGER,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    organization UUID REFERENCES organization (id),
    created_by  VARCHAR REFERENCES employee (username),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR REFERENCES employee (username)
);


-- +migrate Down
drop table if exists tenders;


