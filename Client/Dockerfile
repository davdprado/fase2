FROM node:17-alpine3.14
RUN mkdir -p /cliente
WORKDIR /cliente
COPY . . 

RUN npm install

EXPOSE 3000

CMD [ "node","api.js" ]