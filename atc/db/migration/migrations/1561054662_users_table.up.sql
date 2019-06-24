BEGIN;
  CREATE TABLE users (
    "id" serial NOT NULL PRIMARY KEY,
    "username" text NOT NULL,
    "connector" text NOT NULL,
    "last_login" timestamp with time zone DEFAULT now() NOT NULL
  );
COMMIT;
