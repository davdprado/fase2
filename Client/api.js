const express = require('express');
const app = express();
var cors = require('cors');
const client = require('./client')
var body_parser = require('body-parser');

//Settings
const port = 3000;

//Middlewares
app.use(express.json());
app.use(body_parser.urlencoded({extended:true}));
app.use(cors());



app.post('/data',  function(req, res) {
    const datos= {
        game_id : req.body.game_id,
        players : req.body.players,
    }
    client.runService('Sendgame', datos, (err, res) => {
        console.log(res.msg);
    })

    res.send('Enviado Correctamente')
});

app.get('/', function(req, res) {
    res.json({ mensaje: '¡Hola Mundo!' })   
  })
/*  */
app.listen(port,()=>{
    console.log('Servidor en el puerto', port);
}); 