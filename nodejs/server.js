const express = require('express');
const bodyParser = require('body-parser');
const morgan = require('morgan');
//import logger from 'morgan';

const app = express();
const port = 12001;

app.use(morgan('dev', {}));
//app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

require ('./routes/jsonfilter')(app);
app.listen(port, () => {
  console.log('Listening on port ' + port);
});
