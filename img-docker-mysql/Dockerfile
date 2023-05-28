FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=seidor
ENV MYSQL_DATABASE=test
ENV MYSQL_USER=lucho
ENV MYSQL_PASSWORD=seidor

COPY ./database.sql /docker-entrypoint-initdb.d/