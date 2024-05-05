-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reports(
                          id bigserial PRIMARY KEY,
                          patient_id UUID NOT NULL,
                          image_path TEXT NOT NULL,
                          report_text TEXT NOT NULL,
                          approved BOOLEAN NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS patients(
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id BIGINT NULL,
                           creator_id BIGINT NOT NULL,
                           "name" TEXT NULL,
                           patient_identifier TEXT NULL,
                           created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users(
                        id bigserial PRIMARY KEY,
                        first_name TEXT NOT NULL,
                        last_name TEXT NOT NULL,
                        email TEXT NOT NULL UNIQUE,
                        password_hash BYTEA NOT NULL,
                        is_doctor BOOLEAN NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    patients ADD CONSTRAINT patients_creator_id_foreign FOREIGN KEY(creator_id) REFERENCES users(id);
ALTER TABLE
    reports ADD CONSTRAINT reports_patient_id_foreign FOREIGN KEY(patient_id) REFERENCES patients(id);
ALTER TABLE
    patients ADD CONSTRAINT patients_user_id_foreign FOREIGN KEY(user_id) REFERENCES users(id);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_modified_time_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER update_modified_time_patients BEFORE UPDATE ON patients FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER update_modified_time_reports BEFORE UPDATE ON reports FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CASCADE IF EXISTS reports;
DROP TABLE CASCADE IF EXISTS patients;
DROP TABLE CASCADE IF EXISTS users;
-- +goose StatementEnd
