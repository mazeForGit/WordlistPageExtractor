
class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = { 
			time: Date.now(),
			requestexecution: true,
			pagetoscan: "enter url",
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
		//alert('requested url = ' + this.state.pagetoscan);
		//event.preventDefault();
		//this.sendConfigData();
		console.log('.. startExecution')
		this.state.requestexecution = true;
		this.startExecution();
	}
	componentDidMount() {
		//console.log('componentDidMount')
		//this.readConfigData()
		this.interval = setInterval(() => this.readConfigData(), 5000);
		//this.interval = setInterval(() => this.setState({ time: Date.now() }), 10000);
		//setInterval(this.loadData, 5000);
	}
	componentWillUnmount() {
		clearInterval(this.interval);
	}
	async startExecution() {
		try {
			console.log('startExecution ..')
			const res = await fetch('http://localhost:8080/config?execution=true', {
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
			//console.log(blocks)
			this.setState({
				requestexecution: true,
			})
			
			//console.log(this.state)
		} catch (e) {
			console.log(e);
		}
	}
	async readConfigData() {
		try {
			if (this.state.requestexecution) {
				console.log('readConfigData')
				const res = await fetch('http://localhost:8080/config');
				const blocks = await res.json();
				const NumberLinksFound = blocks.numberlinksfound;
				const NumberLinksVisited = blocks.numberlinksvisited;
				const WordsScanned = blocks.wordsscanned;
				const ExecutionStarted = blocks.executionstarted;
				const ExecutionFinished = blocks.executionfinished;
				console.log(blocks)

				this.setState({
					time: Date.now(),
					numberlinksfound: NumberLinksFound,
					numberlinksvisited: NumberLinksVisited,
					wordsscanned: WordsScanned,
					executionstarted: ExecutionStarted,
					executionfinished: ExecutionFinished,
				})
				console.log(this.state)
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
				
				{(this.state.requestexecution && this.state.executionfinished) ? <Home /> : null}
		
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
			words: []
		};	
	}
	render() {
		return (
			<div className="col-xs-3">
				<div className="panel panel-default">
					<div className="panel-heading">{this.props.word.name}</div>
					<div className="panel-body">{this.props.word.occurance}</div>
				</div>
			</div>
		)
	}
}
ReactDOM.render(<App />, document.getElementById('app'));
