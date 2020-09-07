#!/usr/bin/env bash
cd linux
upx linplast_userimport_service_linux
cd ..
docker rmi -f petrjahoda/linplast_userimport_service:latest
docker build -t petrjahoda/linplast_userimport_service:latest .
docker push petrjahoda/linplast_userimport_service:latest

docker rmi -f petrjahoda/linplast_userimport_service:2020.3.3
docker build -t petrjahoda/linplast_userimport_service:2020.3.3 .
docker push petrjahoda/linplast_userimport_service:2020.3.3
