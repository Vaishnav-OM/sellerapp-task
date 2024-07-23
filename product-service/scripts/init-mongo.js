db = db.getSiblingDB("productdb"); // Switch to your database

if (!db.getCollectionNames().includes("products")) {
	db.createCollection("products");
	db.yourCollectionName.insertOne({
		id: "abcdef-12345678",
		name: "product-names",
		description: "product description",
		colour: "red",
		dimensions: "12 cm x 25 cm x 31 cm",
		price: 14.5,
		currencyUnit: "USD",
	}); // Example document
}
