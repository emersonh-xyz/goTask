import { MongoClient } from "mongodb";

const uri = process.env.MONGODB_URI || "mongodb+srv://user:user123@cluster0.78576.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
let client: MongoClient;
let clientPromise: Promise<MongoClient>;

if (!global._mongoClientPromise) {
  client = new MongoClient(uri);
  global._mongoClientPromise = client.connect();
}
clientPromise = global._mongoClientPromise;

export default clientPromise;