CREATE TABLE congregations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NULL,
    phone VARCHAR(30) NULL,
    email VARCHAR(255) NULL,
    zip_code VARCHAR(20) NULL,
    address VARCHAR(255) NULL,
    number VARCHAR(20) NULL,
    complement VARCHAR(255) NULL,
    neighborhood VARCHAR(255) NULL,
    city VARCHAR(255) NULL,
    state VARCHAR(100) NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_congregations PRIMARY KEY (id)
);

CREATE TABLE roles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_roles PRIMARY KEY (id)
);

CREATE TABLE permissions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `permission_key` VARCHAR(150) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(1000) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_permissions PRIMARY KEY (id)
);

CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT fk_users_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE refresh_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_refresh_tokens PRIMARY KEY (id),
    CONSTRAINT uq_refresh_tokens_token_hash UNIQUE (token_hash),
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE password_reset_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_password_reset_tokens PRIMARY KEY (id),
    CONSTRAINT uq_password_reset_tokens_token_hash UNIQUE (token_hash),
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE user_roles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_user_roles PRIMARY KEY (id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_refresh_tokens PRIMARY KEY (id),
    CONSTRAINT uq_refresh_tokens_token_hash UNIQUE (token_hash),
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_password_reset_tokens PRIMARY KEY (id),
    CONSTRAINT uq_password_reset_tokens_token_hash UNIQUE (token_hash),
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE role_permissions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_role_permissions PRIMARY KEY (id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id)
);

CREATE TABLE members (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NULL,
    gender VARCHAR(30) NULL,
    marital_status VARCHAR(50) NULL,
    cpf VARCHAR(14) NULL,
    rg VARCHAR(20) NULL,
    phone VARCHAR(30) NULL,
    email VARCHAR(255) NULL,
    zip_code VARCHAR(20) NULL,
    address VARCHAR(255) NULL,
    number VARCHAR(20) NULL,
    complement VARCHAR(255) NULL,
    neighborhood VARCHAR(255) NULL,
    city VARCHAR(255) NULL,
    state VARCHAR(100) NULL,
    joined_at DATE NULL,
    member_status VARCHAR(50) NOT NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_members PRIMARY KEY (id),
    CONSTRAINT fk_members_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE visitors (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    assigned_user_id BIGINT UNSIGNED NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(30) NULL,
    email VARCHAR(255) NULL,
    zip_code VARCHAR(20) NULL,
    address VARCHAR(255) NULL,
    number VARCHAR(20) NULL,
    complement VARCHAR(255) NULL,
    neighborhood VARCHAR(255) NULL,
    city VARCHAR(255) NULL,
    state VARCHAR(100) NULL,
    visit_source VARCHAR(100) NULL,
    first_visit_date DATE NULL,
    notes TEXT NULL,
    visitor_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_visitors PRIMARY KEY (id),
    CONSTRAINT fk_visitors_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_visitors_assigned_user FOREIGN KEY (assigned_user_id) REFERENCES users (id)
);

CREATE TABLE ministries (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    leader_member_id BIGINT UNSIGNED NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_ministries PRIMARY KEY (id),
    CONSTRAINT fk_ministries_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_ministries_leader_member FOREIGN KEY (leader_member_id) REFERENCES members (id)
);

CREATE TABLE member_ministries (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    member_id BIGINT UNSIGNED NOT NULL,
    ministry_id BIGINT UNSIGNED NOT NULL,
    role_name VARCHAR(100) NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_member_ministries PRIMARY KEY (id),
    CONSTRAINT fk_member_ministries_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_member_ministries_ministry FOREIGN KEY (ministry_id) REFERENCES ministries (id)
);

CREATE TABLE events (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    responsible_user_id BIGINT UNSIGNED NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NULL,
    type VARCHAR(100) NOT NULL,
    start_datetime DATETIME NOT NULL,
    end_datetime DATETIME NULL,
    location_name VARCHAR(255) NULL,
    zip_code VARCHAR(20) NULL,
    address VARCHAR(255) NULL,
    number VARCHAR(20) NULL,
    complement VARCHAR(255) NULL,
    neighborhood VARCHAR(255) NULL,
    city VARCHAR(255) NULL,
    state VARCHAR(100) NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_events PRIMARY KEY (id),
    CONSTRAINT fk_events_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_events_responsible_user FOREIGN KEY (responsible_user_id) REFERENCES users (id)
);

CREATE TABLE event_participants (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    event_id BIGINT UNSIGNED NOT NULL,
    member_id BIGINT UNSIGNED NOT NULL,
    participant_role VARCHAR(100) NULL,
    attendance_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_event_participants PRIMARY KEY (id),
    CONSTRAINT fk_event_participants_event FOREIGN KEY (event_id) REFERENCES events (id),
    CONSTRAINT fk_event_participants_member FOREIGN KEY (member_id) REFERENCES members (id)
);

CREATE TABLE contributions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    member_id BIGINT UNSIGNED NULL,
    recorded_by_user_id BIGINT UNSIGNED NULL,
    contribution_type VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    contribution_date DATE NOT NULL,
    payment_method VARCHAR(50) NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_contributions PRIMARY KEY (id),
    CONSTRAINT fk_contributions_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_contributions_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_contributions_recorded_by_user FOREIGN KEY (recorded_by_user_id) REFERENCES users (id)
);

CREATE TABLE financial_entries (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    congregation_id BIGINT UNSIGNED NOT NULL,
    created_by_user_id BIGINT UNSIGNED NULL,
    entry_type VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    entry_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_financial_entries PRIMARY KEY (id),
    CONSTRAINT fk_financial_entries_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_financial_entries_created_by_user FOREIGN KEY (created_by_user_id) REFERENCES users (id)
);

CREATE INDEX idx_users_congregation_id ON users (congregation_id);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions (permission_id);
CREATE INDEX idx_members_congregation_id ON members (congregation_id);
CREATE INDEX idx_visitors_congregation_id ON visitors (congregation_id);
CREATE INDEX idx_visitors_assigned_user_id ON visitors (assigned_user_id);
CREATE INDEX idx_ministries_congregation_id ON ministries (congregation_id);
CREATE INDEX idx_ministries_leader_member_id ON ministries (leader_member_id);
CREATE INDEX idx_member_ministries_member_id ON member_ministries (member_id);
CREATE INDEX idx_member_ministries_ministry_id ON member_ministries (ministry_id);
CREATE INDEX idx_events_congregation_id ON events (congregation_id);
CREATE INDEX idx_events_responsible_user_id ON events (responsible_user_id);
CREATE INDEX idx_event_participants_event_id ON event_participants (event_id);
CREATE INDEX idx_event_participants_member_id ON event_participants (member_id);
CREATE INDEX idx_contributions_congregation_id ON contributions (congregation_id);
CREATE INDEX idx_contributions_member_id ON contributions (member_id);
CREATE INDEX idx_contributions_recorded_by_user_id ON contributions (recorded_by_user_id);
CREATE INDEX idx_financial_entries_congregation_id ON financial_entries (congregation_id);
CREATE INDEX idx_financial_entries_created_by_user_id ON financial_entries (created_by_user_id);
