FROM debian:buster

COPY install-speedtestcli.sh ./install-speedtestcli.sh
RUN chmod +x ./install-speedtestcli.sh
RUN ./install-speedtestcli.sh
ENTRYPOINT [ "speedtest" ]