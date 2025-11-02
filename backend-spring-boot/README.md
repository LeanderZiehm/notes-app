```
curl https://start.spring.io/starter.zip \
  -d dependencies=web \
  -d language=java \
  -d type=maven-project \
  -d groupId=com.leanderziehm \
  -d artifactId=notes-app \
  -o template.zip

unzip template.zip
rm template.zip

```


```pom.xml

<dependency>
    <groupId>org.springdoc</groupId>
    <artifactId>springdoc-openapi-starter-webmvc-ui</artifactId>
    <version>2.2.0</version>
</dependency>

```

```
mkdir -p src/main/java/com/example/controller
vim src/main/java/com/example/controller/HelloController.java

```
