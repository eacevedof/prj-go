-- app_users: quien entra al panel. user_role decide las politicas (UserPolicyEnum): root (unico, por consola),
-- admin y owner. api_token_hash es la sesion del panel (rota en cada log-in); NO sirve para enviar mensajes.

CREATE TABLE app_users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uuid VARCHAR(50) NOT NULL UNIQUE,
    user_name VARCHAR(100) NOT NULL,
    user_email VARCHAR(150) NOT NULL UNIQUE COLLATE NOCASE,
    password_hash VARCHAR(255) NOT NULL,
    api_token_hash CHAR(64) NULL UNIQUE,
    user_role VARCHAR(10) NOT NULL DEFAULT 'owner'
        CHECK (user_role IN ('root', 'admin', 'owner')),
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- un solo root
CREATE UNIQUE INDEX uk_users_root ON app_users (user_role) WHERE user_role = 'root';
