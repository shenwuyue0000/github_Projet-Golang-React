import React, { Component } from "react";
import "./ChatHistory.scss";
import Message from "../Message/Message";

class ChatHistory extends Component {
  render() {
    // const messages = this.props.chatHistory.map((msg, index) => (
    //   <p key={index}>{msg.data}</p>
    // ));
    console.log(this.props.chatHistory);
    const messages = this.props.chatHistory.map(msg => <Message message={msg.data} />);


    return (
      <div className="ChatHistory">
        <h2>L'Histoire du Débat</h2>
        {messages}
      </div>
    );
  }
}

export default ChatHistory;
