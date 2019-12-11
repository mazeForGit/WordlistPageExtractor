
class ProgressComponent extends React.Component {
	constructor(props){
		super(props);
		this.state = { 
			time: Date.now(),
			numberlinksfound: 0,
			numberlinksvisited: 0,
			numberlinksvisited: 0,
			wordsscanned: 0,
			executionfinished: false
		};
	}
	componentDidMount() {
		//console.log('componentDidMount')
		this.loadData()
		this.interval = setInterval(() => this.loadData(), 2000);
		//this.interval = setInterval(() => this.setState({ time: Date.now() }), 10000);
		//setInterval(this.loadData, 5000);
	}
	componentWillUnmount() {
		clearInterval(this.interval);
	}
	async loadData() {
		try {
			//console.log('loadData')
			const res = await fetch('http://localhost:8080/config');
			const blocks = await res.json();
			const NumberLinksFound = blocks.numberlinksfound;
			const NumberLinksVisited = blocks.numberlinksvisited;
			const WordsScanned = blocks.wordsscanned;
			const ExecutionFinished = blocks.executionfinished;
			//console.log(blocks)

			this.setState({
				time: Date.now(),
				numberlinksfound: NumberLinksFound,
				numberlinksvisited: NumberLinksVisited,
				wordsscanned: WordsScanned,
				executionfinished: ExecutionFinished,
			})
		
			//console.log(this.state)
		} catch (e) {
			console.log(e);
		}
	}
	render(){
		return(
			<div> progress from backend : 
				 links found = { this.state.numberlinksfound }
				, links visited = { this.state.numberlinksvisited }
				, words extracted = { this.state.wordsscanned }
				, execution = { this.state.executionfinished.toString()  }
			</div>
		);
	}
}


class TargetUrlForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {value: 'test'};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  handleSubmit(event) {
    alert('A name was submitted: ' + this.state.value);
    event.preventDefault();
  }

  render() {
    return (
		<div>
		<form onSubmit={this.handleSubmit}>
			extract from url = &nbsp;
			<input type="text" value={this.state.value} onChange={this.handleChange} />
			<input type="submit" value="extract" />
      </form>
	  
		<ProgressComponent />
		</div>
    );
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      words: []
    };

    this.serverRequest = this.serverRequest.bind(this);
    this.logout = this.logout.bind(this);
  }

  logout() {
    location.reload();
  }

  serverRequest() {
    $.get("http://localhost:8080/wordlist", res => {
      this.setState({
        words: res
      });
    });
	
  }

  componentDidMount() {
    this.serverRequest();
	console.log(this.state)
  }

  render() {
    return (
      <div className="container">
        <br />
        <span className="pull-right">
          <a onClick={this.logout}>Log out</a>
        </span>
		
		<TargetUrlForm />
		
        <p>list of words detected</p>
        <div className="row">
          <div className="container">
            {this.state.words.map(function(word, i) {
              return <Word key={i} word={word} />;
            })}
          </div>
        </div>
      </div>
    );
  }
}
class Word extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: "",
      words: []
    };
    this.like = this.like.bind(this);
    this.serverRequest = this.serverRequest.bind(this);
  }

	dislike() {
    //let joke = this.props.joke;
    //this.serverRequest(joke);
  }
  like() {
    //let joke = this.props.joke;
    //this.serverRequest(joke);
  }
  serverRequest(word) {
    $.post(
      "http://localhost:8080/api/jokes/like/" + word.id,
      { like: 1 },
      res => {
        console.log("res... ", res);
        this.setState({ liked: "Liked!", words: res });
        this.props.words = res;
      }
    );
  }

  render() {
    return (
      <div className="col-xs-3">
        <div className="panel panel-default">
          <div className="panel-heading">
            #{this.props.word.id}{" "}
			</div>
          <div className="panel-body">{this.props.word.name}</div>
          <div className="panel-footer">
            {this.props.word.occurance} Likes &nbsp;
            <a onClick={this.like} className="btn btn-default">
              <span className="glyphicon glyphicon-thumbs-up" />
            </a>
			<a onClick={this.dislike} className="btn btn-default">
              <span className="glyphicon glyphicon-thumbs-down" />
            </a>
          </div>
        </div>
      </div>
    )
  }
}
ReactDOM.render(<Home />, document.getElementById('app'));
