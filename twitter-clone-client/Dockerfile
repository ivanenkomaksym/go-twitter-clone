FROM node:19-alpine

WORKDIR /app

COPY package*.json ./

# Set NODE_ENV to production
ENV NODE_ENV production

RUN npm install

# Create a temporary directory to hold the build context
RUN mkdir /tmp/context

# Copy only the necessary files for the build to the temporary directory
COPY public /tmp/context/public
COPY src /tmp/context/src
COPY package.json /tmp/context

# Copy node_modules to the temporary directory
RUN cp -r node_modules /tmp/context/node_modules

# Build the optimized production build
RUN cd /tmp/context && npm run build

# Install serve as a development dependency
RUN npm install --save-dev serve

# Copy the build files from the temporary directory to /app
RUN cp -a /tmp/context/build /app

EXPOSE 3000

# Start the serve command to serve the production build
CMD ["npx", "serve", "-s", "build"]
