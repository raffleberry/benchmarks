require('dotenv').config()
const mysql = require('mysql2')
const express = require('express')
const app = express()

const port = 3000

const dbPool = mysql.createPool({
  host: process.env.DATABASE_HOST,
  user: process.env.DATABASE_USERNAME,
  database: process.env.DATABASE_NAME,
  password: process.env.DATABASE_PASS,
  waitForConnections: true,
  connectionLimit: 10,
  queueLimit: 50
}).promise()

const isPrime = num => {
  for(let i = 2, s = Math.sqrt(num); i <= s; i++)
      if(num % i === 0) return false; 
  return num > 1;
}

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

app.get('/', (req, res) => {
  res.send('Hello World!')
})

app.get('/bench', async (req, res) => {
  var num = getRandomInt(100)
  if (req.query.id)
    num = Number(req.query.id)
  console.log(req.query)
  var prime = isPrime(num)
  await dbPool.execute(`select ${num}`);
  res.send(prime)
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})

