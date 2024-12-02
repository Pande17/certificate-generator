# stage 1 - build stage
FROM node:18-alpine

# set working directory
WORKDIR /app

# install dependencies
COPY package*.json ./

RUN npm install
RUN npm install serve -g

# copy all source code to build it
COPY . ./

# build the react app
RUN npm run build

EXPOSE 5173

CMD [ "serve", "-s", "dist" ]
