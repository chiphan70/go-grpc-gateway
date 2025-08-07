// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

db = db.getSiblingDB("grpcgateway_db");

// Create a user collection with some indexes
db.createCollection("users");

// Create indexes for better performance
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ created_at: -1 });
db.users.createIndex({ updated_at: -1 });

// Insert some sample data
db.users.insertMany([
  {
    name: "John Doe",
    email: "john.doe@example.com",
    phone: "+1234567890",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    name: "Jane Smith",
    email: "jane.smith@example.com",
    phone: "+0987654321",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

print("Database initialization completed!");
print("Created users collection with indexes");
print("Inserted sample users");
