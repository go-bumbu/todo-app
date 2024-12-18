FROM bumbu-todo-builder:latest

WORKDIR /project
COPY . .

RUN make package-ui
RUN goreleaser build --auto-snapshot --clean
RUN chmod -R 0777 dist/