-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "clients" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "name" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "images" (
    "id" bigserial,
    "image_string" text,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "bank_types" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "description" text,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "banks" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "type" bigserial,
    "address" text,
    "province" varchar(255),
    "city" varchar(255),
    "pic" varchar(255),
    "phone" varchar(255),
    "adminfee_setup" varchar(255),
    "convfee_setup" varchar(255),
    "services" int ARRAY,
    "products" int ARRAY,
    FOREIGN KEY ("type") REFERENCES bank_types(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "services" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "image_id" bigserial,
    "status" varchar(255),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("image_id") REFERENCES images(id)
) WITH (OIDS = FALSE);

CREATE TABLE "products" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    "service_id" bigserial,
    "min_timespan" int,
    "max_timespan" int,
    "interest" int,
    "min_loan" int,
    "max_loan" int,
    "fees" jsonb DEFAULT '[]',
    "collaterals" varchar(255) ARRAY,
    "financing_sector" varchar(255) ARRAY,
    "assurance" varchar(255),
    FOREIGN KEY ("service_id") REFERENCES services(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "borrowers" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "status" varchar(255),
    "fullname" varchar(255) NOT NULL,
    "gender" varchar(1) NOT NULL,
    "idcard_number" varchar(255) NOT NULL,
    "idcard_image" int,
    "taxid_number" varchar(255),
    "taxid_image" int,
    "email" varchar(255) NOT NULL,
    "birthday" DATE NOT NULL,
    "birthplace" varchar(255) NOT NULL,
    "last_education" varchar(255) NOT NULL,
    "mother_name" varchar(255) NOT NULL,
    "phone" varchar(255) NOT NULL,
    "marriage_status" varchar(255) NOT NULL,
    "spouse_name" varchar(255),
    "spouse_birthday" DATE,
    "spouse_lasteducation" varchar(255),
    "dependants" int DEFAULT (0),
    "address" text NOT NULL,
    "province" varchar(255) NOT NULL,
    "city" varchar(255) NOT NULL,
    "neighbour_association" varchar(255) NOT NULL,
    "hamlets" varchar(255) NOT NULL,
    "home_phonenumber" varchar(255) NOT NULL,
    "subdistrict" varchar(255) NOT NULL,
    "urban_village" varchar(255) NOT NULL,
    "home_ownership" varchar(255) NOT NULL,
    "lived_for" int NOT NULL,
    "occupation" varchar(255) NOT NULL,
    "employee_id" varchar(255),
    "employer_name" varchar(255) NOT NULL,
    "employer_address" text NOT NULL,
    "department" varchar(255) NOT NULL,
    "been_workingfor" int NOT NULL,
    "direct_superiorname" varchar(255),
    "employer_number" varchar(255) NOT NULL,
    "monthly_income" int NOT NULL,
    "other_income" int,
    "other_incomesource" varchar(255),
    "field_of_work" varchar(255) NOT NULL,
    "related_personname" varchar(255) NOT NULL,
    "related_relation" varchar(255) NOT NULL,
    "related_phonenumber" varchar(255) NOT NULL,
    "related_homenumber" varchar(255),
    "related_address" text,
    "bank" bigserial,
    "bank_accountnumber" varchar(255),
    FOREIGN KEY ("bank") REFERENCES banks(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loans" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "owner" bigserial,
    "owner_name" varchar(255),
    "bank" bigserial,
    "product" bigserial,
    "status" varchar(255) DEFAULT ('processing'),
    "loan_amount" FLOAT NOT NULL,
    "installment" int NOT NULL,
    "fees" jsonb DEFAULT '[]',
    "interest" FLOAT NOT NULL,
    "total_loan" FLOAT NOT NULL,
    "due_date" timestamptz,
    "layaway_plan" FLOAT NOT NULL,
    "loan_intention" varchar(255) NOT NULL,
    "intention_details" text NOT NULL,
    "borrower_info" jsonb DEFAULT '[]',
    "disburse_date" timestamptz,
    "disburse_date_changed" BOOLEAN,
    "disburse_status" varchar(255) DEFAULT ('processing'),
    "approval_date" timestamptz,
    "reject_reason" text,
    FOREIGN KEY ("owner") REFERENCES borrowers(id),
    FOREIGN KEY ("bank") REFERENCES banks(id),
    FOREIGN KEY ("product") REFERENCES products(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loan_purposes" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "roles" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "name" varchar(255) NOT NULL,
    "description" text,
    "system" varchar(255),
    "status" varchar(255),
    "permissions" varchar(255) ARRAY,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "users" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "roles" int ARRAY,
    "username" varchar(255) NOT NULL UNIQUE,
    "password" text NOT NULL,
    "email" varchar(255) UNIQUE,
    "phone" varchar(255) UNIQUE,
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "bank_representatives" (
    "id" bigserial,
    "bank_id" bigserial,
    "user_id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("bank_id") REFERENCES banks(id),
    FOREIGN KEY ("user_id") REFERENCES users(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "agent_providers" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "pic" varchar(255),
    "phone" varchar(255) UNIQUE,
    "address" text,
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "agents" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "username" varchar(255) UNIQUE,
    "password" text,
    "email" varchar(255) UNIQUE,
    "phone" varchar(255) UNIQUE,
    "category" varchar(255),
    "agent_provider" bigint,
    "banks" int ARRAY,
    "status" varchar(255),
    FOREIGN KEY ("agent_provider") REFERENCES agent_providers(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS "products" CASCADE;
DROP TABLE IF EXISTS "services" CASCADE;
DROP TABLE IF EXISTS "banks" CASCADE;
DROP TABLE IF EXISTS "bank_types" CASCADE;
DROP TABLE IF EXISTS "borrowers" CASCADE;
DROP TABLE IF EXISTS "loan_purposes" CASCADE;
DROP TABLE IF EXISTS "loans" CASCADE;
DROP TABLE IF EXISTS "images" CASCADE;
DROP TABLE IF EXISTS "clients" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "bank_representatives" CASCADE;
DROP TABLE IF EXISTS "agent_providers" CASCADE;
DROP TABLE IF EXISTS "agents" CASCADE;