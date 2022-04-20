const GRPCClient = require('node-grpc-client')

const path = require("path")

const PROTO_PATH = path.resolve(__dirname, '../proto/test.proto')

var client = new GRPCClient(PROTO_PATH, 'main', 'testApi', '0.0.0.0:8080')


module.exports = client;