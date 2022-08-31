
Simple User Management API in Golang and HANA Cloud Database


# Prerequisites
- SAP HANA Cloud DB instance
- SAP BTP Kyma environment instance
- 'DigiCertGlobalRootCA.pem' file should be downloaded from [DigiCert](https://www.digicert.com/digicert-root-certificates.htm) in PEM format and added to the root directory of the project(same level of go.mod).


# Schema creation query
```
CREATE SCHEMA USERAPI;
```

# Table creation query
```
CREATE COLUMN TABLE UserApi.User(
username VARCHAR(36) PRIMARY KEY,
email VARCHAR(20) ,
firstName VARCHAR(30),
lastName VARCHAR(30) ,
address VARCHAR(50) ,
mobile INTEGER
);
```
# Required modifications in secret.yaml

- HDB_USER -> HANA DB user name(e.g. DBADMIN )
- HDB_PASSWORD -> HDB_USER its password
- HDB_HOST & HDB_PORT -> SQL Endpoint of HANA DB