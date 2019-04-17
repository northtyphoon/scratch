FROM arm32v7/openjdk:8-jdk-slim-stretch
COPY . /usr/src/myapp
WORKDIR /usr/src/myapp
RUN javac HelloWorld.java
CMD ["java", "HelloWorld"]