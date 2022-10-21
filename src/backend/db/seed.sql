
/* admin user */
INSERT INTO namespace (id, path, owner_id, type) VALUES (1, 'admin', 1, 1);

INSERT INTO `user` (id, email, encrypted_password, username, name, public_email,
                    last_login_ip, created_at, deleted_at, verified_at, last_login_at, register_ip, is_admin, namespace_id)
VALUES (1, 'admin@admin.com', '$argon2id$v=19$m=65536,t=1,p=4$r2yY6zOj4vCuQVQ9/71t/Q$FLzA2sWdvOGU4uelTlAWZjnth1C+LDjOfDqDPszvDqA', 'admin', 'admin', 'admin@admin.com',
        null, UNIX_TIMESTAMP(now()), null, UNIX_TIMESTAMP(now()), null, '127.0.0.1', true, 1);

INSERT INTO server (summary, host, port, state, created_at, deleted_at)
VALUES ('local', '127.0.0.1', 9000, 1, UNIX_TIMESTAMP(now()), null);
