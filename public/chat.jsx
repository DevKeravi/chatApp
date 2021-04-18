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

	createItem = (msg) => {
		var temp ='<div class="MsgContainer"><p class=ChatAvatar><img class=AvatarImg src='+ msg.AvatarURL +' /></p><p class=ChatMsgBody><p class=ChatMsg><strong>' + msg.Name + '</strong>' + ' : ' + msg.Message + '</p><p class=ChatWhen>' + msg.When + '</p></p></div>';
		return temp;
	}

	componentDidMount(){
		this.connect();
		axios.get("https://cmkrosp.iptime.org:8080/Chatlog")
			.then(resp => {
				resp["data"].forEach(e => {
					var item = this.createItem(e);
					$('.ChatBody').prepend(item)
				})
			})
	}

	timeout = 250;

	connect = () => {
		var ws = new WebSocket('wss://cmkrosp.iptime.org:8080/room');
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
			//some day refactoring to react style

			var msg = JSON.parse(evt.data)
			console.log("onmessage: ", evt.data);

			//				var elem = <ChatMsg AvatarURL={msg.AvatarURL} Name={msg.Name} Message={msg.Message} When={msg.When} />
			$('.ChatBody').prepend(
				//				ReactDOMServer.renderToStaticMarkup(elem)
				'<div class="MsgContainer"><p class=ChatAvatar><img class=AvatarImg src='+ msg.AvatarURL +' /></p><p class=ChatMsgBody><p class=ChatMsg><strong>' + msg.Name + '</strong>' + ' : ' + msg.Message + '</p><p class=ChatWhen>' + msg.When + '</p></p></div>');
			// Don't use when keep using socket
			//	ws.close();
		}
	}

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

class ChatMsg extends React.Component {
	render() {
		return (
			<div class="MsgContainer">
				<p class="ChatAvatar">
					<img src={this.props.AvatarURL} alt="Avatar" />
				</p>
				<p class="ChatMsgBody">
					<p class="ChatMsg"><strong>{this.props.Name}</strong> : {this.props.Message}</p>
					<p class="ChatWhen">{this.props.When}</p>
				</p>
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
				websocket.send(JSON.stringify({"Message": item}));
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
						<input type="text" id="ChatInput" class="form-control" autocomplete="off" />
						<div id="inputHelpBlock" class="form-text">
							This chat program built by WebSocket
						</div>
						<a href="/logout" class="logout">Logout</a>
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
