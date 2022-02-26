import './App.css';
import React, { Component } from "react";
import { connect, sendMsg } from './api';
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
// import ChatInput from './components/ChatInput/ChatInput';

const prisonerBtns = [
  {
      name: 'prisonnier muet',
      value: 1,
  },
  {
      name: 'prisonnier méchant',
      value: 2,
  },
  {
      name: 'prisonnier oeil pour oeil',
      value: 3,
  },
  {
      name: 'prisonnier aléatoire',
      value: 4,
  },
  {
      name: 'prisonnier intelligent',
      value: 5,
  },
  {
      name: 'prisonnier plus intelligent',
      value: 6,
  },
]

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: [],
      clickedBtn: null
    };
  }

  componentDidMount() {
    connect((msg) => {
        console.log("New Message")
        this.setState(prevState => ({
            chatHistory: [...this.state.chatHistory, msg]
        }))
        console.log(this.state.chatHistory);
    });
  }

  send(text) {
    sendMsg(text);
  }
  
  handlePrisonerBtnClick(prisoner) {
    this.send(prisoner.value);
    this.setState({ ...this.state, clickedBtn: prisoner.name });
  }

  render() {
    return (
      <div className="App">
        <Header />
        <div className="btnsContainer">
                    {prisonerBtns.map((prisoner) => (
                        <button
                            className={`prisonerBtn ${this.state.clickedBtn === prisoner.name ? 'selectedPrisoner' : null}`}
                            key={prisoner.name}
                            onClick={() => this.handlePrisonerBtnClick(prisoner)}
                            disabled={this.state.clickedBtn ? true : false}
                        >
                            {prisoner.name}
                        </button>
                    ))}
                </div>
        <ChatHistory chatHistory={this.state.chatHistory} />
        {/* <ChatInput send={this.send} /> */}
      </div>
    );
  }
}

export default App;