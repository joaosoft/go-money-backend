FROM postgres:9.5

ADD ./schema/setup/postgres/* /docker-entrypoint-initdb.d/