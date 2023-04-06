CREATE TABLE users (
    "id" varchar(36) NOT NULL,
    "name" varchar(50) NOT NULL,
    "username" varchar(50) NOT NULL,
    "email" varchar(50) NOT NULL,
    "password" varchar(256) NOT NULL,
    "phone" varchar(16),
    "role_id" varchar(16) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE roles (
    "id" varchar(16) NOT NULL,
    "name" varchar(50) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

ALTER TABLE users
    ADD CONSTRAINT "users_role_id_fkey" 
    FOREIGN KEY ("role_id") 
    REFERENCES roles("id");

INSERT INTO roles ("id", "name", "created_at", "updated_at") VALUES
    ('role001', 'Admin', NOW(), NOW()),
    ('role002', 'Client', NOW(), NOW());
    
INSERT INTO users ("id", "name", "username", "email", "password", "phone", "role_id", "created_at", "updated_at") VALUES
    ('0ed2fd29-d080-4cf3-9019-744deef9c9d8',
    'Diky Afamby Yodihamzah',
    'dikyayodihamzah',
    'dikyayodihamzah@gmail.com',
    '$16$G2okalS3XbmDEznak77Kr.rRdoCk39oT2s4.V/C8TnRqWFHoCqq1a',
    '089684279559',
    'role001',
    NOW(),
    NOW());