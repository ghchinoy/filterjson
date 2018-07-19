import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
//import Dialog from '@material-ui/core/Dialog';
//import DialogTitle from '@material-ui/core/DialogTitle';
//import DialogContent from '@material-ui/core/DialogContent';
//import DialogContentText from '@material-ui/core/DialogContentText';
//import DialogActions from '@material-ui/core/DialogActions';
import Typography from '@material-ui/core/Typography';
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Paper from '@material-ui/core/Paper'
import Grid from '@material-ui/core/Grid'
import TextField from '@material-ui/core/TextField';
import { withStyles } from '@material-ui/core/styles';
import withRoot from '../withRoot';

const styles = theme => ({
  root: {
    textAlign: 'center',
    //paddingTop: theme.spacing.unit * 20,
  },
  paper: {
    ...theme.mixins.gutters(),
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    //width: 200,
  },
  button: {
    margin: theme.spacing.unit,
  },
});

class Index extends React.Component {
  state = {
    //open: false,
    filter: "",
    remove: "",
    unfiltered: {
      "hello": 1,
      "there": 2
    },
    filtered: {}
  }; 

  handleFilterInvokation = () => {
    let service = `${window.location.protocol}//${window.location.hostname}:12001`;
    //console.log(this.state.unfiltered);
    try {
      //console.log("Attempting fetch...")
      let body = JSON.stringify(this.state.unfiltered)
      fetch(`${service}/?filter=${this.state.filter}`,{
        method: 'post',
        headers: {'content-type': 'application/json'},
        body: body
      })
      .then( response => response.json())
      .then( data => {
        //console.log(data);
        this.setState({filtered: data});
      });
    } catch(err) {
      console.log("error",err);
    }
  };

  handleParameterChanges = key => (event, value) => {
    //console.log("paramchange", key)
    //console.log(key,event.target.value)
    this.setState({
      [key]: event.target.value,
    });
  };
  handleInputJSONChange = (e) => {
    //console.log("jsonchange")
    try {
      let unfiltered = JSON.parse(e.target.value)
      this.setState( {
        unfiltered: unfiltered
      });
    } catch(err) {
      console.log("error",err);
    }
  };

  hello = () => {
    console.log("hello");
  };

  

  render() {
    const { classes } = this.props;
    //const { open } = this.state;

    return (
      <div className={classes.root}>
        <AppBar position="static" color="default">
          <Toolbar>
            <Typography variant="title">
              FilterJSON tester
            </Typography>
          </Toolbar>
        </AppBar>
        <Paper className={classes.paper} elevation={1}>
          <Grid container spacing={24}>
            <Grid item xs={12}>
              <TextField
                label="Filter"
                id="param-filter"
                className={classes.textField}
                helperText="comma-separated keys to select"
                margin="dense"
                value={this.state.filter || ''}
                onChange={this.handleParameterChanges('filter')}
              />
              <TextField
                label="Remove"
                id="param-remove"
                className={classes.textField}
                helperText="comma-separated keys to remove"
                margin="dense"
                value={this.state.remove || ''}
                onChange={this.handleParameterChanges('remove')}
              />
              <Button 
                variant="contained" 
                color="primary" 
                onClick={this.handleFilterInvokation.bind(this)}
                className={classes.button}
                >
                Test
              </Button>
            </Grid>
            <Grid item xs={6}>
              <TextField
                id="multiline-flexible-unfiltered"
                label="Input"
                multiline
                rowsMax="20"
                fullWidth
                className={classes.textField}
                margin="normal"
                helperText="Unfiltered JSON"
                value={JSON.stringify(this.state.unfiltered, null, "  ") || ''}
                onChange={this.handleInputJSONChange}
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                id="multiline-flexible-filtered"
                label="Output"
                multiline
                rowsMax="20"
                fullWidth
                className={classes.textField}
                margin="normal"
                helperText="Filter JSON" 
                value={JSON.stringify(this.state.filtered, null, "  ") || ''}
              />
            </Grid>
          </Grid>
        </Paper>
      </div>
    );
  }
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withRoot(withStyles(styles)(Index));
