1. Create a local folder:

mkdir ~/ReceiptProcessor

cd ~/ReceiptProcessor

2. Pull the codes

git init

git clone https://github.com/yanghybrid/Receipt-Processor.git


3. Run codes

go run main.go


4.Test the App 

4.1. POST /receipts/process

Using curl:

curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
      { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
    ],
    "total": "18.74"
}'

Expected response:

{"id": "some-uuid"}

4.2. GET /receipts/{id}/points

Replace <id> with the actual UUID you got from the previous call:

curl http://localhost:8080/receipts/<id>/points
