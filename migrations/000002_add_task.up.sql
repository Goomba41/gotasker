CREATE TABLE tasks (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	title text NOT NULL,
	description text,
  -- status text NOT NULL DEFAULT 'todo'
  asignee_id bigint NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,

	CONSTRAINT tasks_pk PRIMARY KEY (id),
	CONSTRAINT users_fk
        FOREIGN KEY (asignee_id)
        REFERENCES public.users(id)
);
