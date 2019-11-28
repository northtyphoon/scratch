FROM mcr.microsoft.com/dotnet/core/runtime:3.0.1-buster-slim-arm32v7 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/core/sdk:3.0.101-buster AS build
WORKDIR /src
COPY . .
# https://docs.microsoft.com/en-us/dotnet/core/rid-catalog#linux-rids
RUN dotnet publish "helloworld.csproj" -r linux-arm -c Release -o /app --self-contained true /p:PublishTrimmed=true

FROM base AS final
COPY --from=build /app .
ENTRYPOINT ["dotnet", "/app/helloworld.dll"]