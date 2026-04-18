PRAGMA foreign_keys = ON;

CREATE TABLE congregations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    phone VARCHAR(30),
    email VARCHAR(255),
    zip_code VARCHAR(20),
    address VARCHAR(255),
    number VARCHAR(20),
    complement VARCHAR(255),
    neighborhood VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    permission_key VARCHAR(150) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE role_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id)
);

CREATE TABLE members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE,
    gender VARCHAR(30),
    marital_status VARCHAR(50),
    cpf VARCHAR(14),
    rg VARCHAR(20),
    phone VARCHAR(30),
    email VARCHAR(255),
    zip_code VARCHAR(20),
    address VARCHAR(255),
    number VARCHAR(20),
    complement VARCHAR(255),
    neighborhood VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(100),
    joined_at DATE,
    member_status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_members_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE visitors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    assigned_user_id INTEGER,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(255),
    zip_code VARCHAR(20),
    address VARCHAR(255),
    number VARCHAR(20),
    complement VARCHAR(255),
    neighborhood VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(100),
    visit_source VARCHAR(100),
    first_visit_date DATE,
    notes TEXT,
    visitor_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_visitors_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_visitors_assigned_user FOREIGN KEY (assigned_user_id) REFERENCES users (id)
);

CREATE TABLE ministries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    leader_member_id INTEGER,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ministries_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_ministries_leader_member FOREIGN KEY (leader_member_id) REFERENCES members (id)
);

CREATE TABLE member_ministries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id INTEGER NOT NULL,
    ministry_id INTEGER NOT NULL,
    role_name VARCHAR(100),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_member_ministries_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_member_ministries_ministry FOREIGN KEY (ministry_id) REFERENCES ministries (id)
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    responsible_user_id INTEGER,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    type VARCHAR(100) NOT NULL,
    start_datetime DATETIME NOT NULL,
    end_datetime DATETIME,
    location_name VARCHAR(255),
    zip_code VARCHAR(20),
    address VARCHAR(255),
    number VARCHAR(20),
    complement VARCHAR(255),
    neighborhood VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_events_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_events_responsible_user FOREIGN KEY (responsible_user_id) REFERENCES users (id)
);

CREATE TABLE event_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    participant_role VARCHAR(100),
    attendance_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_event_participants_event FOREIGN KEY (event_id) REFERENCES events (id),
    CONSTRAINT fk_event_participants_member FOREIGN KEY (member_id) REFERENCES members (id)
);

CREATE TABLE contributions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    member_id INTEGER,
    recorded_by_user_id INTEGER,
    contribution_type VARCHAR(100) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    contribution_date DATE NOT NULL,
    payment_method VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_contributions_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_contributions_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_contributions_recorded_by_user FOREIGN KEY (recorded_by_user_id) REFERENCES users (id)
);

CREATE TABLE financial_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    congregation_id INTEGER NOT NULL,
    created_by_user_id INTEGER,
    entry_type VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    entry_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_financial_entries_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_financial_entries_created_by_user FOREIGN KEY (created_by_user_id) REFERENCES users (id)
);

CREATE INDEX idx_users_congregation_id ON users (congregation_id);
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
