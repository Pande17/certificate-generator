# stage 1 - build stage
FROM node:18-alpine

# set working directory
WORKDIR /app

# install dependencies
COPY package.json ./

RUN npm install

# copy all source code to build it
COPY . ./

# expose port
EXPOSE 5173

# start nginx
CMD [ "npm", "run", "dev" ]
