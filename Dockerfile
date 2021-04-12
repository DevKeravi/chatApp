FROM scratch
MAINTAINER cmkrosp <cmkrosp@naver.com>
ADD serv serv
EXPOSE 8080 8081
ENTRYPOINT ["/serv"]
