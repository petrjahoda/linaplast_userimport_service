#!/usr/bin/env bash
cd linux
upx linplast_userimport_service_linux
cd ..
docker rmi -f petrjahoda/linaplast_userimport_service:latest
docker build -t petrjahoda/linaplast_userimport_service:latest .
docker push petrjahoda/linaplast_userimport_service:latest

docker rmi -f petrjahoda/linaplast_userimport_service:2020.3.3
docker build -t petrjahoda/linaplast_userimport_service:2020.3.3 .
docker push petrjahoda/linaplast_userimport_service:2020.3.3
