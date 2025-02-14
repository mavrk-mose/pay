#!/bin/bash

# Update and install PostgreSQL
sudo yum update -y
sudo amazon-linux-extras enable postgresql14
sudo yum install -y postgresql-server postgresql-contrib

# Initialize and start PostgreSQL
sudo postgresql-setup --initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Set up a PostgreSQL user and database
sudo -i -u postgres psql <<EOF
CREATE DATABASE payment_db;
CREATE USER payment_user WITH ENCRYPTED PASSWORD 'securepassword';
GRANT ALL PRIVILEGES ON DATABASE payment_db TO payment_user;
EOF

# Allow remote connections (modify pg_hba.conf and postgresql.conf)
echo "host all  all  0.0.0.0/0  md5" | sudo tee -a /var/lib/pgsql/data/pg_hba.conf
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /var/lib/pgsql/data/postgresql.conf

# Restart PostgreSQL to apply changes
sudo systemctl restart postgresql
