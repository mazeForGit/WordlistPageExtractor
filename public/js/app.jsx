
class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = { 
			time: Date.now(),
			requestexecution: true,
			pagetoscan: "enter url",
			pagescanned: "",
			numberlinksfound: 0,
			numberlinksvisited: 0,
			wordsscanned: 0,
			executionstarted: true,
			executionfinished: false
		};
		this.handleChange = this.handleChange.bind(this);
		this.handleRun = this.handleRun.bind(this);
	}
	handleChange(event) {
		this.setState({pagetoscan: event.target.value});
	}
	handleRun(event) {
		//console.log('.. startExecution');
		this.state.requestexecution = true;
		this.startExecution();
	}
	componentDidMount() {
		//console.log('App componentDidMount')
		this.interval = setInterval(() => this.readConfigData(), 3000);
		//console.log("this.state.pagetoscan = " + this.state.pagescanned);
	}
	componentWillUnmount() {
		clearInterval(this.interval);
	}
	async startExecution() {
		try {
			//console.log('App startExecution ..')
			var reqUrl = ""
			if (window.location.port == "") {
				reqUrl = window.location.protocol + "//" + window.location.hostname + "/config?execution=true";
			} else {
				reqUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/config?execution=true";
			}
			//console.log('request to url = ' + reqUrl);
			//console.log('post data ..');
			const res = await fetch(reqUrl, {
				method: 'POST',
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					pagetoscan: this.state.pagetoscan,
				})
			});
			const blocks = await res.json();
			//console.log(blocks);
			this.setState({
				requestexecution: true,
			})
			//console.log(this.state);
		} catch (e) {
			console.log(e);
		}
	}
	async readConfigData() {
		try {
			if (this.state.requestexecution) {
				//console.log('App readConfigData')
				var reqUrl = ""
				if (window.location.port == "") {
					reqUrl = window.location.protocol + "//" + window.location.hostname + "/config";
				} else {
					reqUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/config";
				}
				//console.log('request to url = ' + reqUrl);
				//console.log('read data ..');
				const res = await fetch(reqUrl);
				const blocks = await res.json();
				const PageScanned = blocks.pagetoscan;
				const NumberLinksFound = blocks.numberlinksfound;
				const NumberLinksVisited = blocks.numberlinksvisited;
				const WordsScanned = blocks.wordsscanned;
				const ExecutionStarted = blocks.executionstarted;
				const ExecutionFinished = blocks.executionfinished;
				//console.log(blocks);

				this.setState({
					time: Date.now(),
					pagescanned: PageScanned,
					numberlinksfound: NumberLinksFound,
					numberlinksvisited: NumberLinksVisited,
					wordsscanned: WordsScanned,
					executionstarted: ExecutionStarted,
					executionfinished: ExecutionFinished,
				})
				//console.log(this.state);
				//console.log("App this.state.pagescanned = " + this.state.pagescanned);
			}
		} catch (e) {
			console.log(e);
		}
	}
	render() {
		return(
			<div className="container">
				<div className="container">
					extract from url = &nbsp;
					<input type="text" size="40" value={this.state.pagetoscan} onChange={this.handleChange} />
					<button onClick={this.handleRun}>run</button>
				</div>	
				<div className="container"> 
					progress from backend : links found = { this.state.numberlinksfound }
					, links visited = { this.state.numberlinksvisited }
					, words extracted = { this.state.wordsscanned }
					, started = { this.state.executionstarted.toString()  }
					, finished = { this.state.executionfinished.toString()  }
				</div>
				{(this.state.requestexecution && this.state.executionfinished) ? <Home pagescanned={this.state.pagescanned} /> : null}
		
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
	
	async serverRequest() {
		try {
			//console.log('Home readData')
			//console.log("Home this.props.pagescanned = " + this.props.pagescanned);
			var reqUrl = ""
			if (window.location.port == "") {
				reqUrl = window.location.protocol + "//" + window.location.hostname + "/wordlist";
			} else {
				reqUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/wordlist";
			}
			//console.log('Home request to url = ' + reqUrl);
			//console.log('Home read data ..');
				
			const res = await fetch(reqUrl);
			const blocks = await res.json();
			//console.log(blocks);
			
			this.setState({
				words: blocks,
			})
			//console.log(this.state);
		} catch (e) {
			console.log(e);
		}
	}
	componentDidMount() {
		//console.log("Home componentDidMount");
		this.serverRequest();
	}
	render() {
		return (
			<div className="container">
				<p>list of words detected at url = {this.props.pagescanned}</p>
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
			words: []
		};	
	}
	render() {
		return (
			<div className="col-xs-2">
				<div className="panel panel-default">
					<div className="panel-heading">{this.props.word.name}</div>
					<div className="panel-body">{this.props.word.occurance}</div>
				</div>
			</div>
		)
	}
}
ReactDOM.render(<App />, document.getElementById('app'));
