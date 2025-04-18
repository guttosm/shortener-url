// This script initializes a MongoDB database for a URL shortener application.

db = db.getSiblingDB("urlshortener");

db.createCollection("urls");

// If using a replica set
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo:27017" }
  ]
})
