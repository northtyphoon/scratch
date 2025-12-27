FROM mcr.microsoft.com/acr/acr-cli:0.17 AS acr-cli

FROM mcr.microsoft.com/azure-cli:2.81.0-azurelinux3.0
COPY --from=acr-cli /usr/bin/acr /usr/bin/acr
COPY acr.sh /usr/bin/
RUN chmod +x /usr/bin/acr.sh