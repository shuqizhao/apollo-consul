#!/usr/bin/env bash
export app_name=apollo-consul
export service_name=apolloconsul

cp ${service_name}.service /usr/lib/systemd/system/
cp ${service_name}log.conf /etc/rsyslog.d/${service_name}.conf
systemctl daemon-realod
systemctl restart rsyslog
systemctl restart ${service_name}