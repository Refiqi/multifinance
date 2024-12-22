-- Create LoanLimit table
CREATE TABLE loan_limits (
     id VARCHAR(255) PRIMARY KEY,
     tenor_1 FLOAT NOT NULL,
     tenor_2 FLOAT NOT NULL,
     tenor_3 FLOAT NOT NULL,
     tenor_6 FLOAT NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Transaction table
CREATE TABLE transactions (
  id VARCHAR(255) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  otr FLOAT NOT NULL,
  admin_fee FLOAT NOT NULL,
  installments INT NOT NULL,
  interest FLOAT NOT NULL,
  asset_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT unique_transaction UNIQUE (id, user_id)  -- Ensures unique pair of id and user_id
);

-- Create User table
CREATE TABLE users (
   nik VARCHAR(255) PRIMARY KEY,
   full_name VARCHAR(255) NOT NULL,
   legal_name VARCHAR(255) NOT NULL,
   place_of_birth VARCHAR(255) NOT NULL,
   date_of_birth DATE NOT NULL,
   salary FLOAT NOT NULL,
   ktp_photo_url TEXT NOT NULL,
   selfie_photo_url TEXT NOT NULL,
   loan_limit_id VARCHAR(255) NOT NULL,        -- Reference to LoanLimit table
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT fk_loan_limit FOREIGN KEY (loan_limit_id) REFERENCES loan_limits(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_nik ON users(nik);
CREATE INDEX idx_transaction_id_user_id ON transactions(id, user_id);
CREATE INDEX idx_nik_loan_limit_id ON users(nik, loan_limit_id);

