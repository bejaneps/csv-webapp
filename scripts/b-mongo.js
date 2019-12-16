// newdb is the new database we create
var url = "mongodb://localhost:27017/cdr";

// create a client to mongodb
import { MongoClient } from 'mongodb';

// make client connect to mongo service
MongoClient.connect(url, function (err, db) {
	if (err) throw err;
	console.log("Database created!");
	// print database name
	console.log("db object points to the database : " + db.databaseName);

	db.createCollection("users")
	db.users.insert({
		email: 'paul',
		password: 'paul123'
	})
	
	// after completing all the operations with db, close it.
	db.close();
}); 