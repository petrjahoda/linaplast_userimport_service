#!/usr/bin/env bash
cd linux
upx linaplast_userimport_service_linux
cd ..
docker rmi -f petrjahoda/linaplast_userimport:latest
docker build -t petrjahoda/linaplast_userimport:latest .
docker push petrjahoda/linaplast_userimport:latest

docker rmi -f petrjahoda/linaplast_userimport:2021.1.2
docker build -t petrjahoda/linaplast_userimport:2021.1.2 .
docker push petrjahoda/linaplast_userimport:2021.1.2
