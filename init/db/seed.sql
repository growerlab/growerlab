
/* admin user */
TRUNCATE TABLE namespace RESTART IDENTITY;

INSERT INTO namespace ( path, owner_id, type) VALUES ('admin', 1, 1);

INSERT INTO "user" (email,
                    encrypted_password,
                    username,
                    name,
                    public_email,
                    last_login_ip,
                    created_at,
                    deleted_at,
                    verified_at,
                    last_login_at,
                    register_ip,
                    is_admin,
                    namespace_id)
VALUES ('admin@admin.com', '$argon2id$v=19$m=65536,t=1,p=4$r2yY6zOj4vCuQVQ9/71t/Q$FLzA2sWdvOGU4uelTlAWZjnth1C+LDjOfDqDPszvDqA', 'admin', 'admin', 'admin@admin.com', NULL, EXTRACT(EPOCH FROM NOW()), NULL, EXTRACT(EPOCH FROM NOW()), 0, '127.0.0.1', TRUE, (SELECT currval('namespace_id_seq')));

