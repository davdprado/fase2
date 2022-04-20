const grpc = require("grpc");
const protoLoader = require("@grpc/proto-loader");
const packageDef = protoLoader.loadSync("./proto/test.proto", {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});
const grpcObject = grpc.loadPackageDefinition(packageDef);
const mainPackage = grpcObject.main;

const client = new mainPackage.testApi(
  "34.135.125.140:8080",
  grpc.credentials.createInsecure()
);

module.exports = client;
