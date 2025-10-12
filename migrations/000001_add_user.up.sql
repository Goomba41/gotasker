CREATE TABLE users (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	email text NOT NULL UNIQUE,
	"password" text NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id,email)
);