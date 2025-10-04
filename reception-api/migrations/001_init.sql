CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    his_patient_id VARCHAR(100),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hl7_messages (
    id SERIAL PRIMARY KEY,
    message_id VARCHAR(50) UNIQUE NOT NULL,
    patient_id INT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    message_type VARCHAR(10) NOT NULL,
    status VARCHAR(20) DEFAULT 'SENT',
    his_patient_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    ack_received_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_hl7_message_id ON hl7_messages(message_id);
CREATE INDEX IF NOT EXISTS idx_patient_his_id ON patients(his_patient_id);
