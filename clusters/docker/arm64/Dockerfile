FROM arm64v8/alpine:latest
COPY circonus-agentd /
# NOTE: these are the default ports, use -p to map alternatives configured
EXPOSE 2609/tcp
EXPOSE 8125/udp
ENTRYPOINT ["/circonus-agentd"]