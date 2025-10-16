CREATE TABLE audit_logs (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	user_id bigint NULL,
	user_email text NULL,
	"action" text NOT NULL,
	entity_type text NOT NULL,
	entity_id bigint NOT NULL,
	old_values jsonb NULL,
	new_values jsonb NULL,
	ip_address inet NULL,
	user_agent text NULL,
	request_id uuid NULL,
	created_at time with time zone DEFAULT now() NOT NULL,
	CONSTRAINT audit_logs_pk PRIMARY KEY (id)
);

CREATE INDEX audit_entity_idx ON audit_logs (entity_type,entity_id);
CREATE INDEX audit_user_idx ON audit_logs (user_id);
CREATE INDEX audit_action_idx ON audit_logs ("action");