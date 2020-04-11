#See https://aka.ms/containerfastmode to understand how Visual Studio uses this Dockerfile to build your images for faster debugging.

ARG BaseImgae=mcr.microsoft.com/dotnet/core/aspnet:3.1-buster-slim
ARG BuildImage=mcr.microsoft.com/dotnet/core/sdk:3.1-buster

FROM ${BaseImgae} AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM ${BuildImage} AS build
WORKDIR /src
COPY ["TestServer.csproj", ""]
RUN dotnet restore "./TestServer.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "TestServer.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "TestServer.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "TestServer.dll"]