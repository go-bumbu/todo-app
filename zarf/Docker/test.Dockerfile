FROM bumbu-todo-builder:latest

WORKDIR /project
COPY . .

RUN make verify