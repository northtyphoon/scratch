FROM mcr.microsoft.com/acr/acr-cli:0.17 AS acr-cli

FROM mcr.microsoft.com/azure-cli:2.81.0-azurelinux3.0
COPY --from=acr-cli /usr/bin/acr /usr/bin/acr
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]