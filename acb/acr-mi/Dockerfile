FROM mcr.microsoft.com/acr/acr-cli:0.19 AS acr-cli

FROM mcr.microsoft.com/azure-cli:2.90.0-azurelinux3.0
COPY --from=acr-cli /usr/bin/acr /usr/bin/acr
COPY acr.sh /usr/bin/
RUN chmod +x /usr/bin/acr.sh