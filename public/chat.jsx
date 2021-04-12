
class App extends React.Component {
	constructor(props){
		super(props);

	}
	render() {
		return (
			<div class="ChatAppBody">
				<ChatApp />
			</div>
		)
	}
}
class ChatApp extends React.Component {

	constructor(props){
		super(props);
		this.state = {
			ws: null
		};
		
	}

	componentDidMount(){
		this.connect();
	}

	timeout = 250;

	connect = () => {
		var ws = new WebSocket('ws://cmkrosp.iptime.org:8080/room');
		let that = this;
		var connectInterval;

		ws.onopen= () => {
			console.log("connected websocket main component");
			this.setState({ws: ws});
			that.timeout = 250;
			clearTimeout(connectInterval);
		};

		ws.onclose = e => {
			console.log(
				'socket is closed'
			)

			that.timeout= that.timeout + that.timeout;
			connectInterval = setTimeout(this.check,Math.min(10000,that.timeout));
		};
		ws.onerror = err => {
			console.error("Socket encountered error: ",err.message,"Closing socket");
		};
		ws.onmessage = (evt) => {
			console.log("onmessage: ", evt.data)
			$('.ChatBody').prepend('<p>' + evt.data + '</p>');

		}

		// Don't use when keep using socket
	//	ws.close();
	};

 
	check = () => {
		const { ws } = this.state;
		if (!ws || ws.readyState == WebSocket.CLOSED) this.connect();
	};

	render() {
		return (

			<div class="Container">
				<div class="row">
					<div class="col col-md-3"></div>
					<div class="col col-md-6">
						<ChatHead />
						<ChatInsert websocket={this.state.ws} />
						<ChatBody websocket={this.state.ws} />
					</div>
					<div class="col col-md-3"></div>
				</div>
			</div>
		)
	}
}

class ChatHead extends React.Component{

	render() {
		return (
			<div class="ChatHead">
				<h1><b>Chat App WebSocket</b> <span class="badge bg-warning">TEST</span></h1>
			</div>
		)
	}
}
class ChatInsert extends React.Component {
	constructor(props){
		super(props);
	}
	handleSubmit = (event) => {
		event.preventDefault();
		const {websocket} = this.props;
		var item = $('#ChatInput').val();
		if(item){
			try {
				console.log("websocket send:",item);
				websocket.send(item);
			} catch (error) {
				console.log(error);
			}

			$('#ChatInput').val("");
		}

	};
	render() {
		return (
			<div class="ChatForm">
				<form onSubmit={this.handleSubmit}>
					<div class="mb-3">
						<label for="ChatInput" class="form-label">Input your Chat</label>
						<input type="text" id="ChatInput" class="form-control" />
						<div id="inputHelpBlock" class="form-text">
							This chat program built by WebSocket
						</div>
						<button type="submit" class="btn btn-primary chatsub">Submit</button>
					</div>
				</form>
			</div>


		)
	}
}

class ChatBody extends React.Component {
	render() {
		return (
			<div class="ChatBody">
			</div>

		)

	}
}

ReactDOM.render(<App />, document.getElementById('root'));
