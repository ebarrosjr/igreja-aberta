CREATE TABLE congregations (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    name VARCHAR2(255 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR),
    phone VARCHAR2(30 CHAR),
    email VARCHAR2(255 CHAR),
    zip_code VARCHAR2(20 CHAR),
    address VARCHAR2(255 CHAR),
    number VARCHAR2(20 CHAR),
    complement VARCHAR2(255 CHAR),
    neighborhood VARCHAR2(255 CHAR),
    city VARCHAR2(255 CHAR),
    state VARCHAR2(100 CHAR),
    status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_congregations PRIMARY KEY (id)
);

CREATE TABLE roles (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    name VARCHAR2(100 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_roles PRIMARY KEY (id)
);

CREATE TABLE permissions (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    permission_key VARCHAR2(150 CHAR) NOT NULL,
    name VARCHAR2(150 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_permissions PRIMARY KEY (id)
);

CREATE TABLE users (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    name VARCHAR2(255 CHAR) NOT NULL,
    email VARCHAR2(255 CHAR) NOT NULL,
    password VARCHAR2(255 CHAR) NOT NULL,
    status VARCHAR2(50 CHAR) NOT NULL,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT fk_users_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE user_roles (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    user_id NUMBER(19) NOT NULL,
    role_id NUMBER(19) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_user_roles PRIMARY KEY (id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE role_permissions (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    role_id NUMBER(19) NOT NULL,
    permission_id NUMBER(19) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_role_permissions PRIMARY KEY (id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id)
);

CREATE TABLE members (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    full_name VARCHAR2(255 CHAR) NOT NULL,
    birth_date DATE,
    gender VARCHAR2(30 CHAR),
    marital_status VARCHAR2(50 CHAR),
    cpf VARCHAR2(14 CHAR),
    rg VARCHAR2(20 CHAR),
    phone VARCHAR2(30 CHAR),
    email VARCHAR2(255 CHAR),
    zip_code VARCHAR2(20 CHAR),
    address VARCHAR2(255 CHAR),
    number VARCHAR2(20 CHAR),
    complement VARCHAR2(255 CHAR),
    neighborhood VARCHAR2(255 CHAR),
    city VARCHAR2(255 CHAR),
    state VARCHAR2(100 CHAR),
    joined_at DATE,
    member_status VARCHAR2(50 CHAR) NOT NULL,
    notes CLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_members PRIMARY KEY (id),
    CONSTRAINT fk_members_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id)
);

CREATE TABLE visitors (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    assigned_user_id NUMBER(19),
    full_name VARCHAR2(255 CHAR) NOT NULL,
    phone VARCHAR2(30 CHAR),
    email VARCHAR2(255 CHAR),
    zip_code VARCHAR2(20 CHAR),
    address VARCHAR2(255 CHAR),
    number VARCHAR2(20 CHAR),
    complement VARCHAR2(255 CHAR),
    neighborhood VARCHAR2(255 CHAR),
    city VARCHAR2(255 CHAR),
    state VARCHAR2(100 CHAR),
    visit_source VARCHAR2(100 CHAR),
    first_visit_date DATE,
    notes CLOB,
    visitor_status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_visitors PRIMARY KEY (id),
    CONSTRAINT fk_visitors_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_visitors_assigned_user FOREIGN KEY (assigned_user_id) REFERENCES users (id)
);

CREATE TABLE ministries (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    leader_member_id NUMBER(19),
    name VARCHAR2(255 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR),
    status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_ministries PRIMARY KEY (id),
    CONSTRAINT fk_ministries_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_ministries_leader_member FOREIGN KEY (leader_member_id) REFERENCES members (id)
);

CREATE TABLE member_ministries (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    member_id NUMBER(19) NOT NULL,
    ministry_id NUMBER(19) NOT NULL,
    role_name VARCHAR2(100 CHAR),
    start_date DATE,
    end_date DATE,
    status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_member_ministries PRIMARY KEY (id),
    CONSTRAINT fk_member_ministries_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_member_ministries_ministry FOREIGN KEY (ministry_id) REFERENCES ministries (id)
);

CREATE TABLE events (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    responsible_user_id NUMBER(19),
    title VARCHAR2(255 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR),
    event_type VARCHAR2(100 CHAR) NOT NULL,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP,
    location_name VARCHAR2(255 CHAR),
    zip_code VARCHAR2(20 CHAR),
    address VARCHAR2(255 CHAR),
    number VARCHAR2(20 CHAR),
    complement VARCHAR2(255 CHAR),
    neighborhood VARCHAR2(255 CHAR),
    city VARCHAR2(255 CHAR),
    state VARCHAR2(100 CHAR),
    status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_events PRIMARY KEY (id),
    CONSTRAINT fk_events_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_events_responsible_user FOREIGN KEY (responsible_user_id) REFERENCES users (id)
);

CREATE TABLE event_participants (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    event_id NUMBER(19) NOT NULL,
    member_id NUMBER(19) NOT NULL,
    participant_role VARCHAR2(100 CHAR),
    attendance_status VARCHAR2(50 CHAR) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_event_participants PRIMARY KEY (id),
    CONSTRAINT fk_event_participants_event FOREIGN KEY (event_id) REFERENCES events (id),
    CONSTRAINT fk_event_participants_member FOREIGN KEY (member_id) REFERENCES members (id)
);

CREATE TABLE contributions (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    member_id NUMBER(19),
    recorded_by_user_id NUMBER(19),
    contribution_type VARCHAR2(100 CHAR) NOT NULL,
    amount NUMBER(15,2) NOT NULL,
    contribution_date DATE NOT NULL,
    payment_method VARCHAR2(50 CHAR),
    notes CLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_contributions PRIMARY KEY (id),
    CONSTRAINT fk_contributions_congregation FOREIGN KEY (congregation_id) REFERENCES congregations (id),
    CONSTRAINT fk_contributions_member FOREIGN KEY (member_id) REFERENCES members (id),
    CONSTRAINT fk_contributions_recorded_by_user FOREIGN KEY (recorded_by_user_id) REFERENCES users (id)
);

CREATE TABLE financial_entries (
    id NUMBER(19) GENERATED ALWAYS AS IDENTITY NOT NULL,
    congregation_id NUMBER(19) NOT NULL,
    created_by_user_id NUMBER(19),
    entry_type VARCHAR2(50 CHAR) NOT NULL,
    category VARCHAR2(100 CHAR) NOT NULL,
    description VARCHAR2(1000 CHAR) NOT NULL,
    amount NUMBER(15,2) NOT NULL,
    entry_date DATE NOT NULL,
    status VARCHAR2(50 CHAR) NOT NULL,
    notes CLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_financial_entries PRIMARY KEY (id),
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
