express = require("express")
app = express()
app.use "/components", express.static(__dirname + "/components")
app.use "/js", express.static(__dirname + "/js")
app.use "/css", express.static(__dirname + "/css")
app.use express.static(__dirname + "/client")
app.use express.static(__dirname + "/client")
app.listen 3001
