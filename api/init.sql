SET timezone = 'Europe/Paris';

----------- CITIES -----------

CREATE SEQUENCE CitiesIdSeq;

CREATE TABLE Cities (
    Id          INTEGER NOT NULL DEFAULT nextval('CitiesIdSeq'),
    Timestamp   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Name        VARCHAR(255) NOT NULL,
    Touristic   BOOLEAN NOT NULL,
    PRIMARY KEY(Id)
);

ALTER SEQUENCE CitiesIdSeq
OWNED BY Cities.id;

----------- HOUSES -----------

CREATE SEQUENCE HousesIdSeq;

CREATE TABLE Houses (
    Id          INTEGER NOT NULL DEFAULT nextval('HousesIdSeq'),
    Timestamp   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Address     VARCHAR(255) NOT NULL,
    Cityid      INTEGER REFERENCES Cities (Id),
    Inhabitants INTEGER NOT NULL,
    PRIMARY KEY(Id)
);

ALTER SEQUENCE HousesIdSeq
OWNED BY Houses.id;

----------- STORES -----------

CREATE SEQUENCE StoresIdSeq;

CREATE TABLE Stores (
    Id          INTEGER NOT NULL DEFAULT nextval('StoresIdSeq'),
    Timestamp   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Name        VARCHAR(255) NOT NULL,
    Type        VARCHAR(255) NOT NULL,
    Address     VARCHAR(255) NOT NULL,
    Cityid      INTEGER REFERENCES Cities (Id),
    PRIMARY KEY(Id)
);

ALTER SEQUENCE StoresIdSeq
OWNED BY Stores.id;
