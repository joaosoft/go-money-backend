FROM postgres:latest

ADD ./schema/setup/postgres/* /docker-entrypoint-initdb.d/