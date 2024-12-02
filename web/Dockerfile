# stage 1 - build stage
FROM node:18-alpine

# set working directory
WORKDIR /app

# install dependencies
COPY package*.json ./

RUN npm install

# copy all source code to build it
COPY . ./

# build the react app
RUN npm run build

# Debugging: Check the output of the build
RUN ls -la /app/dist

ENTRYPOINT [ "/app/dist" ]